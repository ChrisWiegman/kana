package wordpress

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"
)

func NewWordPress(siteName string) {

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
		HostName:    fmt.Sprintf("kana_%s_mariadb", siteName),
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
		HostName:    fmt.Sprintf("kana_%s_wordpress", siteName),
		Env: []string{
			fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_mariadb", siteName),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", siteName): "web",
			fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", siteName):        fmt.Sprintf("Host(`%s.sites.cfw.li`)", siteName),
			fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", siteName):      "websecure",
			fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", siteName):             fmt.Sprintf("Host(`%s.sites.cfw.li`)", siteName),
			fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", siteName):              "true",
		},
	}

	_, err = controller.ContainerRun(wordPressConfig)
	if err != nil {
		panic(err)
	}

}
