package docker

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/spf13/viper"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/pkg/jsonmessage"
	"github.com/moby/term"
)

var displayJSONMessagesStream = jsonmessage.DisplayJSONMessagesStream

// https://gist.github.com/miguelmota/4980b18d750fb3b1eb571c3e207b1b92
// https://riptutorial.com/docker/example/31980/image-pulling-with-progress-bars--written-in-go
func (d *DockerClient) EnsureImage(imageName string, updateDays int, consoleOutput *console.Console) (err error) {
	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	lastUpdated := d.imageUpdateData.GetTime(imageName)

	imageList, err := d.moby.ImageList(context.Background(), types.ImageListOptions{})

	hasImage := false
	checkForUpdate := false

	for i := range imageList {
		imageRepoLabel := imageList[i]
		for _, repoTag := range imageRepoLabel.RepoTags {
			if repoTag == imageName {
				hasImage = true
			}
		}
	}

	if updateDays > 0 {
		hours := 24 * updateDays
		checkForUpdate = lastUpdated.Compare(time.Now().Add(time.Duration(-hours)*time.Hour)) == -1
	}

	if !hasImage || checkForUpdate {
		reader, err := d.moby.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
		if err != nil {
			return err
		}

		defer reader.Close()

		// Discard the download information unless we're debugging
		out := os.Stdout

		d.imageUpdateData.Set(imageName, time.Now())
		err = d.imageUpdateData.WriteConfig()
		if err != nil {
			return err
		}

		termFd, isTerm := term.GetFdInfo(os.Stdout)
		return displayJSONMessagesStream(reader, out, termFd, isTerm, nil)
	}

	return nil
}

func (d *DockerClient) removeImage(image string) (removed bool, err error) {
	removedResponse, err := d.moby.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})

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

func (d *DockerClient) loadImageUpdateData(appDirectory string) (*viper.Viper, error) {
	imageUpdateData := viper.New()

	imageUpdateData.SetConfigName("images")
	imageUpdateData.SetConfigType("json")
	imageUpdateData.AddConfigPath(path.Join(appDirectory, "config"))

	err := imageUpdateData.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
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
