package config

type KanaConfig struct {
	SiteDomain string
	ConfigRoot string
}

func GetKanaConfig() (KanaConfig, error) {

	configRoot, err := GetConfigRoot()
	if err != nil {
		return KanaConfig{}, err
	}

	kanaConfig := KanaConfig{
		SiteDomain: "sites.cfw.li",
		ConfigRoot: configRoot,
	}

	return kanaConfig, nil

}
