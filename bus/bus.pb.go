// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.3
// source: bus/bus.proto

package bus

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CommonErrorCode int32

const (
	CommonErrorCode_UNKNOWN      CommonErrorCode = 0
	CommonErrorCode_INVALID_TYPE CommonErrorCode = 1
	CommonErrorCode_TIMEOUT      CommonErrorCode = 2
	CommonErrorCode_NOT_FOUND    CommonErrorCode = 3
)

// Enum value maps for CommonErrorCode.
var (
	CommonErrorCode_name = map[int32]string{
		0: "UNKNOWN",
		1: "INVALID_TYPE",
		2: "TIMEOUT",
		3: "NOT_FOUND",
	}
	CommonErrorCode_value = map[string]int32{
		"UNKNOWN":      0,
		"INVALID_TYPE": 1,
		"TIMEOUT":      2,
		"NOT_FOUND":    3,
	}
)

func (x CommonErrorCode) Enum() *CommonErrorCode {
	p := new(CommonErrorCode)
	*p = x
	return p
}

func (x CommonErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommonErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_bus_bus_proto_enumTypes[0].Descriptor()
}

func (CommonErrorCode) Type() protoreflect.EnumType {
	return &file_bus_bus_proto_enumTypes[0]
}

func (x CommonErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommonErrorCode.Descriptor instead.
func (CommonErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{0}
}

type ExternalMessageType int32

const (
	ExternalMessageType_UNSPECIFIED ExternalMessageType = 0
	ExternalMessageType_HAS_TOPIC   ExternalMessageType = 1
	ExternalMessageType_SUBSCRIBE   ExternalMessageType = 2
	ExternalMessageType_UNSUBSCRIBE ExternalMessageType = 3
)

// Enum value maps for ExternalMessageType.
var (
	ExternalMessageType_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "HAS_TOPIC",
		2: "SUBSCRIBE",
		3: "UNSUBSCRIBE",
	}
	ExternalMessageType_value = map[string]int32{
		"UNSPECIFIED": 0,
		"HAS_TOPIC":   1,
		"SUBSCRIBE":   2,
		"UNSUBSCRIBE": 3,
	}
)

func (x ExternalMessageType) Enum() *ExternalMessageType {
	p := new(ExternalMessageType)
	*p = x
	return p
}

func (x ExternalMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ExternalMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_bus_bus_proto_enumTypes[1].Descriptor()
}

func (ExternalMessageType) Type() protoreflect.EnumType {
	return &file_bus_bus_proto_enumTypes[1]
}

func (x ExternalMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ExternalMessageType.Descriptor instead.
func (ExternalMessageType) EnumDescriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{1}
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code           int32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Detail         *string `protobuf:"bytes,2,opt,name=detail,proto3,oneof" json:"detail,omitempty"`
	UserMessage    *string `protobuf:"bytes,3,opt,name=user_message,json=userMessage,proto3,oneof" json:"user_message,omitempty"`
	NotCommonError bool    `protobuf:"varint,4,opt,name=NotCommonError,proto3" json:"NotCommonError,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_bus_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_bus_bus_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{0}
}

func (x *Error) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetDetail() string {
	if x != nil && x.Detail != nil {
		return *x.Detail
	}
	return ""
}

func (x *Error) GetUserMessage() string {
	if x != nil && x.UserMessage != nil {
		return *x.UserMessage
	}
	return ""
}

func (x *Error) GetNotCommonError() bool {
	if x != nil {
		return x.NotCommonError
	}
	return false
}

type BusMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic   string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Type    int32  `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	Error   *Error `protobuf:"bytes,3,opt,name=error,proto3,oneof" json:"error,omitempty"`
	Message []byte `protobuf:"bytes,4,opt,name=message,proto3,oneof" json:"message,omitempty"`
	ReplyTo *int64 `protobuf:"varint,5,opt,name=reply_to,json=replyTo,proto3,oneof" json:"reply_to,omitempty"`
}

func (x *BusMessage) Reset() {
	*x = BusMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_bus_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusMessage) ProtoMessage() {}

func (x *BusMessage) ProtoReflect() protoreflect.Message {
	mi := &file_bus_bus_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusMessage.ProtoReflect.Descriptor instead.
func (*BusMessage) Descriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{1}
}

func (x *BusMessage) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *BusMessage) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *BusMessage) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *BusMessage) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *BusMessage) GetReplyTo() int64 {
	if x != nil && x.ReplyTo != nil {
		return *x.ReplyTo
	}
	return 0
}

type HasTopicRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic     string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	TimeoutMs int32  `protobuf:"varint,2,opt,name=timeout_ms,json=timeoutMs,proto3" json:"timeout_ms,omitempty"`
}

func (x *HasTopicRequest) Reset() {
	*x = HasTopicRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_bus_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HasTopicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HasTopicRequest) ProtoMessage() {}

func (x *HasTopicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bus_bus_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HasTopicRequest.ProtoReflect.Descriptor instead.
func (*HasTopicRequest) Descriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{2}
}

func (x *HasTopicRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *HasTopicRequest) GetTimeoutMs() int32 {
	if x != nil {
		return x.TimeoutMs
	}
	return 0
}

type HasTopicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic    string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	HasTopic bool   `protobuf:"varint,2,opt,name=has_topic,json=hasTopic,proto3" json:"has_topic,omitempty"`
}

func (x *HasTopicResponse) Reset() {
	*x = HasTopicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_bus_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HasTopicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HasTopicResponse) ProtoMessage() {}

func (x *HasTopicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bus_bus_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HasTopicResponse.ProtoReflect.Descriptor instead.
func (*HasTopicResponse) Descriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{3}
}

func (x *HasTopicResponse) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *HasTopicResponse) GetHasTopic() bool {
	if x != nil {
		return x.HasTopic
	}
	return false
}

type SubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_bus_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bus_bus_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{4}
}

func (x *SubscribeRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type UnsubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *UnsubscribeRequest) Reset() {
	*x = UnsubscribeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_bus_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnsubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnsubscribeRequest) ProtoMessage() {}

func (x *UnsubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bus_bus_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnsubscribeRequest.ProtoReflect.Descriptor instead.
func (*UnsubscribeRequest) Descriptor() ([]byte, []int) {
	return file_bus_bus_proto_rawDescGZIP(), []int{5}
}

func (x *UnsubscribeRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

var File_bus_bus_proto protoreflect.FileDescriptor

var file_bus_bus_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x62, 0x75, 0x73, 0x2f, 0x62, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x62, 0x75, 0x73, 0x22, 0xa4, 0x01, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x1b, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12,
	0x26, 0x0a, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0e, 0x4e, 0x6f, 0x74, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0e, 0x4e, 0x6f, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xbf, 0x01, 0x0a, 0x0a,
	0x42, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x62, 0x75, 0x73, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48,
	0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x01, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x08, 0x72, 0x65,
	0x70, 0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x07,
	0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x22, 0x46, 0x0a,
	0x0f, 0x48, 0x61, 0x73, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x5f, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x4d, 0x73, 0x22, 0x45, 0x0a, 0x10, 0x48, 0x61, 0x73, 0x54, 0x6f, 0x70, 0x69,
	0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12,
	0x1b, 0x0a, 0x09, 0x68, 0x61, 0x73, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x68, 0x61, 0x73, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x28, 0x0a, 0x10,
	0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x2a, 0x0a, 0x12, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x2a, 0x4c, 0x0a, 0x0f, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x54, 0x49, 0x4d, 0x45, 0x4f, 0x55, 0x54, 0x10,
	0x02, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x03,
	0x2a, 0x55, 0x0a, 0x13, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x48, 0x41, 0x53, 0x5f,
	0x54, 0x4f, 0x50, 0x49, 0x43, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x55, 0x42, 0x53, 0x43,
	0x52, 0x49, 0x42, 0x45, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x55, 0x42, 0x53,
	0x43, 0x52, 0x49, 0x42, 0x45, 0x10, 0x03, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6e, 0x6f, 0x6d, 0x6f, 0x75, 0x73,
	0x6b, 0x6f, 0x69, 0x2f, 0x61, 0x6b, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x62, 0x75, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bus_bus_proto_rawDescOnce sync.Once
	file_bus_bus_proto_rawDescData = file_bus_bus_proto_rawDesc
)

func file_bus_bus_proto_rawDescGZIP() []byte {
	file_bus_bus_proto_rawDescOnce.Do(func() {
		file_bus_bus_proto_rawDescData = protoimpl.X.CompressGZIP(file_bus_bus_proto_rawDescData)
	})
	return file_bus_bus_proto_rawDescData
}

var file_bus_bus_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_bus_bus_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_bus_bus_proto_goTypes = []any{
	(CommonErrorCode)(0),       // 0: bus.CommonErrorCode
	(ExternalMessageType)(0),   // 1: bus.ExternalMessageType
	(*Error)(nil),              // 2: bus.Error
	(*BusMessage)(nil),         // 3: bus.BusMessage
	(*HasTopicRequest)(nil),    // 4: bus.HasTopicRequest
	(*HasTopicResponse)(nil),   // 5: bus.HasTopicResponse
	(*SubscribeRequest)(nil),   // 6: bus.SubscribeRequest
	(*UnsubscribeRequest)(nil), // 7: bus.UnsubscribeRequest
}
var file_bus_bus_proto_depIdxs = []int32{
	2, // 0: bus.BusMessage.error:type_name -> bus.Error
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_bus_bus_proto_init() }
func file_bus_bus_proto_init() {
	if File_bus_bus_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bus_bus_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Error); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_bus_bus_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*BusMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_bus_bus_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*HasTopicRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_bus_bus_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*HasTopicResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_bus_bus_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_bus_bus_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*UnsubscribeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_bus_bus_proto_msgTypes[0].OneofWrappers = []any{}
	file_bus_bus_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_bus_bus_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bus_bus_proto_goTypes,
		DependencyIndexes: file_bus_bus_proto_depIdxs,
		EnumInfos:         file_bus_bus_proto_enumTypes,
		MessageInfos:      file_bus_bus_proto_msgTypes,
	}.Build()
	File_bus_bus_proto = out.File
	file_bus_bus_proto_rawDesc = nil
	file_bus_bus_proto_goTypes = nil
	file_bus_bus_proto_depIdxs = nil
}
