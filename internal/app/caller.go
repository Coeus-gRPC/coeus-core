package app

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
)

type Caller struct {
	// In future version, we should support concurrency, thus we may need more than one stub
	stubs  grpcdynamic.Stub
	method *desc.MethodDescriptor
}

func NewCaller(runtimeConfig *CoeusRuntimeConfig) *Caller {
	c := &Caller{
		stubs:  grpcdynamic.NewStub(&grpc.ClientConn{}),
		method: runtimeConfig.MethodDesc,
	}

	// methodInputs are essentially messages
	methodInputs := c.method.GetInputType()

	dynamic.NewMessage(methodInputs)

	return c
}

func (*Caller) Run() {
	println("Caller running!")
}
