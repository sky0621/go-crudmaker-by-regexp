package config

import "github.com/spf13/viper"

// Config ...
type Config struct {
	SubA *SubConfigA
	SubB *SubConfigB
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		SubA: NewSubConfigA(),
		SubB: NewSubConfigB(),
	}
}

type SubConfigA struct {
	SubAText string
}

func NewSubConfigA() *SubConfigA {
	return &SubConfigA{
		SubAText: viper.GetString("cfg.suba.text"),
	}
}

type SubConfigB struct {
	SubBText string
}

func NewSubConfigB() *SubConfigB {
	return &SubConfigB{
		/* SubText */
		SubBText: viper.GetString("cfg.subb.text"),
	}
}

// ReadConfig ...
func ReadConfig(configFilePath string) error {
	/*
	 * SetSetSet
	 */
	viper.SetConfigFile(configFilePath)
	return viper.ReadInConfig()
}
