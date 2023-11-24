package settings

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func NewSettings() (*Settings, error) {
	kanaSettings := new(Settings)

	kanaSettings.AppDomain = domain
	kanaSettings.RootKey = rootKey
	kanaSettings.RootCert = rootCert
	kanaSettings.SiteCert = siteCert
	kanaSettings.SiteKey = siteKey

	cwd, err := os.Getwd()
	if err != nil {
		return kanaSettings, err
	}

	kanaSettings.WorkingDirectory = cwd

	home, err := homedir.Dir()
	if err != nil {
		return kanaSettings, err
	}

	kanaSettings.AppDirectory = filepath.Join(home, configFolderName)

	err = kanaSettings.EnsureStaticConfigFiles()

	return kanaSettings, err
}

// GetDefaultPermissions returns the default directory permissions and the default file permissions.
func GetDefaultPermissions() (dirPerms, filePerms int) {
	return defaultDirPermissions, defaultFilePermissions
}
