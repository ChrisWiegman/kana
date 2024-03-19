package site

import (
	"strings"

	"github.com/ChrisWiegman/kana/internal/console"
)

// IsXdebugRunning returns true if Xdebug is already running or false if not.
func (s *Site) IsXdebugRunning(consoleOutput *console.Console) bool {
	output, err := s.runCli("pecl list | grep xdebug", false, false)
	if err != nil {
		return false
	}

	return strings.Contains(output.StdOut, "xdebug")
}

// StartXdebug installs and starts xdebug in the site's PHP container.
func (s *Site) StartXdebug(consoleOutput *console.Console) error {
	commands := []string{
		"pecl list | grep xdebug",
		"pecl install xdebug",
		"docker-php-ext-enable xdebug",
		"echo 'xdebug.start_with_request=yes' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.mode=debug,develop,trace' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.client_host=host.docker.internal' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.discover_client_host=on' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.start_with_request=trigger' >> /usr/local/etc/php/php.ini",
		"echo 'xdebug.show_local_vars=1' >> /usr/local/etc/php/php.ini",
		"echo 'html_errors = On' >> /usr/local/etc/php/conf.d/z-custom.ini", // Ensure custom overrides happen
	}

	for i, command := range commands {
		restart := false

		if i+1 == len(commands) {
			restart = true
		}

		output, err := s.runCli(command, restart, true)
		if err != nil {
			return err
		}

		// Verify that the command ran correctly
		if i == 0 && strings.Contains(output.StdOut, "xdebug") {
			return nil
		}
	}

	return nil
}

// StopXdebug stops Xdebug by restarting the WordPress containers.
func (s *Site) StopXdebug(consoleOutput *console.Console) error {
	err := s.stopWordPress()
	if err != nil {
		return err
	}

	return s.startWordPress(consoleOutput)
}
