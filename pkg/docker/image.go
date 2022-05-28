package docker

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

//https://gist.github.com/miguelmota/4980b18d750fb3b1eb571c3e207b1b92
func (c *Controller) EnsureImage(image string) (err error) {
	reader, err := c.cli.ImagePull(context.Background(), image, types.ImagePullOptions{})

	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)
	return nil
}

func (c *Controller) RemoveImage(image string) (removed bool, err error) {

	removedResponse, err := c.cli.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})

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
