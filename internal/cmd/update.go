package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ChrisWiegman/kana-cli/internal/console"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

func newUpdateCommand(consoleOutput *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Kana to the latest version.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runUpdate(true, cmd.Use, consoleOutput)
		},
	}

	return cmd
}

// runUpdate checks for updates to the binary and allows for self-updates.
func runUpdate(runUpdate bool, cmd string, consoleOutput *console.Console) {
	if strings.Contains(Version, "-") {
		if runUpdate {
			consoleOutput.Warn("Kana cannot update a development version. Please build your new version of Kana from source.")
		}
		return
	}

	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("ChrisWiegman/kana-cli"))
	if err != nil || !found {
		consoleOutput.Warn("Kana could not check GitHub for the latest version.")
		return
	}

	if runUpdate {
		consoleOutput.Println(fmt.Sprintf("Updating Kana to version %s", latest.Version()))

		if latest.LessOrEqual(Version) {
			consoleOutput.Success(fmt.Sprintf("Current version (%s) is the latest", Version))
		}

		exe, err := os.Executable()
		if err != nil {
			consoleOutput.Error(errors.New("could not locate executable path"))
		}

		err = selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe)
		if err != nil {
			consoleOutput.Error(fmt.Errorf("error occurred while updating binary: %w", err))
		}

		consoleOutput.Success(fmt.Sprintf("Successfully updated to version %s", latest.Version()))
	} else if !latest.LessOrEqual(Version) && cmd != "update" {
		consoleOutput.Warn(fmt.Sprintf("A new version, %s, is available. Run `kana changelog` for details or update via Homebrew with `brew update` and `brew upgrade` or use `kana update` to update.", latest.Version())) //nolint: lll
	}
}
