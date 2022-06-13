package helper

import (
	"errors"
	"fmt"
)

//var ErrEmptyConfigFilePath = errors.New("provided configuration file path is empty")

//var ErrConfigLoadFailed = errors.New("failed to load configuration file")

func ErrConfigLoadFailed(fileName string) error {
	if len(fileName) != 0 {
		return errors.New(fmt.Sprintf("failed to load configuration file: %s", fileName))
	} else {
		return errors.New("empty path to configuration file")
	}
}

var ErrConfigReadFailed = errors.New("failed to read configuration file")

func ErrProtobufParseFailed(fileName string) error {
	if len(fileName) != 0 {
		return errors.New(fmt.Sprintf("failed to parse Protobuf file: %s", fileName))
	} else {
		return errors.New("empty path to Protobuf file")
	}
}

func ErrProtobufServiceNotExist(serviceName string) error {
	return errors.New(fmt.Sprintf("provided service: `%s` does not exist in protobuf file", serviceName))
}

func ErrProtobufMethodNotExist(methodName string) error {
	return errors.New(fmt.Sprintf("provided method: `%s` does not exist in protobuf file", methodName))
}

var ErrProtobufMethodIsEmpty = errors.New("no protobuf method is provided, please specify a method")

var ErrInvalidProtobufMethodName = errors.New("invalid protobuf method name, please check the config file")

func ErrDataFileLoadFailed(fileName string) error {
	if len(fileName) != 0 {
		return errors.New(fmt.Sprintf("failed to load data file: %s", fileName))
	} else {
		return errors.New("empty path to data file")
	}
}

var ErrFailedToParseData = errors.New("failed to parse call data, please check data json file")
