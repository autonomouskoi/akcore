package akcore

import (
	_ "embed"
	"errors"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:embed VERSION
var Version string

var (
	ErrNotFound = errors.New("not found")
)

func init() {
	Version = strings.TrimSpace(Version)
}

type ProtoMessagePointer[M any] interface {
	*M
	protoreflect.ProtoMessage
}
