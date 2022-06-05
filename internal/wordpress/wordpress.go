package wordpress

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"
)

func NewWordPress(controller *docker.Controller) error {

	wordPressContainers := []docker.ContainerConfig{
		{
			Name:        fmt.Sprintf("kana_%s_database", controller.Config.CurrentDirectory),
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_database", controller.Config.CurrentDirectory),
			Env: []string{
				"MARIADB_ROOT_PASSWORD=password",
				"MARIADB_DATABASE=wordpress",
				"MARIADB_USER=wordpress",
				"MARIADB_PASSWORD=wordpress",
			},
		},
		{
			Name:        fmt.Sprintf("kana_%s_wordpress", controller.Config.CurrentDirectory),
			Image:       "wordpress",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", controller.Config.CurrentDirectory),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_database", controller.Config.CurrentDirectory),
				"WORDPRESS_DB_USER=wordpress",
				"WORDPRESS_DB_PASSWORD=wordpress",
				"WORDPRESS_DB_NAME=wordpress",
			},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", controller.Config.CurrentDirectory): "web",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", controller.Config.CurrentDirectory):        fmt.Sprintf("Host(`%s.%s`)", controller.Config.CurrentDirectory, controller.Config.SiteDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", controller.Config.CurrentDirectory):      "websecure",
				fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", controller.Config.CurrentDirectory):             fmt.Sprintf("Host(`%s.%s`)", controller.Config.CurrentDirectory, controller.Config.SiteDomain),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", controller.Config.CurrentDirectory):              "true",
			},
		},
	}

	for _, container := range wordPressContainers {

		err := controller.EnsureImage(container.Image)
		if err != nil {
			return err
		}

		_, err = controller.ContainerRun(container)
		if err != nil {
			return err
		}
	}

	return nil

}
