package config

type ConfigParams struct {
	AuthToken string
	EncKey    string
	Env       string
	Host      string
	Port      int
}

var Config = &ConfigParams{}

func NewCofig(config *ConfigParams) {
	Config = config
}
