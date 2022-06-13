package app

import (
	"Coeus/internal/helper"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"os"
)

type Caller struct {
	Config *CoeusRuntimeConfig
	// In future version, we should support concurrency, thus we may need more than one stub
	Stubs  grpcdynamic.Stub
	Method *desc.MethodDescriptor
}

func newMessageFromData(runtimeConfig *CoeusRuntimeConfig) (*dynamic.Message, error) {
	msgDesc := runtimeConfig.MethodDesc.GetInputType()

	dynamicMsg := dynamic.NewMessage(msgDesc) //msgFactory.NewMessage(msgDesc)
	err := dynamicMsg.UnmarshalJSON(runtimeConfig.MethodData)
	if err != nil {
		return &dynamic.Message{}, helper.ErrFailedToParseData
	}

	return dynamicMsg, nil
}

func NewCaller(runtimeConfig *CoeusRuntimeConfig) *Caller {
	c := &Caller{
		Config: runtimeConfig,
		Stubs:  grpcdynamic.NewStub(&grpc.ClientConn{}),
		Method: runtimeConfig.MethodDesc,
	}

	return c
}

func (c *Caller) Run() {
	_, err := newMessageFromData(c.Config)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	println("Caller running!")
}
