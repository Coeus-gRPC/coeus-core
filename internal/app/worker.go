package app

import "github.com/jhump/protoreflect/dynamic/grpcdynamic"

type Worker struct {
	ID   int
	Stub grpcdynamic.Stub
}