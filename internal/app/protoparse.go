package app

import (
	"Coeus/internal/helper"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

func ParseProtobufFile(path string) (*desc.FileDescriptor, error) {
	if path == "" {
		return nil, helper.ErrProtobufParseFailed(path)
	}

	parser := protoparse.Parser{}
	fileDesc, err := parser.ParseFiles(path)
	if err != nil {
		return nil, helper.ErrProtobufParseFailed(path)
	}

	return fileDesc[0], nil
}

func CheckProtobufMethod(fileDesc *desc.FileDescriptor, methodName string) (*desc.MethodDescriptor, error) {
	// First, parse method to find service name and method name individually
	serviceStr, methodStr, err := parseMethodName(methodName)
	if err != nil {
		return nil, err
	}

	// Then, find the corresponding (service) descriptor
	dsc := fileDesc.FindSymbol(serviceStr)
	if dsc == nil {
		return nil, helper.ErrProtobufServiceNotExist(serviceStr)
	}

	// Then, use the service descriptor to find method
	// Cast the generic descriptor to a service descriptor
	serviceDes := dsc.(*desc.ServiceDescriptor)

	methodDes := serviceDes.FindMethodByName(methodStr)
	if methodDes == nil {
		return nil, helper.ErrProtobufMethodNotExist(methodStr)
	}

	return methodDes, nil
}

// parseMethodName parses the full method name into Package+Service and Method Name
// Valid inputs:
// package.Service.Method
//   .package.Service.Method
//   package.Service/Method
//   .package.Service/Method
// packageName.ServiceName.MethodName
// packageName.ServiceName/MethodName
func parseMethodName(fullMethodName string) (string, string, error) {
	if len(fullMethodName) == 0 {
		return "", "", helper.ErrProtobufMethodIsEmpty
	}

	var delimiter string

	if strings.Count(fullMethodName, ".") == 2 {
		delimiter = "."
	} else if strings.Count(fullMethodName, ".") == 1 && strings.Count(fullMethodName, "/") == 1 {
		delimiter = "/"
	} else {
		return "", "", helper.ErrInvalidProtobufMethodName
	}

	pos := strings.LastIndex(fullMethodName, delimiter)
	if pos == -1 {
		return "", "", helper.ErrInvalidProtobufMethodName
	}

	return fullMethodName[:pos], fullMethodName[pos+1:], nil
}
