package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChrisWiegman/kana/internal/console"

	"github.com/docker/docker/api/types/image"
	"github.com/moby/moby/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/spf13/viper"
)

var displayJSONMessagesStream = jsonmessage.DisplayJSONMessagesStream

// https://gist.github.com/miguelmota/4980b18d750fb3b1eb571c3e207b1b92
// https://riptutorial.com/docker/example/31980/image-pulling-with-progress-bars--written-in-go
func (d *Client) EnsureImage(imageName string, updateDays int, consoleOutput *console.Console) (err error) {
	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	// Skip more complicated checks if we can
	for _, checkedImage := range d.checkedImages {
		if checkedImage == imageName {
			return nil
		}
	}

	return d.maybeUpdateImage(imageName, updateDays, consoleOutput.JSON)
}

func ValidateImage(imageName, imageTag string) error {
	requestURL := fmt.Sprintf("https://hub.docker.com/v2/namespaces/library/repositories/%s/tags/%s", imageName, imageTag)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var data map[string]interface{}

	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return err
	}

	if data["errinfo"] != nil {
		return fmt.Errorf("image not found for %s:%s", imageName, imageTag)
	}

	return err
}

func (d *Client) maybeUpdateImage(imageName string, updateDays int, suppressOutput bool) error {
	lastUpdated := d.imageUpdateData.GetTime(imageName)

	imageList, err := d.apiClient.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return err
	}

	hasImage := false
	checkForUpdate := false

	// Make sure we've actually downloaded the image
	for i := range imageList {
		imageRepoLabel := imageList[i]
		for _, repoTag := range imageRepoLabel.RepoTags {
			if repoTag == imageName {
				hasImage = true
			}
		}
	}

	// Check the image for updates if needed
	if updateDays > 0 {
		hours := 24 * updateDays
		checkForUpdate = lastUpdated.Compare(time.Now().Add(time.Duration(-hours)*time.Hour)) == -1
	}

	// Pull the image or a newer image if needed
	if !hasImage || checkForUpdate {
		reader, err := d.apiClient.ImagePull(context.Background(), imageName, image.PullOptions{})
		if err != nil {
			return err
		}

		defer func() {
			if err = reader.Close(); err != nil {
				panic(err)
			}
		}()

		out := os.Stdout

		// Discard the download information if set to suppress
		if suppressOutput {
			out, _ = os.Open(os.DevNull)
		}

		d.imageUpdateData.Set(imageName, time.Now())
		err = d.imageUpdateData.WriteConfig()
		if err != nil {
			return err
		}

		termFd, isTerm := term.GetFdInfo(os.Stdout)

		d.checkedImages = append(d.checkedImages, imageName)

		return displayJSONMessagesStream(reader, out, termFd, isTerm, nil)
	}

	d.checkedImages = append(d.checkedImages, imageName)

	return nil
}

func (d *Client) removeImage(imageName string) (removed bool, err error) {
	removedResponse, err := d.apiClient.ImageRemove(context.Background(), imageName, image.RemoveOptions{})

	if err != nil {
		if !strings.Contains(err.Error(), "No such image:") {
			return false, err
		}
	}

	if len(removedResponse) > 0 {
		return true, nil
	}

	return false, nil
}

func (d *Client) loadImageUpdateData(appDirectory string) (*viper.Viper, error) {
	imageUpdateData := viper.New()

	imageUpdateData.SetConfigName("images")
	imageUpdateData.SetConfigType("json")
	imageUpdateData.AddConfigPath(filepath.Join(appDirectory, "config"))

	err := imageUpdateData.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if errors.As(err, &configFileNotFoundError) {
			err = imageUpdateData.SafeWriteConfig()
			if err != nil {
				return imageUpdateData, err
			}
		} else {
			return imageUpdateData, err
		}
	}

	return imageUpdateData, nil
}
