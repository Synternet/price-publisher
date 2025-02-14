package config

import (
	"os"

	"github.com/spf13/viper"
)

type CmcConfig struct {
	Ids       string
	IdsSingle string
	ApiKey    string
}

type NatsConfig struct {
	Urls          string
	NKey          string
	PrefixName    string
	PublisherName string
}

type Config struct {
	PublishIntervalSec int64
	NatsConfig         NatsConfig
	CmcConfig          CmcConfig
}

func isDotEnvPresent() bool {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return false
	}

	return true
}

func Init() (config *Config, err error) {
	if isDotEnvPresent() {
		viper.AddConfigPath(".")
		viper.SetConfigName(".env")
		viper.SetConfigType("dotenv")

		err = viper.ReadInConfig()
		if err != nil {
			return
		}
	}

	viper.SetDefault("NATS_URLS", "nats://europe-west3-gcp-dl-testnet-brokernode-frankfurt01.synternet.com")
	viper.SetDefault("NATS_PREFIX_NAME", "syntropy_defi")
	viper.SetDefault("NATS_PUB_NAME", "price")
	viper.SetDefault("PUBLISH_INTERVAL_SEC", 60)
	viper.SetDefault("CMC_IDS_SINGLE", "")

	viper.AutomaticEnv()
	c := parseOsEnv(viper.GetViper())

	return c, err
}

func parseOsEnv(v *viper.Viper) *Config {
	return &Config{
		PublishIntervalSec: v.GetInt64("PUBLISH_INTERVAL_SEC"),
		NatsConfig: NatsConfig{
			Urls:          v.GetString("NATS_URLS"),
			NKey:          v.GetString("NATS_NKEY"),
			PrefixName:    v.GetString("NATS_PREFIX_NAME"),
			PublisherName: v.GetString("NATS_PUB_NAME"),
		},
		CmcConfig: CmcConfig{
			Ids:       v.GetString("CMC_IDS"),
			IdsSingle: v.GetString("CMC_IDS_SINGLE"),
			ApiKey:    v.GetString("CMC_API_KEY"),
		},
	}
}
