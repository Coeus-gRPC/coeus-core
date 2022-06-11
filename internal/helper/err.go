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

func ErrProtobufParseFailed(fileName string) error {
	if fileName != "" {
		return errors.New(fmt.Sprintf("failed to parse Protobuf file: %s", fileName))
	} else {
		return errors.New("empty path to Protobuf file")
	}
}

func ErrProtobufMethodNotExist(methodName string) error {
	return errors.New(fmt.Sprintf("provided method: `%s` does not exist in protobuf file", methodName))
}

var ErrProtobufMethodIsEmpty = errors.New("no protobuf method is provided, please specify a method")

var ErrInvalidProtobufMethodName = errors.New("invalid protobuf method name provided, please check the config file")
