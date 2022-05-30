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
		HostName:    "kanawordpress",
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.middlewares.wordpress-https.redirectscheme.scheme": "https",
			"traefik.http.routers.wordpress-http.entrypoints":                "web",
			"traefik.http.routers.wordpress-http.rule":                       "Host(`wordpress.local.dev`)",
			"traefik.http.routers.wordpress-http.middlewares":                "wordpress-https@docker",
			"traefik.http.routers.wordpress.entrypoints":                     "web-secure",
			"traefik.http.routers.wordpress.rule":                            "Host(`wordpress.local.dev`)",
			"traefik.http.routers.wordpress.tls":                             "true",
		},
	}

	_, err = controller.ContainerRun(wordPressConfig)
	if err != nil {
		panic(err)
	}

}
