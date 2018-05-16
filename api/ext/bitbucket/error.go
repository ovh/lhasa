package bitbucket

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

// Error options
type Error struct {
	Message string
	Fields  map[string][]string
}

// DecodeError options
func DecodeError(e map[string]interface{}) error {
	var bitbucketError Error
	if err := mapstructure.Decode(e["error"], &bitbucketError); err != nil {
		return err
	}

	return errors.New(bitbucketError.Message)
}
