// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: trade_bot_protobuf/src/entiretick.proto

package entiretick

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

type EntireTickArrProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*EntireTickProto `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *EntireTickArrProto) Reset() {
	*x = EntireTickArrProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_bot_protobuf_src_entiretick_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntireTickArrProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntireTickArrProto) ProtoMessage() {}

func (x *EntireTickArrProto) ProtoReflect() protoreflect.Message {
	mi := &file_trade_bot_protobuf_src_entiretick_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntireTickArrProto.ProtoReflect.Descriptor instead.
func (*EntireTickArrProto) Descriptor() ([]byte, []int) {
	return file_trade_bot_protobuf_src_entiretick_proto_rawDescGZIP(), []int{0}
}

func (x *EntireTickArrProto) GetData() []*EntireTickProto {
	if x != nil {
		return x.Data
	}
	return nil
}

type EntireTickProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts        int64   `protobuf:"varint,1,opt,name=ts,proto3" json:"ts,omitempty"`
	Close     float64 `protobuf:"fixed64,2,opt,name=close,proto3" json:"close,omitempty"`
	Volume    int64   `protobuf:"varint,3,opt,name=volume,proto3" json:"volume,omitempty"`
	BidPrice  float64 `protobuf:"fixed64,4,opt,name=bid_price,json=bidPrice,proto3" json:"bid_price,omitempty"`
	BidVolume int64   `protobuf:"varint,5,opt,name=bid_volume,json=bidVolume,proto3" json:"bid_volume,omitempty"`
	AskPrice  float64 `protobuf:"fixed64,6,opt,name=ask_price,json=askPrice,proto3" json:"ask_price,omitempty"`
	AskVolume int64   `protobuf:"varint,7,opt,name=ask_volume,json=askVolume,proto3" json:"ask_volume,omitempty"`
	TickType  int64   `protobuf:"varint,8,opt,name=tick_type,json=tickType,proto3" json:"tick_type,omitempty"`
}

func (x *EntireTickProto) Reset() {
	*x = EntireTickProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trade_bot_protobuf_src_entiretick_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntireTickProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntireTickProto) ProtoMessage() {}

func (x *EntireTickProto) ProtoReflect() protoreflect.Message {
	mi := &file_trade_bot_protobuf_src_entiretick_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntireTickProto.ProtoReflect.Descriptor instead.
func (*EntireTickProto) Descriptor() ([]byte, []int) {
	return file_trade_bot_protobuf_src_entiretick_proto_rawDescGZIP(), []int{1}
}

func (x *EntireTickProto) GetTs() int64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *EntireTickProto) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *EntireTickProto) GetVolume() int64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

func (x *EntireTickProto) GetBidPrice() float64 {
	if x != nil {
		return x.BidPrice
	}
	return 0
}

func (x *EntireTickProto) GetBidVolume() int64 {
	if x != nil {
		return x.BidVolume
	}
	return 0
}

func (x *EntireTickProto) GetAskPrice() float64 {
	if x != nil {
		return x.AskPrice
	}
	return 0
}

func (x *EntireTickProto) GetAskVolume() int64 {
	if x != nil {
		return x.AskVolume
	}
	return 0
}

func (x *EntireTickProto) GetTickType() int64 {
	if x != nil {
		return x.TickType
	}
	return 0
}

var File_trade_bot_protobuf_src_entiretick_proto protoreflect.FileDescriptor

var file_trade_bot_protobuf_src_entiretick_proto_rawDesc = []byte{
	0x0a, 0x27, 0x74, 0x72, 0x61, 0x64, 0x65, 0x5f, 0x62, 0x6f, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x74,
	0x69, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x65, 0x6e, 0x74, 0x69, 0x72,
	0x65, 0x74, 0x69, 0x63, 0x6b, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x4e,
	0x0a, 0x12, 0x45, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x41, 0x72, 0x72, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x38, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x24, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x74, 0x69, 0x63, 0x6b, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x54,
	0x69, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xe4,
	0x01, 0x0a, 0x0f, 0x45, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x62, 0x69, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x08, 0x62, 0x69, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x62, 0x69, 0x64, 0x5f, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x62, 0x69, 0x64, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x73, 0x6b, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x08, 0x61, 0x73, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x73, 0x6b,
	0x5f, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x61,
	0x73, 0x6b, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x69, 0x63, 0x6b,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x74, 0x69, 0x63,
	0x6b, 0x54, 0x79, 0x70, 0x65, 0x42, 0x17, 0x5a, 0x15, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x74, 0x69, 0x63, 0x6b, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trade_bot_protobuf_src_entiretick_proto_rawDescOnce sync.Once
	file_trade_bot_protobuf_src_entiretick_proto_rawDescData = file_trade_bot_protobuf_src_entiretick_proto_rawDesc
)

func file_trade_bot_protobuf_src_entiretick_proto_rawDescGZIP() []byte {
	file_trade_bot_protobuf_src_entiretick_proto_rawDescOnce.Do(func() {
		file_trade_bot_protobuf_src_entiretick_proto_rawDescData = protoimpl.X.CompressGZIP(file_trade_bot_protobuf_src_entiretick_proto_rawDescData)
	})
	return file_trade_bot_protobuf_src_entiretick_proto_rawDescData
}

var file_trade_bot_protobuf_src_entiretick_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_trade_bot_protobuf_src_entiretick_proto_goTypes = []interface{}{
	(*EntireTickArrProto)(nil), // 0: entiretick_protobuf.EntireTickArrProto
	(*EntireTickProto)(nil),    // 1: entiretick_protobuf.EntireTickProto
}
var file_trade_bot_protobuf_src_entiretick_proto_depIdxs = []int32{
	1, // 0: entiretick_protobuf.EntireTickArrProto.data:type_name -> entiretick_protobuf.EntireTickProto
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_trade_bot_protobuf_src_entiretick_proto_init() }
func file_trade_bot_protobuf_src_entiretick_proto_init() {
	if File_trade_bot_protobuf_src_entiretick_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trade_bot_protobuf_src_entiretick_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntireTickArrProto); i {
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
		file_trade_bot_protobuf_src_entiretick_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntireTickProto); i {
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
			RawDescriptor: file_trade_bot_protobuf_src_entiretick_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trade_bot_protobuf_src_entiretick_proto_goTypes,
		DependencyIndexes: file_trade_bot_protobuf_src_entiretick_proto_depIdxs,
		MessageInfos:      file_trade_bot_protobuf_src_entiretick_proto_msgTypes,
	}.Build()
	File_trade_bot_protobuf_src_entiretick_proto = out.File
	file_trade_bot_protobuf_src_entiretick_proto_rawDesc = nil
	file_trade_bot_protobuf_src_entiretick_proto_goTypes = nil
	file_trade_bot_protobuf_src_entiretick_proto_depIdxs = nil
}
