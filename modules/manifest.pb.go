// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.3
// source: modules/manifest.proto

package modules

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

type ManifestWebPathType int32

const (
	ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_GENERAL       ManifestWebPathType = 0
	ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_OBS_OVERLAY   ManifestWebPathType = 1
	ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_EMBED_CONTROL ManifestWebPathType = 2
)

// Enum value maps for ManifestWebPathType.
var (
	ManifestWebPathType_name = map[int32]string{
		0: "MANIFEST_WEB_PATH_TYPE_GENERAL",
		1: "MANIFEST_WEB_PATH_TYPE_OBS_OVERLAY",
		2: "MANIFEST_WEB_PATH_TYPE_EMBED_CONTROL",
	}
	ManifestWebPathType_value = map[string]int32{
		"MANIFEST_WEB_PATH_TYPE_GENERAL":       0,
		"MANIFEST_WEB_PATH_TYPE_OBS_OVERLAY":   1,
		"MANIFEST_WEB_PATH_TYPE_EMBED_CONTROL": 2,
	}
)

func (x ManifestWebPathType) Enum() *ManifestWebPathType {
	p := new(ManifestWebPathType)
	*p = x
	return p
}

func (x ManifestWebPathType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ManifestWebPathType) Descriptor() protoreflect.EnumDescriptor {
	return file_modules_manifest_proto_enumTypes[0].Descriptor()
}

func (ManifestWebPathType) Type() protoreflect.EnumType {
	return &file_modules_manifest_proto_enumTypes[0]
}

func (x ManifestWebPathType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ManifestWebPathType.Descriptor instead.
func (ManifestWebPathType) EnumDescriptor() ([]byte, []int) {
	return file_modules_manifest_proto_rawDescGZIP(), []int{0}
}

type ManifestWebPath struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path        string              `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Type        ManifestWebPathType `protobuf:"varint,2,opt,name=type,proto3,enum=modules.ManifestWebPathType" json:"type,omitempty"`
	Description string              `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *ManifestWebPath) Reset() {
	*x = ManifestWebPath{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modules_manifest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManifestWebPath) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManifestWebPath) ProtoMessage() {}

func (x *ManifestWebPath) ProtoReflect() protoreflect.Message {
	mi := &file_modules_manifest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManifestWebPath.ProtoReflect.Descriptor instead.
func (*ManifestWebPath) Descriptor() ([]byte, []int) {
	return file_modules_manifest_proto_rawDescGZIP(), []int{0}
}

func (x *ManifestWebPath) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ManifestWebPath) GetType() ManifestWebPathType {
	if x != nil {
		return x.Type
	}
	return ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_GENERAL
}

func (x *ManifestWebPath) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type Manifest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string             `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string             `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	WebPaths    []*ManifestWebPath `protobuf:"bytes,4,rep,name=web_paths,json=webPaths,proto3" json:"web_paths,omitempty"`
}

func (x *Manifest) Reset() {
	*x = Manifest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modules_manifest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Manifest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Manifest) ProtoMessage() {}

