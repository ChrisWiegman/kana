package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

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
			checkForUpdate(true, consoleOutput)
		},
	}

	return cmd
}

// checkForUpdate checks for updates to the binary and allows for self-updates.
func checkForUpdate(runUpdate bool, consoleOutput *console.Console) {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("ChrisWiegman/kana-cli"))
	if err != nil || !found {
		consoleOutput.Warn("Kana could not check GitHub for the latest version.")
	}

	if !latest.LessOrEqual(Version) {
		consoleOutput.Warn(fmt.Sprintf("A new version, %s, is available. Run `kana changelog` for details or `kana update` to update.", latest.Version())) //nolint: lll
	}

	if runUpdate {
		if latest.LessOrEqual(Version) {
			consoleOutput.Success(fmt.Sprintf("Current version (%s) is the latest", Version))
		}

		exe, err := os.Executable()
		if err != nil {
			consoleOutput.Error(errors.New("could not locate executable path"))
		}

		if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
			consoleOutput.Error(fmt.Errorf("error occurred while updating binary: %w", err))
		}

		consoleOutput.Success(fmt.Sprintf("Successfully updated to version %s", latest.Version()))
	}
}
