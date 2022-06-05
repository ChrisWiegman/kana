package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ChrisWiegman/kana/internal/cursor"
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

//https://gist.github.com/miguelmota/4980b18d750fb3b1eb571c3e207b1b92
//https://riptutorial.com/docker/example/31980/image-pulling-with-progress-bars--written-in-go
func (c *Controller) EnsureImage(imageName string) (err error) {

	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	images, err := c.client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return err
	}

	for _, image := range images {
		for _, imageTag := range image.RepoTags {
			if imageTag == imageName {
				return nil
			}
		}
	}

	events, err := c.client.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer events.Close()

	cursor := cursor.Cursor{}
	layers := make([]string, 0)
	oldIndex := len(layers)

	var event *pullEvent
	decoder := json.NewDecoder(events)

	cursor.Hide()

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
		if strings.HasPrefix(event.Status, "Digest:") || strings.HasPrefix(event.Status, "Status:") {
			fmt.Printf("%s\n", event.Status)
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

		if event.Status == "Pull complete" {
			fmt.Printf("%s: %s\n", event.ID, event.Status)
		} else {
			fmt.Printf("%s: %s %s\n", event.ID, event.Status, event.Progress)
		}
	}

	cursor.Show()

	return nil

}

func (c *Controller) RemoveImage(image string) (removed bool, err error) {

	removedResponse, err := c.client.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})

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
