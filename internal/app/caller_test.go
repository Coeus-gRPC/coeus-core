package app

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	_ "github.com/Coeus-gRPC/coeus-core/test"
	pb "github.com/Coeus-gRPC/coeus-core/test/testdata/proto"
	"google.golang.org/grpc"
)

func generateCorrectConfig(config *CoeusConfig, runtimeConfig *CoeusRuntimeConfig) {
	LoadConfigFromFile(testCorrectConfigFile, config, runtimeConfig)
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Reply: "Hello " + in.GetName()}, nil
}

func DummyServerRun() (*grpc.Server, string) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
		}
	}()

	return s, lis.Addr().String()
}

// TestSuccessfulSecureInit dials the sample host, thus, it needs network connection to succeed
func TestSuccessfulSecureInit(t *testing.T) {
	config := CoeusConfig{}
	runtimeConfig := CoeusRuntimeConfig{}
	generateCorrectConfig(&config, &runtimeConfig)
	caller := Caller{Config: &config}

	// Secure call
	secureErr := caller.InitCaller(&runtimeConfig)
	if secureErr != nil {
		t.Errorf(`TestSuccessfulSecureInit should not fail when correct configs are provided`)
	}
}

func TestSuccessfulInsecureInit(t *testing.T) {
	config := CoeusConfig{}
	runtimeConfig := CoeusRuntimeConfig{}
	generateCorrectConfig(&config, &runtimeConfig)
	caller := Caller{Config: &config}

	service, host := DummyServerRun()
	defer service.Stop()

	// Insecure call
	caller.Config.Insecure = true
	caller.Config.Timeout = -1
	caller.Config.TargetHost = host

	insecureErr := caller.InitCaller(&runtimeConfig)
	if insecureErr != nil {
		t.Errorf(`TestSuccessfulInsecureInit should not fail when correct configs are provided`)
	}
}

func TestFailedInit(t *testing.T) {
	config := CoeusConfig{}
	runtimeConfig := CoeusRuntimeConfig{}
	generateCorrectConfig(&config, &runtimeConfig)
	caller := Caller{Config: &config}

	service, host := DummyServerRun()
	defer service.Stop()

	// This line intends to make the test run faster
	caller.Config.Timeout = 100
	caller.Config.Insecure = true

	// Modify the target host so that it has an incorrect format (lack port number)
	caller.Config.TargetHost = host[:strings.LastIndex(host, ":")]

	err := caller.InitCaller(&runtimeConfig)
	if err == nil {
		t.Errorf(`TestFailedInit should yield err when it cannot dial target host`)
	}
}

func TestSuccessfulRun(t *testing.T) {
	config := CoeusConfig{}
	runtimeConfig := CoeusRuntimeConfig{}
	generateCorrectConfig(&config, &runtimeConfig)
	caller := Caller{Config: &config}

	service, host := DummyServerRun()
	defer service.Stop()

	caller.Config.Insecure = true
	caller.Config.Timeout = -1
	caller.Config.TargetHost = host

	_ = caller.InitCaller(&runtimeConfig)
	err := caller.Run()
	if err != nil {
		t.Errorf(`TestSuccessfulRun should not fail when correct configs are provided`)
	}
}

// func TesFailedRun(t *testing.T) {
// 	// config := &CoeusConfig{}
// 	// runtimeConfig := &CoeusRuntimeConfig{}

// 	// err := LoadConfigFromFile(testCorrectConfigFile, config, runtimeConfig)
// 	// if err != nil {
// 	// 	println(err.Error())
// 	// 	t.Errorf(`TestCorrectConfigFile should not return error when correct config file path is provided`)
// 	// }
// }
