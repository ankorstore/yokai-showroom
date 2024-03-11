// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.25.3
// source: proto/transform.proto

package proto

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

type Transformer int32

const (
	Transformer_TRANSFORMER_UNSPECIFIED Transformer = 0
	Transformer_TRANSFORMER_UPPERCASE   Transformer = 1
	Transformer_TRANSFORMER_LOWERCASE   Transformer = 2
)

// Enum value maps for Transformer.
var (
	Transformer_name = map[int32]string{
		0: "TRANSFORMER_UNSPECIFIED",
		1: "TRANSFORMER_UPPERCASE",
		2: "TRANSFORMER_LOWERCASE",
	}
	Transformer_value = map[string]int32{
		"TRANSFORMER_UNSPECIFIED": 0,
		"TRANSFORMER_UPPERCASE":   1,
		"TRANSFORMER_LOWERCASE":   2,
	}
)

func (x Transformer) Enum() *Transformer {
	p := new(Transformer)
	*p = x
	return p
}

func (x Transformer) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Transformer) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_transform_proto_enumTypes[0].Descriptor()
}

func (Transformer) Type() protoreflect.EnumType {
	return &file_proto_transform_proto_enumTypes[0]
}

func (x Transformer) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Transformer.Descriptor instead.
func (Transformer) EnumDescriptor() ([]byte, []int) {
	return file_proto_transform_proto_rawDescGZIP(), []int{0}
}

type TransformTextRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transformer Transformer `protobuf:"varint,1,opt,name=transformer,proto3,enum=transform.Transformer" json:"transformer,omitempty"`
	Text        string      `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *TransformTextRequest) Reset() {
	*x = TransformTextRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transform_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransformTextRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformTextRequest) ProtoMessage() {}

func (x *TransformTextRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transform_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransformTextRequest.ProtoReflect.Descriptor instead.
func (*TransformTextRequest) Descriptor() ([]byte, []int) {
	return file_proto_transform_proto_rawDescGZIP(), []int{0}
}

func (x *TransformTextRequest) GetTransformer() Transformer {
	if x != nil {
		return x.Transformer
	}
	return Transformer_TRANSFORMER_UNSPECIFIED
}

func (x *TransformTextRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type TransformTextResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *TransformTextResponse) Reset() {
	*x = TransformTextResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transform_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransformTextResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformTextResponse) ProtoMessage() {}

func (x *TransformTextResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transform_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransformTextResponse.ProtoReflect.Descriptor instead.
func (*TransformTextResponse) Descriptor() ([]byte, []int) {
	return file_proto_transform_proto_rawDescGZIP(), []int{1}
}

func (x *TransformTextResponse) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

var File_proto_transform_proto protoreflect.FileDescriptor

var file_proto_transform_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f,
	0x72, 0x6d, 0x22, 0x64, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54,
	0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x0b, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f,
	0x72, 0x6d, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x2b, 0x0a, 0x15, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x65, 0x78, 0x74, 0x2a, 0x60, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f,
	0x72, 0x6d, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x17, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x46, 0x4f, 0x52,
	0x4d, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x19, 0x0a, 0x15, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x46, 0x4f, 0x52, 0x4d, 0x45, 0x52,
	0x5f, 0x55, 0x50, 0x50, 0x45, 0x52, 0x43, 0x41, 0x53, 0x45, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15,
	0x54, 0x52, 0x41, 0x4e, 0x53, 0x46, 0x4f, 0x52, 0x4d, 0x45, 0x52, 0x5f, 0x4c, 0x4f, 0x57, 0x45,
	0x52, 0x43, 0x41, 0x53, 0x45, 0x10, 0x02, 0x32, 0xce, 0x01, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x54, 0x0a, 0x0d, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x78,
	0x74, 0x12, 0x1f, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x20, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x60, 0x0a, 0x15, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x6f, 0x72, 0x6d, 0x41, 0x6e, 0x64, 0x53, 0x70, 0x6c, 0x69, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12,
	0x1f, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x20, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x6b, 0x6f, 0x72, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2f, 0x79, 0x6f, 0x6b, 0x61, 0x69, 0x2d, 0x73, 0x68, 0x6f, 0x77, 0x72, 0x6f, 0x6f, 0x6d,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_transform_proto_rawDescOnce sync.Once
	file_proto_transform_proto_rawDescData = file_proto_transform_proto_rawDesc
)

func file_proto_transform_proto_rawDescGZIP() []byte {
	file_proto_transform_proto_rawDescOnce.Do(func() {
		file_proto_transform_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_transform_proto_rawDescData)
	})
	return file_proto_transform_proto_rawDescData
}

var file_proto_transform_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_transform_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_transform_proto_goTypes = []interface{}{
	(Transformer)(0),              // 0: transform.Transformer
	(*TransformTextRequest)(nil),  // 1: transform.TransformTextRequest
	(*TransformTextResponse)(nil), // 2: transform.TransformTextResponse
}
var file_proto_transform_proto_depIdxs = []int32{
	0, // 0: transform.TransformTextRequest.transformer:type_name -> transform.Transformer
	1, // 1: transform.TransformTextService.TransformText:input_type -> transform.TransformTextRequest
	1, // 2: transform.TransformTextService.TransformAndSplitText:input_type -> transform.TransformTextRequest
	2, // 3: transform.TransformTextService.TransformText:output_type -> transform.TransformTextResponse
	2, // 4: transform.TransformTextService.TransformAndSplitText:output_type -> transform.TransformTextResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_transform_proto_init() }
func file_proto_transform_proto_init() {
	if File_proto_transform_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_transform_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransformTextRequest); i {
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
		file_proto_transform_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransformTextResponse); i {
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
			RawDescriptor: file_proto_transform_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_transform_proto_goTypes,
		DependencyIndexes: file_proto_transform_proto_depIdxs,
		EnumInfos:         file_proto_transform_proto_enumTypes,
		MessageInfos:      file_proto_transform_proto_msgTypes,
	}.Build()
	File_proto_transform_proto = out.File
	file_proto_transform_proto_rawDesc = nil
	file_proto_transform_proto_goTypes = nil
	file_proto_transform_proto_depIdxs = nil
}
