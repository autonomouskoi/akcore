//*
// Direct messages are sent from AK itself directly to a module. The module can
// receive these messages by subscribing to its own ID has a topic.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.3
// source: bus/direct.proto

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

type MessageTypeDirect int32

const (
	MessageTypeDirect_WEBHOOK_CALL_REQ  MessageTypeDirect = 0
	MessageTypeDirect_WEBHOOK_CALL_RESP MessageTypeDirect = 1
)

// Enum value maps for MessageTypeDirect.
var (
	MessageTypeDirect_name = map[int32]string{
		0: "WEBHOOK_CALL_REQ",
		1: "WEBHOOK_CALL_RESP",
	}
	MessageTypeDirect_value = map[string]int32{
		"WEBHOOK_CALL_REQ":  0,
		"WEBHOOK_CALL_RESP": 1,
	}
)

func (x MessageTypeDirect) Enum() *MessageTypeDirect {
	p := new(MessageTypeDirect)
	*p = x
	return p
}

func (x MessageTypeDirect) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageTypeDirect) Descriptor() protoreflect.EnumDescriptor {
	return file_bus_direct_proto_enumTypes[0].Descriptor()
}

func (MessageTypeDirect) Type() protoreflect.EnumType {
	return &file_bus_direct_proto_enumTypes[0]
}

func (x MessageTypeDirect) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageTypeDirect.Descriptor instead.
func (MessageTypeDirect) EnumDescriptor() ([]byte, []int) {
	return file_bus_direct_proto_rawDescGZIP(), []int{0}
}

type WebhookValues struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *WebhookValues) Reset() {
	*x = WebhookValues{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_direct_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookValues) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookValues) ProtoMessage() {}

func (x *WebhookValues) ProtoReflect() protoreflect.Message {
	mi := &file_bus_direct_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookValues.ProtoReflect.Descriptor instead.
func (*WebhookValues) Descriptor() ([]byte, []int) {
	return file_bus_direct_proto_rawDescGZIP(), []int{0}
}

func (x *WebhookValues) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type WebhookCallRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Params map[string]*WebhookValues `protobuf:"bytes,1,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *WebhookCallRequest) Reset() {
	*x = WebhookCallRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_direct_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookCallRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookCallRequest) ProtoMessage() {}

func (x *WebhookCallRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bus_direct_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookCallRequest.ProtoReflect.Descriptor instead.
func (*WebhookCallRequest) Descriptor() ([]byte, []int) {
	return file_bus_direct_proto_rawDescGZIP(), []int{1}
}

func (x *WebhookCallRequest) GetParams() map[string]*WebhookValues {
	if x != nil {
		return x.Params
	}
	return nil
}

type WebhookCallResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WebhookCallResponse) Reset() {
	*x = WebhookCallResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bus_direct_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookCallResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookCallResponse) ProtoMessage() {}

func (x *WebhookCallResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bus_direct_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookCallResponse.ProtoReflect.Descriptor instead.
func (*WebhookCallResponse) Descriptor() ([]byte, []int) {
	return file_bus_direct_proto_rawDescGZIP(), []int{2}
}

var File_bus_direct_proto protoreflect.FileDescriptor

var file_bus_direct_proto_rawDesc = []byte{
	0x0a, 0x10, 0x62, 0x75, 0x73, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x62, 0x75, 0x73, 0x22, 0x27, 0x0a, 0x0d, 0x57, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x22, 0xa0, 0x01, 0x0a, 0x12, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x43, 0x61, 0x6c, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3b, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x62, 0x75, 0x73, 0x2e, 0x57, 0x65,
	0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x1a, 0x4d, 0x0a, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x62, 0x75, 0x73, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x15, 0x0a, 0x13, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x43, 0x61,
	0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x40, 0x0a, 0x11, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12,
	0x14, 0x0a, 0x10, 0x57, 0x45, 0x42, 0x48, 0x4f, 0x4f, 0x4b, 0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f,
	0x52, 0x45, 0x51, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x57, 0x45, 0x42, 0x48, 0x4f, 0x4f, 0x4b,
	0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x10, 0x01, 0x42, 0x25, 0x5a, 0x23,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6e,
	0x6f, 0x6d, 0x6f, 0x75, 0x73, 0x6b, 0x6f, 0x69, 0x2f, 0x61, 0x6b, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x62, 0x75, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bus_direct_proto_rawDescOnce sync.Once
	file_bus_direct_proto_rawDescData = file_bus_direct_proto_rawDesc
)

func file_bus_direct_proto_rawDescGZIP() []byte {
	file_bus_direct_proto_rawDescOnce.Do(func() {
		file_bus_direct_proto_rawDescData = protoimpl.X.CompressGZIP(file_bus_direct_proto_rawDescData)
	})
	return file_bus_direct_proto_rawDescData
}

var file_bus_direct_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_bus_direct_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_bus_direct_proto_goTypes = []any{
	(MessageTypeDirect)(0),      // 0: bus.MessageTypeDirect
	(*WebhookValues)(nil),       // 1: bus.WebhookValues
	(*WebhookCallRequest)(nil),  // 2: bus.WebhookCallRequest
	(*WebhookCallResponse)(nil), // 3: bus.WebhookCallResponse
	nil,                         // 4: bus.WebhookCallRequest.ParamsEntry
}
var file_bus_direct_proto_depIdxs = []int32{
	4, // 0: bus.WebhookCallRequest.params:type_name -> bus.WebhookCallRequest.ParamsEntry
	1, // 1: bus.WebhookCallRequest.ParamsEntry.value:type_name -> bus.WebhookValues
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_bus_direct_proto_init() }
func file_bus_direct_proto_init() {
	if File_bus_direct_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bus_direct_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*WebhookValues); i {
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
		file_bus_direct_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*WebhookCallRequest); i {
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
		file_bus_direct_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*WebhookCallResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_bus_direct_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bus_direct_proto_goTypes,
		DependencyIndexes: file_bus_direct_proto_depIdxs,
		EnumInfos:         file_bus_direct_proto_enumTypes,
		MessageInfos:      file_bus_direct_proto_msgTypes,
	}.Build()
	File_bus_direct_proto = out.File
	file_bus_direct_proto_rawDesc = nil
	file_bus_direct_proto_goTypes = nil
	file_bus_direct_proto_depIdxs = nil
}
