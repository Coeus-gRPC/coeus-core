package app

import (
	"context"
	"fmt"
	"sync"
	"time"

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

	return nil
}

func (c *Caller) consumeLatencyChannel(latencies chan time.Duration) *[]time.Duration {
	var latencySlice []time.Duration
	var i uint = 0

	for i < c.Config.TotalCallNum {
		latencySlice = append(latencySlice, <-latencies)
		i++
	}

	return &latencySlice
}

func (c *Caller) sendRequest(limiter chan bool, latencies chan time.Duration, input *dynamic.Message, wg *sync.WaitGroup, count uint) {
	limiter <- true

	stubCount := int(count) % c.Config.Concurrent

	start := time.Now()
	resp, _ := c.Stubs[stubCount].InvokeRpc(context.Background(), c.RuntimeConfig.MethodDesc, input)
	//if err != nil {
	//	return err
	//}

	elapsed := time.Since(start)

	println(resp.String())

	defer func() {
		latencies <- elapsed
		<-limiter
		wg.Done()
	}()
}

func (c *Caller) SendRequests(input *dynamic.Message) *[]time.Duration {
	var i uint
	ch := make(chan bool, c.Config.Concurrent)
	wg := &sync.WaitGroup{}
	latencies := make(chan time.Duration, c.Config.TotalCallNum)

	for i = 0; i < c.Config.TotalCallNum; i++ {
		wg.Add(1)
		i := i
		go c.sendRequest(ch, latencies, input, wg, i)
	}
	wg.Wait()

	latencySlice := c.consumeLatencyChannel(latencies)

	return latencySlice
}

func (c *Caller) Run() error {
	start := time.Now()

	latencies := c.SendRequests(c.RuntimeConfig.MethodMessage)

	fmt.Printf("Latencies: %v\n", latencies)

	total := time.Since(start)
	fmt.Printf("This call cost a total of %d ms.\n", total.Milliseconds())

	return nil
}
