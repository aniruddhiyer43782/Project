package entity

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// parentSavePath is the main directory to store backup files
const parentSavePath = "backup-x-files"

// Config represents the YAML configuration file
// In Go, struct field names must be capitalized to map to the lowercase keys in config.yml
type Config struct {
	User
	BackupConfig []BackupConfig
	Webhook
	S3Config
	EncryptKey string // Encryption key
}

// cacheType holds the cached configuration
type cacheType struct {
	ConfigSingle *Config
	Err          error
	Lock         sync.Mutex
}

var cache = &cacheType{}

// GetConfigCache retrieves the configuration from cache or file
func GetConfigCache() (conf Config, err error) {

	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	if cache.ConfigSingle != nil {
		return *cache.ConfigSingle, cache.Err
	}

	// Initialize config
	cache.ConfigSingle = &Config{}

	configFilePath := getConfigFilePath()
	_, err = os.Stat(configFilePath)
	if err != nil {
		log.Println("Configuration file not found! Please enter it via the web interface.")
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	byt, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("Failed to read config.yaml")
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	err = yaml.Unmarshal(byt, cache.ConfigSingle)
	if err != nil {
		log.Println("Failed to deserialize configuration file", err)
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	// Clear previous error
	cache.Err = nil
	return *cache.ConfigSingle, err
}

// SaveConfig saves the configuration to file
func (conf *Config) SaveConfig() (err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	byt, err := yaml.Marshal(conf)
	if err != nil {
		log.Println(err)
		return err
	}

	err = ioutil.WriteFile(getConfigFilePath(), byt, 0600)
	if err != nil {
		log.Println(err)
		return err
	}

	// Clear cached configuration
	cache.ConfigSingle = nil

	return
}

// getConfigFilePath returns the path to the config file inside the backup directory
func getConfigFilePath() string {
	_, err := os.Stat(parentSavePath)
	if err != nil {
		os.Mkdir(parentSavePath, 0750)
	}
	return parentSavePath + string(os.PathSeparator) + ".backup_x_config.yaml"
}
