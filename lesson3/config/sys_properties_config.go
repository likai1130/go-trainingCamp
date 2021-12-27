package config

import (
	"github.com/flyleft/gprofile"
	"log"
	"os"
)

var AppConfig = ApplicationConfig{}

func InitConfig() {
	getwd, _ := os.Getwd()
	config, err := gprofile.Profile(&ApplicationConfig{}, getwd + "/lesson3/application.yaml", true)
	if err != nil {
		log.Fatalf("Profile execute error: %s", err.Error())
	}
	AppConfig = *config.(*ApplicationConfig)

}
