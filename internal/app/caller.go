package app

import (
	"context"
	"os"
	"time"

	"github.com/Coeus-gRPC/coeus-core/internal/helper"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/johnsiilver/getcert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Caller struct {
	// configs
	Config        *CoeusConfig
	RuntimeConfig *CoeusRuntimeConfig
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

func (c *Caller) InitCaller(runtimeConfig *CoeusRuntimeConfig) error {
	var opts []grpc.DialOption
	var ctx context.Context
	var cancel context.CancelFunc

	if c.Config.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// TODO: Current implementation does not use an actual tsl cert, it retrieve TLS cert from destination server
		tlsCert, _, _ := getcert.FromTLSServer(c.Config.TargetHost, true)
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewServerTLSFromCert(&tlsCert)))
	}

	println("Dialing ", c.Config.TargetHost)

	if c.Config.Timeout != -1 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Millisecond)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	conn, err := grpc.DialContext(ctx, c.Config.TargetHost, opts...)
	if err != nil {
		return err
	}

	c.RuntimeConfig = runtimeConfig
	c.Stubs = grpcdynamic.NewStub(conn)
	c.Method = runtimeConfig.MethodDesc

	return nil
}

func (c *Caller) SendRequest(input *dynamic.Message) error {
	resp, err := c.Stubs.InvokeRpc(context.Background(), c.Method, input)
	if err != nil {
		panic(err)
	}

	println(resp.String())

	return nil
}

func (c *Caller) Run() error {
	input, err := newMessageFromData(c.RuntimeConfig)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	err = c.SendRequest(input)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	return nil
}
