// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.3
// source: webclient.proto

package svc

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

type WebclientStaticDownloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	URL string `protobuf:"bytes,1,opt,name=URL,proto3" json:"URL,omitempty"`
}

func (x *WebclientStaticDownloadRequest) Reset() {
	*x = WebclientStaticDownloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webclient_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebclientStaticDownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebclientStaticDownloadRequest) ProtoMessage() {}

func (x *WebclientStaticDownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_webclient_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebclientStaticDownloadRequest.ProtoReflect.Descriptor instead.
func (*WebclientStaticDownloadRequest) Descriptor() ([]byte, []int) {
	return file_webclient_proto_rawDescGZIP(), []int{0}
}

func (x *WebclientStaticDownloadRequest) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

type WebclientStaticDownloadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *WebclientStaticDownloadResponse) Reset() {
	*x = WebclientStaticDownloadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_webclient_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebclientStaticDownloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebclientStaticDownloadResponse) ProtoMessage() {}

func (x *WebclientStaticDownloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_webclient_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebclientStaticDownloadResponse.ProtoReflect.Descriptor instead.
func (*WebclientStaticDownloadResponse) Descriptor() ([]byte, []int) {
	return file_webclient_proto_rawDescGZIP(), []int{1}
}

func (x *WebclientStaticDownloadResponse) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

var File_webclient_proto protoreflect.FileDescriptor

var file_webclient_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x77, 0x65, 0x62, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x73, 0x76, 0x63, 0x22, 0x32, 0x0a, 0x1e, 0x57, 0x65, 0x62, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x52, 0x4c, 0x22, 0x35, 0x0a, 0x1f, 0x57, 0x65,
	0x62, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74,
	0x68, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x75, 0x74, 0x6f, 0x6e, 0x6f, 0x6d, 0x6f, 0x75, 0x73, 0x6b, 0x6f, 0x69, 0x2f, 0x61, 0x6b,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x73, 0x76, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x73, 0x76, 0x63, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_webclient_proto_rawDescOnce sync.Once
	file_webclient_proto_rawDescData = file_webclient_proto_rawDesc
)

func file_webclient_proto_rawDescGZIP() []byte {
	file_webclient_proto_rawDescOnce.Do(func() {
		file_webclient_proto_rawDescData = protoimpl.X.CompressGZIP(file_webclient_proto_rawDescData)
	})
	return file_webclient_proto_rawDescData
}

var file_webclient_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_webclient_proto_goTypes = []any{
	(*WebclientStaticDownloadRequest)(nil),  // 0: svc.WebclientStaticDownloadRequest
	(*WebclientStaticDownloadResponse)(nil), // 1: svc.WebclientStaticDownloadResponse
}
var file_webclient_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_webclient_proto_init() }
func file_webclient_proto_init() {
	if File_webclient_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_webclient_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*WebclientStaticDownloadRequest); i {
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
		file_webclient_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*WebclientStaticDownloadResponse); i {
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
			RawDescriptor: file_webclient_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_webclient_proto_goTypes,
		DependencyIndexes: file_webclient_proto_depIdxs,
		MessageInfos:      file_webclient_proto_msgTypes,
	}.Build()
	File_webclient_proto = out.File
	file_webclient_proto_rawDesc = nil
	file_webclient_proto_goTypes = nil
	file_webclient_proto_depIdxs = nil
}
