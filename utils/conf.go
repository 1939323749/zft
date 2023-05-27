package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Upyun struct {
	Operator  string `json:"operator"`
	Secret    string `json:"secret"`
	Bucket    string `json:"bucket"`
	Bucketurl string `json:"bucketurl"`
}

var (
	defaultConfigFile string
	Conf              Upyun
)

func InitConfig() error {
	// Get the home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to find user home directory: %w", err)
	}

	// Set the default config file to a file in the home directory
	defaultConfigFile = filepath.Join(home, ".zft.json")

	// Check if the config file exists
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		// If not, create an empty one
		f, err := os.Create(defaultConfigFile)
		if err != nil {
			return err
		}
		// Write an empty JSON array to the new file
		if _, err := f.WriteString("{\n    \"operator\":\"\",\n    \"secret\":\"\",\n    \"bucket\":\"\",\n    \"bucketurl\":\"\"\n}"); err != nil {
			return err
		} else {
			return fmt.Errorf("please write configuration file in $HOME/.zft.json")
		}
	}

	viper.SetConfigFile(defaultConfigFile)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %w", err)
	}
	return nil
}

func GetConf() (Upyun, error) {
	err := InitConfig()
	if err != nil {
		return Upyun{}, err
	}
	return Conf, nil
}
