package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ChrisWiegman/kana-cli/pkg/console"

	"github.com/docker/docker/api/types"
)

type pullEvent struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

// https://gist.github.com/miguelmota/4980b18d750fb3b1eb571c3e207b1b92
// https://riptutorial.com/docker/example/31980/image-pulling-with-progress-bars--written-in-go
func (d *DockerClient) EnsureImage(imageName string) (err error) {
	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	events, err := d.client.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer events.Close()

	decoder := json.NewDecoder(events)

	err = displayEventStatus(imageName, decoder)

	return err
}

func (d *DockerClient) RemoveImage(image string) (removed bool, err error) {
	removedResponse, err := d.client.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})

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

func displayEventStatus(imageName string, decoder *json.Decoder) error {
	cursor := console.Cursor{}
	layers := make([]string, 0)
	oldIndex := len(layers)

	var event *pullEvent

	cursor.Hide()

	downloading := ""

	for {
		err := decoder.Decode(&event)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		imageID := event.ID

		// Check if the line is one of the final two ones
		if event.Status != "Downloading" && event.Status != "Extracting" {
			continue
		}

		// Check if ID has already passed once
		index := 0
		for i, v := range layers {
			if v == imageID {
				index = i + 1
				break
			}
		}

		// Move the cursor
		if index > 0 {
			diff := index - oldIndex

			if diff > 1 {
				down := diff - 1
				cursor.MoveDown(down)
			} else if diff < 1 {
				up := diff*(-1) + 1
				cursor.MoveUp(up)
			}

			oldIndex = index
		} else {
			layers = append(layers, event.ID)
			diff := len(layers) - oldIndex

			if diff > 1 {
				cursor.MoveDown(diff) // Return to the last row
			}

			oldIndex = len(layers)
		}

		cursor.ClearLine()

		if event.Status == "Downloading" || event.Status == "Extracting" {
			if downloading != imageName {
				fmt.Printf("Downloading latest docker image: %s\n", imageName)
				downloading = imageName
			}

			fmt.Printf("%s: %s %s\n", event.ID, event.Status, event.Progress)
		}
	}

	cursor.Show()

	return nil
}
