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
	viper.SetDefault("K8S_TIMEOUT", 60)
	viper.SetDefault("FURION_HOST", "furion:8080")
	viper.SetDefault("GLOBAL_CLUSTER_NAME", "ace")
	viper.SetDefault("NSXT_GO_ENDPOINT", "http://nsxt-go:12500")
	viper.AutomaticEnv()
}
