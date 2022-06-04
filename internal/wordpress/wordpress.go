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

	err = controller.EnsureImage("mariadb")
	if err != nil {
		panic(err)
	}

	databaseConfig := docker.ContainerConfig{
		Image:       "mariadb",
		NetworkName: "kana",
		HostName:    "kanamariadb",
		Env: []string{
			"MARIADB_ROOT_PASSWORD=password",
			"MARIADB_DATABASE=wordpress",
			"MARIADB_USER=wordpress",
			"MARIADB_PASSWORD=wordpress",
		},
	}

	_, err = controller.ContainerRun(databaseConfig)
	if err != nil {
		panic(err)
	}

	wordPressConfig := docker.ContainerConfig{
		Image:       "wordpress",
		NetworkName: "kana",
		HostName:    "kanawordpress",
		Env: []string{
			"WORDPRESS_DB_HOST=kanamariadb",
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers.wordpress-http.entrypoints": "web",
			"traefik.http.routers.wordpress-http.rule":        "Host(`wordpress.dev.local`)",
			"traefik.http.routers.wordpress.entrypoints":      "websecure",
			"traefik.http.routers.wordpress.rule":             "Host(`wordpress.dev.local`)",
			"traefik.http.routers.wordpress.tls":              "true",
		},
	}

	_, err = controller.ContainerRun(wordPressConfig)
	if err != nil {
		panic(err)
	}

}
