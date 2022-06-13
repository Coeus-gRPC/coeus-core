package app

import (
	"Coeus/internal/helper"
	"encoding/json"
	"github.com/jhump/protoreflect/desc"
	"os"
)

type CoeusRuntimeConfig struct {

	// data
	MethodData []byte
	// methods
	MethodDesc *desc.MethodDescriptor
}

type CoeusConfig struct {
	TotalCallNum    uint   `json:"totalCallNum"`
	Concurrent      int    `json:"concurrent"`
	TargetHost      string `json:"targetHost"`
	Insecure        bool   `json:"insecure"`
	ProtoFile       string `json:"protoFile"`
	MethodName      string `json:"methodName"`
	MessageDataFile string `json:"messageDataFile"`
}

func LoadConfigFromFile(path string, config *CoeusConfig, runtimeConfig *CoeusRuntimeConfig) error {
	jsonConfig, err := os.ReadFile(path)
	if err != nil {
		return helper.ErrConfigLoadFailed(path)
	}

	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		return helper.ErrConfigReadFailed
	}

	protobufPath := config.ProtoFile
	fileDes, err := ParseProtobufFile(protobufPath)
	if err != nil {
		return helper.ErrProtobufParseFailed(protobufPath)
	}

	methodName := config.MethodName
	methodDes, err := CheckProtobufMethod(fileDes, methodName)
	if err != nil {
		return helper.ErrProtobufMethodNotExist(methodName)
	}

	dataFile := config.MessageDataFile
	messageDataByte, err := os.ReadFile(dataFile)
	if err != nil {
		return helper.ErrDataFileLoadFailed(dataFile)
	}

	runtimeConfig.MethodData = messageDataByte
	runtimeConfig.MethodDesc = methodDes

	return nil
}
