package compendium

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func GetFromFile(filePath string) *Compendium {
	var compendium Compendium

	// set the config path for viper to the specified filePath
	viper.SetConfigFile(filePath)

	// attempt to read the specified filePath
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Could not read the compendium from: \""+filePath+"\"", err)
		os.Exit(1)
	}

	// attempt to unmarshal the compendium into a known struct
	if err := viper.Unmarshal(&compendium); err != nil {
		log.Println("Could not parse the compendium from: \"" + filePath + "\"")
		os.Exit(1)
	}

	return &compendium
}
