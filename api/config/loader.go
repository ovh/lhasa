package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ovh/lhasa/api/security"
)

// LoadFromFile extract configuration file
func LoadFromFile(configFile *os.File) (config Lhasa, err error) {
	// Init config file
	b, err := ioutil.ReadFile(configFile.Name())
	if err != nil {
		return Lhasa{}, err
	}
	err = json.Unmarshal(b, &config)
	config.Policy = security.Compile(config.Security)
	return config, err
}
