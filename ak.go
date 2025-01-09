package akcore

import (
	_ "embed"
	"errors"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var Version = "DEV"

var (
	ErrNotFound = errors.New("not found")
)

type ProtoMessagePointer[M any] interface {
	*M
	protoreflect.ProtoMessage
}
