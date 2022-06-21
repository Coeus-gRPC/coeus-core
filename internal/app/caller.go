package app

import (
	"context"
	"fmt"
	"sync"
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
	Connections []*grpc.ClientConn
	Stubs       []grpcdynamic.Stub
	Method      *desc.MethodDescriptor
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

	opts = append(opts, grpc.WithReturnConnectionError())

	if c.Config.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// TODO: Current implementation does not use an actual tsl cert, it retrieve TLS cert from destination server
		tlsCert, _, _ := getcert.FromTLSServer(c.Config.TargetHost, true)
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewServerTLSFromCert(&tlsCert)))
	}

	if c.Config.Timeout != -1 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Millisecond)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	for i := 0; i < c.Config.Concurrent; i++ {
		conn, err := grpc.DialContext(ctx, c.Config.TargetHost, opts...)
		if err != nil {
			return err
		}

		c.Stubs = append(c.Stubs, grpcdynamic.NewStub(conn))
		c.Connections = append(c.Connections, conn)
	}

	c.RuntimeConfig = runtimeConfig
	c.Method = runtimeConfig.MethodDesc

	return nil
}

func (c *Caller) sendRequest(limiter chan bool, input *dynamic.Message, wg *sync.WaitGroup, count uint) {
	limiter <- true

	stubCount := int(count) % c.Config.Concurrent
	fmt.Printf("Using Stub Num: %d\n", stubCount)

	resp, _ := c.Stubs[stubCount].InvokeRpc(context.Background(), c.Method, input)
	//if err != nil {
	//	return err
	//}

	println(resp.String())

	defer func() {
		<-limiter
		wg.Done()
	}()
}

func (c *Caller) SendRequests(input *dynamic.Message) error {
	var i uint
	ch := make(chan bool, c.Config.Concurrent)
	wg := &sync.WaitGroup{}
	println()

	for i = 0; i < c.Config.TotalCallNum; i++ {
		wg.Add(1)
		go c.sendRequest(ch, input, wg, i)
	}
	wg.Wait()

	return nil
}

func (c *Caller) Run() error {
	start := time.Now()

	input, err := newMessageFromData(c.RuntimeConfig)
	if err != nil {
		return err
	}

	err = c.SendRequests(input)
	if err != nil {
		return err
	}

	total := time.Since(start)
	fmt.Printf("This call cost a total of %d ms.\n", total.Milliseconds())

	return nil
}
