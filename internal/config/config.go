package config

type KanaConfig struct {
	SiteDomain string
	ConfigRoot string
	RootCert   string
	RootKey    string
	SiteCert   string
	SiteKey    string
}

func GetKanaConfig() (KanaConfig, error) {

	configRoot, err := GetConfigRoot()
	if err != nil {
		return KanaConfig{}, err
	}

	kanaConfig := KanaConfig{
		SiteDomain: "sites.cfw.li",
		ConfigRoot: configRoot,
		RootKey:    "kana.root.key",
		RootCert:   "kana.root.pem",
		SiteCert:   "kana.site.pem",
		SiteKey:    "kana.site.key",
	}

	return kanaConfig, nil

}
