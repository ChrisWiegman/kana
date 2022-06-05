package config

import (
	"os"
	"path/filepath"
)

type KanaConfig struct {
	SiteDomain       string
	CurrentDirectory string
	ConfigRoot       string
}

func GetKanaConfig() (KanaConfig, error) {

	configRoot, err := GetConfigRoot()
	if err != nil {
		return KanaConfig{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	kanaConfig := KanaConfig{
		SiteDomain:       "sites.cfw.li",
		CurrentDirectory: filepath.Base(cwd),
		ConfigRoot:       configRoot,
	}

	return kanaConfig, nil

}
