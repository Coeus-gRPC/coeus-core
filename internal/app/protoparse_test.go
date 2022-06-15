package app

import "testing"

var testProtoBufFilePath = "./testdata/proto/greeter.proto"

// TestEmptyProtobufFile calls greetings.Hello with an empty string,
// checking for an error.
func TestEmptyProtobufFile(t *testing.T) {
	des, err := ParseProtobufFile("")
	if err == nil || des != nil {
		t.Errorf(`ParseProtobufFile should return error when empty protobuf file path is provided`)
	}
}

func TestWrongProtobufFile(t *testing.T) {
	des, err := ParseProtobufFile("./nonexist/file")
	if err == nil || des != nil {
		t.Errorf(`WrongProtobufFile should return error when wrong protobuf file path is provided`)
	}
}

func TestCorrectProtobufFile(t *testing.T) {
	des, err := ParseProtobufFile(testProtoBufFilePath)
	if err != nil || des == nil {
		t.Errorf("CorrcetProtobufFile should not yield any error: %s", err.Error())
	}
}

func TestCheckCorrectProtobufMethod(t *testing.T) {
	des, _ := ParseProtobufFile(testProtoBufFilePath)

	periodMethodDes, err := CheckProtobufMethod(des, "greeterservice.Greeter.SayHello")
	if err != nil || periodMethodDes == nil {
		t.Errorf("CheckCorrectProtobufMethod should not yield any error: %s", err.Error())
	}

	slashMethodDes, err := CheckProtobufMethod(des, "greeterservice.Greeter/SayHello")
	if err != nil || slashMethodDes == nil {
		t.Errorf("CheckCorrectProtobufMethod should not yield any error: %s", err.Error())
	}
}

func TestCheckWrongProtobufService(t *testing.T) {
	des, _ := ParseProtobufFile(testProtoBufFilePath)

	methodDes, err := CheckProtobufMethod(des, "wrongpackage.WrongService.SayHello")
	if err == nil || methodDes != nil {
		t.Errorf("CheckWrongProtobufService should return error when a wrong service name is provided")
	}
}

func TestCheckWrongProtobufMethod(t *testing.T) {
	des, _ := ParseProtobufFile(testProtoBufFilePath)

	emptyMethodDes, err := CheckProtobufMethod(des, "")
	if err == nil || emptyMethodDes != nil {
		t.Errorf("CheckWrongProtobufMethod should return error when an empty method name is provided")
	}

	wrongMethodDelimiterDes, err := CheckProtobufMethod(des, "greeterservice&Greeter*WrongMethodName")
	if err == nil || wrongMethodDelimiterDes != nil {
		t.Errorf("CheckWrongProtobufMethod should return error when a wrong delimiter is presented in full method name")
	}

	wrongMethodDes, err := CheckProtobufMethod(des, "greeterservice.Greeter.WrongMethodName")
	if err == nil || wrongMethodDes != nil {
		t.Errorf("CheckWrongProtobufMethod should return error when a wrong method name is provided")
	}
}
