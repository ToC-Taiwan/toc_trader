// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: trade_bot_protobuf/src/kbar.proto

package kbar

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

type KbarArrProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*KbarProto `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *KbarArrProto) Reset() {
	*x = KbarArrProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_bot_protobuf_src_kbar_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KbarArrProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KbarArrProto) ProtoMessage() {}

func (x *KbarArrProto) ProtoReflect() protoreflect.Message {
	mi := &file_trade_bot_protobuf_src_kbar_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KbarArrProto.ProtoReflect.Descriptor instead.
func (*KbarArrProto) Descriptor() ([]byte, []int) {
	return file_trade_bot_protobuf_src_kbar_proto_rawDescGZIP(), []int{0}
}

func (x *KbarArrProto) GetData() []*KbarProto {
	if x != nil {
		return x.Data
	}
	return nil
}

type KbarProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts     int64   `protobuf:"varint,1,opt,name=ts,proto3" json:"ts,omitempty"`
	Close  float64 `protobuf:"fixed64,2,opt,name=Close,proto3" json:"Close,omitempty"`
	Open   float64 `protobuf:"fixed64,3,opt,name=Open,proto3" json:"Open,omitempty"`
	High   float64 `protobuf:"fixed64,4,opt,name=High,proto3" json:"High,omitempty"`
	Low    float64 `protobuf:"fixed64,5,opt,name=Low,proto3" json:"Low,omitempty"`
	Volume int64   `protobuf:"varint,6,opt,name=Volume,proto3" json:"Volume,omitempty"`
}

func (x *KbarProto) Reset() {
	*x = KbarProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_bot_protobuf_src_kbar_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KbarProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KbarProto) ProtoMessage() {}

func (x *KbarProto) ProtoReflect() protoreflect.Message {
	mi := &file_trade_bot_protobuf_src_kbar_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KbarProto.ProtoReflect.Descriptor instead.
func (*KbarProto) Descriptor() ([]byte, []int) {
	return file_trade_bot_protobuf_src_kbar_proto_rawDescGZIP(), []int{1}
}

func (x *KbarProto) GetTs() int64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *KbarProto) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *KbarProto) GetOpen() float64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *KbarProto) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *KbarProto) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *KbarProto) GetVolume() int64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

var File_trade_bot_protobuf_src_kbar_proto protoreflect.FileDescriptor

var file_trade_bot_protobuf_src_kbar_proto_rawDesc = []byte{
	0x0a, 0x21, 0x74, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x62, 0x6f, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x6b, 0x62, 0x61, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x6b, 0x62, 0x61, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x22, 0x3c, 0x0a, 0x0c, 0x4b, 0x62, 0x61, 0x72, 0x41, 0x72, 0x72, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x2c, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x6b, 0x62, 0x61, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4b, 0x62, 0x61, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x83, 0x01, 0x0a, 0x09, 0x4b, 0x62, 0x61, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x0a, 0x02, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x74, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4f, 0x70, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x04, 0x4f, 0x70, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x69, 0x67, 0x68,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x48, 0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03,
	0x4c, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x4c, 0x6f, 0x77, 0x12, 0x16,
	0x0a, 0x06, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2f, 0x6b, 0x62, 0x61, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_trade_bot_protobuf_src_kbar_proto_rawDescOnce sync.Once
	file_trade_bot_protobuf_src_kbar_proto_rawDescData = file_trade_bot_protobuf_src_kbar_proto_rawDesc
)

func file_trade_bot_protobuf_src_kbar_proto_rawDescGZIP() []byte {
	file_trade_bot_protobuf_src_kbar_proto_rawDescOnce.Do(func() {
		file_trade_bot_protobuf_src_kbar_proto_rawDescData = protoimpl.X.CompressGZIP(file_trade_bot_protobuf_src_kbar_proto_rawDescData)
	})
	return file_trade_bot_protobuf_src_kbar_proto_rawDescData
}

var file_trade_bot_protobuf_src_kbar_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_trade_bot_protobuf_src_kbar_proto_goTypes = []interface{}{
	(*KbarArrProto)(nil), // 0: kbar_protobuf.KbarArrProto
	(*KbarProto)(nil),    // 1: kbar_protobuf.KbarProto
}
var file_trade_bot_protobuf_src_kbar_proto_depIdxs = []int32{
	1, // 0: kbar_protobuf.KbarArrProto.data:type_name -> kbar_protobuf.KbarProto
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_trade_bot_protobuf_src_kbar_proto_init() }
func file_trade_bot_protobuf_src_kbar_proto_init() {
	if File_trade_bot_protobuf_src_kbar_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trade_bot_protobuf_src_kbar_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KbarArrProto); i {
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
		file_trade_bot_protobuf_src_kbar_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KbarProto); i {
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
			RawDescriptor: file_trade_bot_protobuf_src_kbar_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trade_bot_protobuf_src_kbar_proto_goTypes,
		DependencyIndexes: file_trade_bot_protobuf_src_kbar_proto_depIdxs,
		MessageInfos:      file_trade_bot_protobuf_src_kbar_proto_msgTypes,
	}.Build()
	File_trade_bot_protobuf_src_kbar_proto = out.File
	file_trade_bot_protobuf_src_kbar_proto_rawDesc = nil
	file_trade_bot_protobuf_src_kbar_proto_goTypes = nil
	file_trade_bot_protobuf_src_kbar_proto_depIdxs = nil
}
