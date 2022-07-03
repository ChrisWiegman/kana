package setup

import "github.com/ChrisWiegman/kana/internal/docker"

func SetupApp(controller *docker.Controller) {

	EnsureAppConfig(controller.Config)
	EnsureCerts(controller.Config)

}
