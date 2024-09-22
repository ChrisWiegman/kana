package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChrisWiegman/kana-wordpress/internal/console"

	"github.com/docker/docker/api/types/image"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/moby/moby/pkg/jsonmessage"
	"github.com/moby/term"
)

var displayJSONMessagesStream = jsonmessage.DisplayJSONMessagesStream

// https://gist.github.com/miguelmota/4980b18d750fb3b1eb571c3e207b1b92
// https://riptutorial.com/docker/example/31980/image-pulling-with-progress-bars--written-in-go
func (d *Client) EnsureImage(imageName, appDirectory string, updateDays int64, consoleOutput *console.Console) (err error) {
	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	// Skip more complicated checks if we can
	for _, checkedImage := range d.checkedImages {
		if checkedImage == imageName {
			return nil
		}
	}

	return d.maybeUpdateImage(imageName, updateDays, consoleOutput.JSON, appDirectory)
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

func (d *Client) maybeUpdateImage(imageName string, updateDays int64, suppressOutput bool, appDirectory string) error {
	lastUpdated := d.imageUpdateData.Time(imageName, time.RFC3339)

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

		err = d.setImageUpdate(imageName, time.Now(), appDirectory)
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

func (d *Client) loadImageUpdateData(appDirectory string) (*koanf.Koanf, error) {
	imageUpdateData := koanf.New(".")

	configFile := filepath.Join(appDirectory, "config", "images.json")
	configFileExists := true

	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		configFileExists = false
	}

	if configFileExists {
		err = imageUpdateData.Load(file.Provider(configFile), kjson.Parser())
		if err != nil {
			return imageUpdateData, err
		}
	}

	return imageUpdateData, nil
}

func (d *Client) setImageUpdate(imageName string, timeStamp time.Time, appDirectory string) error {
	err := d.imageUpdateData.Set(imageName, timeStamp.Format(time.RFC3339))
	if err != nil {
		return err
	}

	configFile := filepath.Join(appDirectory, "config", "images.json")

	f, _ := os.Create(configFile)
	defer f.Close()

	jsonBytes, err := d.imageUpdateData.Marshal(kjson.Parser())
	if err != nil {
		return err
	}

	_, err = f.Write(jsonBytes)

	return err
}
