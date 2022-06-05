package utilities

import (
	"os/exec"
	"runtime"

	"github.com/pkg/browser"
)

func OpenURL(url string) error {

	if runtime.GOOS == "linux" {
		openCmd := exec.Command("xdg-open", url)
		return openCmd.Run()
	}

	return browser.OpenURL(url)
}
