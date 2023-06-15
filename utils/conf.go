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

type Settings struct {
	Ignore []string `json:"ignore"` //"ignore": ["*.tmp", "*.log"]
}

var (
	defaultConfigFile string
	Conf              Upyun
	SettingsConf      Settings
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
		// Write an empty JSON object to the new file
		if _, err := f.WriteString("{\n    \"operator\":\"\",\n    \"secret\":\"\",\n    \"bucket\":\"\",\n    \"bucketurl\":\"\",\n    \"ignore\": []\n}"); err != nil {
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
		return fmt.Errorf("unable to decode Upyun struct, %w", err)
	}

	err = viper.Unmarshal(&SettingsConf)
	if err != nil {
		return fmt.Errorf("unable to decode Settings struct, %w", err)
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

func GetSettings() (Settings, error) {
	err := InitConfig()
	if err != nil {
		return Settings{}, err
	}
	return SettingsConf, nil
}
func SetConf(upyun Upyun) error {
	err := InitConfig()
	if err != nil {
		return err
	}

	// Update the Conf
	Conf = upyun

	// Update the config file with the new Upyun fields
	err = updateConfigFile("operator", Conf.Operator)
	if err != nil {
		return err
	}
	err = updateConfigFile("secret", Conf.Secret)
	if err != nil {
		return err
	}
	err = updateConfigFile("bucket", Conf.Bucket)
	if err != nil {
		return err
	}
	err = updateConfigFile("bucketurl", Conf.Bucketurl)
	if err != nil {
		return err
	}

	return nil
}

func SetSettings(settings Settings) error {
	err := InitConfig()
	if err != nil {
		return err
	}

	// Update the SettingsConf
	SettingsConf = settings

	// Update the config file with the new Settings fields
	err = updateConfigFile("ignore", SettingsConf.Ignore)
	if err != nil {
		return err
	}

	return nil
}

func updateConfigFile(key string, value interface{}) error {
	viper.Set(key, value)

	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("unable to write config file, %w", err)
	}

	return nil
}
