package wordpress

import (
	"github.com/ChrisWiegman/kana/internal/docker"
)

func NewWordPress() {

	controller, err := docker.NewController()
	if err != nil {
		panic(err)
	}

	_, _, err = controller.EnsureNetwork("kana")
	if err != nil {
		panic(err)
	}

	err = controller.EnsureImage("wordpress")
	if err != nil {
		panic(err)
	}

	wordPressConfig := docker.ContainerConfig{
		Image:       "wordpress",
		NetworkName: "kana",
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers.wordpress.entrypoints": "web",
		},
	}

	_, err = controller.ContainerRun(wordPressConfig)
	if err != nil {
		panic(err)
	}

}
