package helper

import (
	"errors"
	"fmt"
)

//var ErrEmptyConfigFilePath = errors.New("provided configuration file path is empty")

//var ErrConfigLoadFailed = errors.New("failed to load configuration file")

func ErrConfigLoadFailed(fileName *string) error {
	if fileName != nil {
		return errors.New(fmt.Sprintf("failed to load configuration file: %s", *fileName))
	} else {
		return errors.New("empty path to configuration file")
	}
}

var ErrConfigReadFailed = errors.New("failed to read configuration file")
