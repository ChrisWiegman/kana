package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

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
			latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("ChrisWiegman/kana-cli"))
			if err != nil {
				consoleOutput.Error(fmt.Errorf("error occurred while detecting version: %w", err))
			}

			if !found {
				consoleOutput.Error(fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH))
			}

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
		},
	}

	return cmd
}
