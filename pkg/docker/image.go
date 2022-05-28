package docker

import (
	"context"
	"io"
	"os"

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
