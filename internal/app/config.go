package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jhump/protoreflect/dynamic"
	"os"

	"github.com/Coeus-gRPC/coeus-core/internal/helper"

	"github.com/jhump/protoreflect/desc"
)

type CoeusRuntimeConfig struct {
	// data
	MethodMessage *dynamic.Message
	// methods
	MethodDesc *desc.MethodDescriptor
}

type CoeusConfig struct {
	ID              uuid.UUID `json:"id"`
	TotalCallNum    int       `json:"totalCallNum"`
	Concurrent      int       `json:"concurrent"`
	TargetHost      string    `json:"targetHost"`
	Insecure        bool      `json:"insecure"`
	Timeout         int       `json:"timeout"`
	ProtoFile       string    `json:"protoFile"`
	MethodName      string    `json:"methodName"`
	MessageDataFile string    `json:"messageDataFile"`
	OutputFilePath  string    `json:"outputFilePath"`
}

func NewMessageFromData(methodDes *desc.MethodDescriptor, messageData []byte) (*dynamic.Message, error) {
	msgDesc := methodDes.GetInputType()

	dynamicMsg := dynamic.NewMessage(msgDesc) //msgFactory.NewMessage(msgDesc)
	err := dynamicMsg.UnmarshalJSON(messageData)
	if err != nil {
		return &dynamic.Message{}, helper.ErrFailedToParseData
	}

	return dynamicMsg, nil
}

func LoadConfigFromFile(path string, config *CoeusConfig, runtimeConfig *CoeusRuntimeConfig) error {
	jsonConfig, err := os.ReadFile(path)
	if err != nil {
		println(err.Error())
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

	methodMessage, err := NewMessageFromData(methodDes, messageDataByte)
	if err != nil {
		return helper.ErrFailedToGenerateProtobufMessage(methodDes.String())
	}

	runtimeConfig.MethodMessage = methodMessage
	runtimeConfig.MethodDesc = methodDes

	return nil
}
