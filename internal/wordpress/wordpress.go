package wordpress

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/traefik"
)

func StopWordPress(controller *docker.Controller) error {

	wordPressContainers := []string{
		fmt.Sprintf("kana_%s_database", controller.Config.CurrentDirectory),
		fmt.Sprintf("kana_%s_wordpress", controller.Config.CurrentDirectory),
	}

	for _, wordPressContainer := range wordPressContainers {
		_, err := controller.ContainerStop(wordPressContainer)
		if err != nil {
			return err
		}
	}

	return traefik.MaybeStopTraefik(controller)

}

func StartWordPress(controller *docker.Controller) error {

	_, _, err := controller.EnsureNetwork("kana")
	if err != nil {
		return err
	}

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
			Labels: map[string]string{
				"kana.site": controller.Config.CurrentDirectory,
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
				"kana.site": controller.Config.CurrentDirectory,
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
