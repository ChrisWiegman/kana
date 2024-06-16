package site

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/docker"
)

type Cli interface {
	WordPress(command string, restart, root bool) (docker.ExecResult, error)
	WPCli(command []string, interactive bool, consoleOutput *console.Console) (statusCode int64, output string, err error)
	Command(name string, arg ...string) *exec.Cmd
}

func Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// RunWPCli Runs a wp-cli command returning it's output and any errors.
func (s *Site) WPCli(command []string, interactive bool, consoleOutput *console.Console) (statusCode int64, output string, err error) {
	mounts := s.dockerClient.ContainerGetMounts(fmt.Sprintf("kana-%s-wordpress", s.Settings.Name))

	for _, mount := range mounts {
		if strings.Contains(mount.Destination, "/var/www/html/wp-content/plugins/") {
			s.Settings.Type = "plugin" //nolint
		}

		if strings.Contains(mount.Destination, "/var/www/html/wp-content/themes/") {
			s.Settings.Type = "theme" //nolint
		}
	}

	wordPressDirectory, err := s.getWordPressDirectory()
	if err != nil {
		return 1, "", err
	}

	appVolumes, err := s.getWordPressMounts(wordPressDirectory)
	if err != nil {
		return 1, "", err
	}

	fullCommand := []string{
		"wp",
		"--path=/var/www/html",
	}

	fullCommand = append(fullCommand, command...)

	envVars := []string{
		"IS_KANA_ENVIRONMENT=true",
	}

	isUsingSQLite, err := s.isUsingSQLite()
	if err != nil {
		return 1, "", err
	}

	if isUsingSQLite {
		envVars = append(envVars, "KANA_SQLITE=true")
	} else {
		envVars = append(envVars,
			fmt.Sprintf("WORDPRESS_DB_HOST=kana-%s-database", s.Settings.Name),
			"WORDPRESS_DB_USER=wordpress",
			"WORDPRESS_DB_PASSWORD=wordpress",
			"WORDPRESS_DB_NAME=wordpress",
			"WORDPRESS_ADMIN_USER=admin")
	}

	container := docker.ContainerConfig{
		Name:        fmt.Sprintf("kana-%s-wordpress_cli", s.Settings.Name),
		Image:       fmt.Sprintf("wordpress:cli-php%s", s.Settings.PHP),
		NetworkName: "kana",
		HostName:    fmt.Sprintf("kana-%s-wordpress_cli", s.Settings.Name),
		Command:     fullCommand,
		Env:         envVars,
		Labels: map[string]string{
			"kana.site": s.Settings.Name,
		},
		Volumes: appVolumes,
	}

	if s.Settings.AutomaticLogin {
		container.Env = append(container.Env, "KANA_ADMIN_LOGIN=true")
	}

	err = s.dockerClient.EnsureImage(container.Image, s.Settings.ImageUpdateDays, consoleOutput)
	if err != nil {
		return 1, "", err
	}

	code, output, err := s.dockerClient.ContainerRunAndClean(&container, interactive)
	if err != nil {
		return code, "", err
	}

	return code, output, nil
}

// runCli Runs an arbitrary CLI command against the site's WordPress container.
func (s *Site) WordPress(command string, restart, root bool) (docker.ExecResult, error) {
	container := fmt.Sprintf("kana-%s-wordpress", s.Settings.Name)

	output, err := s.dockerClient.ContainerExec(container, root, []string{command})
	if err != nil {
		return docker.ExecResult{}, err
	}

	if restart {
		_, err = s.dockerClient.ContainerRestart(container)
		return output, err
	}

	return output, nil
}
