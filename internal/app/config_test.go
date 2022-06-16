package app

import (
	"testing"

	_ "github.com/Coeus-gRPC/coeus-core/test"
)

var testCorrectConfigFile = "./test/testdata/config/testconfig.json"
var testMalformedConfigFile = "./test/testdata/config/testconfig_err.json"
var testEmptyConfigFile = "./test/testdata/config/empty_config.json"
var testIncorrectProtobufMethodConfig = "./test/testdata/config/testconfig_err_method.json"
var testIncorrectDataFileConfig = "./test/testdata/config/testconfig_err_data.json"

func TestCorrectConfigFile(t *testing.T) {
	config := &CoeusConfig{}
	runtimeConfig := &CoeusRuntimeConfig{}

	err := LoadConfigFromFile(testCorrectConfigFile, config, runtimeConfig)
	if err != nil {
		println(err.Error())
		t.Errorf(`TestCorrectConfigFile should not return error when correct config file path is provided`)
	}
}

func TestLoadWrongConfigFile(t *testing.T) {
	config := &CoeusConfig{}
	runtimeConfig := &CoeusRuntimeConfig{}

	nonexistErr := LoadConfigFromFile("./nonexistConfig.json", config, runtimeConfig)
	if nonexistErr == nil {
		t.Errorf(`tLoadWrongConfigFile should return error when incorrect config file path is provided`)
	}

	emptyErr := LoadConfigFromFile(testEmptyConfigFile, config, runtimeConfig)
	if emptyErr == nil {
		t.Errorf(`tLoadWrongConfigFile should return error when an empty config file is provided`)
	}

	malformedErr := LoadConfigFromFile(testMalformedConfigFile, config, runtimeConfig)
	if malformedErr == nil {
		t.Errorf(`tLoadWrongConfigFile should return error when malformed config file is provided`)
	}

	wrongProtobufMethodErr := LoadConfigFromFile(testIncorrectProtobufMethodConfig, config, runtimeConfig)
	if wrongProtobufMethodErr == nil {
		t.Errorf(`tLoadWrongConfigFile should return error when wrong method name is presented in config file`)
	}

	wrongDataErr := LoadConfigFromFile(testIncorrectDataFileConfig, config, runtimeConfig)
	if wrongDataErr == nil {
		t.Errorf(`tLoadWrongConfigFile should return error when incorrect data file is presented in config file`)
	}
}
