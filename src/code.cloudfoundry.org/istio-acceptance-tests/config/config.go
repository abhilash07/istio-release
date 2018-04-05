package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ApiEndpoint   string `json:"cf_api"`
	AdminUser     string `json:"cf_admin_user"`
	AdminPassword string `json:"cf_admin_password"`
	AppsDomain    string `json:"cf_apps_domain"`
}

func NewConfig(path string) (Config, error) {
	var config Config

	configFile, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	return config, err
}

func (c Config) Validate() error {
	missingProperties := []string{}
	if c.ApiEndpoint == "" {
		missingProperties = append(missingProperties, "cf_api")
	}
	if c.AdminUser == "" {
		missingProperties = append(missingProperties, "cf_admin_user")
	}
	if c.AdminPassword == "" {
		missingProperties = append(missingProperties, "cf_admin_password")
	}
	if c.AppsDomain == "" {
		missingProperties = append(missingProperties, "cf_apps_domain")
	}
	if len(missingProperties) > 0 {
		return errors.New(fmt.Sprintf("Missing required config properties: %s", strings.Join(missingProperties, ", ")))
	} else {
		return nil
	}
}
