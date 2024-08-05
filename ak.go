package akcore

import (
	"errors"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrNotFound = errors.New("not found")
)

/*
func ErrorPB(src *akpb.BusMessage, err error) *akpb.BusMessage {
	return &akpb.BusMessage{
		Type: src.Type,
		Error: &akpb.Error{
			Detail: proto.String(err.Error()),
		},
	}
}

type BusError struct {
	Err *akpb.Error
}

func (err BusError) Error() string {
	if err.Err.UserMessage != nil {
		return *err.Err.UserMessage
	}
	if err.Err.Detail != nil {
		return *err.Err.Detail
	}
	return fmt.Sprint("code ", err.Err.Code)
}
*/

type ProtoMessagePointer[M any] interface {
	*M
	protoreflect.ProtoMessage
}
