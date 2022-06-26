package app

import (
	"context"
	"fmt"
	"github.com/Coeus-gRPC/coeus-core/internal/report"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/johnsiilver/getcert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"sync"
	"time"
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

func (c *Caller) consumeReporterChannel(reporters chan report.Reporter) []report.Reporter {
	var reporterSlice []report.Reporter

	for i := 0; i < c.Config.TotalCallNum; i++ {
		reporterSlice = append(reporterSlice, <-reporters)
	}

	return reporterSlice
}

func (c *Caller) sendRequest(limiter chan bool, reporters chan report.Reporter, input *dynamic.Message, wg *sync.WaitGroup, count int) {
	limiter <- true
	stubCount := count % c.Config.Concurrent
	reporter := report.NewReporter()

	start := time.Now()
	resp, err := c.Stubs[stubCount].InvokeRpc(context.Background(), c.RuntimeConfig.MethodDesc, input)

	elapsed := time.Since(start)

	reporter.TimeConsumption = float64(elapsed.Microseconds()) / 1000
	reporter.StatusStr = status.Code(err).String()
	reporter.ReturnStr = resp.String()

	defer func() {
		reporters <- reporter
		<-limiter
		wg.Done()
	}()
}

func (c *Caller) SendRequests(input *dynamic.Message) []report.Reporter {
	ch := make(chan bool, c.Config.Concurrent)
	wg := &sync.WaitGroup{}

	reporters := make(chan report.Reporter, c.Config.TotalCallNum)

	for i := 0; i < c.Config.TotalCallNum; i++ {
		wg.Add(1)
		i := i
		go c.sendRequest(ch, reporters, input, wg, i)
	}
	wg.Wait()

	reporterSlice := c.consumeReporterChannel(reporters)

	return reporterSlice
}

func (c *Caller) Run() error {
	start := time.Now()

	reports := c.SendRequests(c.RuntimeConfig.MethodMessage)

	finalReport := report.GenerateReport(reports, time.Since(start), c.Config.Concurrent)
	fmt.Printf("%v\n", finalReport)

	return nil
}
