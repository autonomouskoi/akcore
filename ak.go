package akcore

import (
	_ "embed"
	"errors"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var Version = "DEV"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrForbidden  = errors.New("forbidden")
)

type ProtoMessagePointer[M any] interface {
	*M
	protoreflect.ProtoMessage
}
