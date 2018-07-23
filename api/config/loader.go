package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ovh/lhasa/api/security"
)

// Empty returns an empty configuration struct
func Empty() Lhasa {
	return Lhasa{
		Policy: make(security.Policy),
	}
}

// LoadFromFile extract configuration file
func LoadFromFile(configFile *os.File) (config Lhasa, err error) {
	config = Empty()
	// Init config file
	b, err := ioutil.ReadFile(configFile.Name())
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(b, &config)
	return config, err
}
