package docker

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ChrisWiegman/kana-cli/pkg/console"
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/pkg/jsonmessage"
	"github.com/moby/term"
)

// https://gist.github.com/miguelmota/4980b18d750fb3b1eb571c3e207b1b92
// https://riptutorial.com/docker/example/31980/image-pulling-with-progress-bars--written-in-go
func (d *DockerClient) EnsureImage(imageName string, consoleOutput *console.Console) (err error) {
	if !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:latest", imageName)
	}

	reader, err := d.client.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer reader.Close()

	consoleOutput.Println("Ensuring all images are present and up to date (this may take a few minutes")

	// Discard the download information unless we're debugging
	out := os.Stdout

	if consoleOutput.JSON || !consoleOutput.Debug {
		out, _ = os.Open(os.DevNull)
	}

	termFd, isTerm := term.GetFdInfo(os.Stdout)
	return jsonmessage.DisplayJSONMessagesStream(reader, out, termFd, isTerm, nil)
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
