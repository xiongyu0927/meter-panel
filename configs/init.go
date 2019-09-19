package configs

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	log.SetPrefix("[meter-panel]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	viper.SetDefault("MONNITOR_FREQUENCE", 15)
	viper.SetDefault("TTL", 1)
	viper.SetDefault("K8S_TIMEOUT", 10)
	viper.AutomaticEnv()
}