func (x *Manifest) ProtoReflect() protoreflect.Message {
	mi := &file_modules_manifest_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Manifest.ProtoReflect.Descriptor instead.
func (*Manifest) Descriptor() ([]byte, []int) {
	return file_modules_manifest_proto_rawDescGZIP(), []int{1}
}

func (x *Manifest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Manifest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Manifest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Manifest) GetWebPaths() []*ManifestWebPath {
	if x != nil {
		return x.WebPaths
	}
	return nil
}

var File_modules_manifest_proto protoreflect.FileDescriptor

var file_modules_manifest_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65,
	0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x73, 0x22, 0x79, 0x0a, 0x0f, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x57, 0x65, 0x62,
	0x50, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x30, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73,
	0x2e, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x57, 0x65, 0x62, 0x50, 0x61, 0x74, 0x68,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x87, 0x01, 0x0a,
	0x08, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x35, 0x0a, 0x09, 0x77, 0x65, 0x62, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2e, 0x4d, 0x61, 0x6e,
	0x69, 0x66, 0x65, 0x73, 0x74, 0x57, 0x65, 0x62, 0x50, 0x61, 0x74, 0x68, 0x52, 0x08, 0x77, 0x65,
	0x62, 0x50, 0x61, 0x74, 0x68, 0x73, 0x2a, 0x8b, 0x01, 0x0a, 0x13, 0x4d, 0x61, 0x6e, 0x69, 0x66,
	0x65, 0x73, 0x74, 0x57, 0x65, 0x62, 0x50, 0x61, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65, 0x12, 0x22,
	0x0a, 0x1e, 0x4d, 0x41, 0x4e, 0x49, 0x46, 0x45, 0x53, 0x54, 0x5f, 0x57, 0x45, 0x42, 0x5f, 0x50,
	0x41, 0x54, 0x48, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x4c,
	0x10, 0x00, 0x12, 0x26, 0x0a, 0x22, 0x4d, 0x41, 0x4e, 0x49, 0x46, 0x45, 0x53, 0x54, 0x5f, 0x57,
	0x45, 0x42, 0x5f, 0x50, 0x41, 0x54, 0x48, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x42, 0x53,
	0x5f, 0x4f, 0x56, 0x45, 0x52, 0x4c, 0x41, 0x59, 0x10, 0x01, 0x12, 0x28, 0x0a, 0x24, 0x4d, 0x41,
	0x4e, 0x49, 0x46, 0x45, 0x53, 0x54, 0x5f, 0x57, 0x45, 0x42, 0x5f, 0x50, 0x41, 0x54, 0x48, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x45, 0x4d, 0x42, 0x45, 0x44, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52,
	0x4f, 0x4c, 0x10, 0x02, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6e, 0x6f, 0x6d, 0x6f, 0x75, 0x73, 0x6b, 0x6f, 0x69,
	0x2f, 0x61, 0x6b, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_modules_manifest_proto_rawDescOnce sync.Once
	file_modules_manifest_proto_rawDescData = file_modules_manifest_proto_rawDesc
)

func file_modules_manifest_proto_rawDescGZIP() []byte {
	file_modules_manifest_proto_rawDescOnce.Do(func() {
		file_modules_manifest_proto_rawDescData = protoimpl.X.CompressGZIP(file_modules_manifest_proto_rawDescData)
	})
	return file_modules_manifest_proto_rawDescData
}

var file_modules_manifest_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_modules_manifest_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_modules_manifest_proto_goTypes = []any{
	(ManifestWebPathType)(0), // 0: modules.ManifestWebPathType
	(*ManifestWebPath)(nil),  // 1: modules.ManifestWebPath
	(*Manifest)(nil),         // 2: modules.Manifest
}
var file_modules_manifest_proto_depIdxs = []int32{
	0, // 0: modules.ManifestWebPath.type:type_name -> modules.ManifestWebPathType
	1, // 1: modules.Manifest.web_paths:type_name -> modules.ManifestWebPath
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_modules_manifest_proto_init() }
func file_modules_manifest_proto_init() {
	if File_modules_manifest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_modules_manifest_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ManifestWebPath); i {
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
		file_modules_manifest_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Manifest); i {
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
			RawDescriptor: file_modules_manifest_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_modules_manifest_proto_goTypes,
		DependencyIndexes: file_modules_manifest_proto_depIdxs,
		EnumInfos:         file_modules_manifest_proto_enumTypes,
		MessageInfos:      file_modules_manifest_proto_msgTypes,
	}.Build()
	File_modules_manifest_proto = out.File
	file_modules_manifest_proto_rawDesc = nil
	file_modules_manifest_proto_goTypes = nil
	file_modules_manifest_proto_depIdxs = nil
}
