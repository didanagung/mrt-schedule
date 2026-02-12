package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			// Tidak fatal, karena bisa saja production
			log.Println(".env not found, using system environment")
		}
	}
}

func GetString(key string) string {
	val := viper.GetString(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return val
}
