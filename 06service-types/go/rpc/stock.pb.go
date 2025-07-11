// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.2
// source: stock.proto

package stv1

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

// 股票订阅请求
type StockSubscribeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Symbols       []string               `protobuf:"bytes,1,rep,name=symbols,proto3" json:"symbols,omitempty"`                   // 股票代码列表，如 ["AAPL", "GOOGL", "TSLA"]
	ClientId      string                 `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"` // 客户端标识
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StockSubscribeRequest) Reset() {
	*x = StockSubscribeRequest{}
	mi := &file_stock_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StockSubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockSubscribeRequest) ProtoMessage() {}

func (x *StockSubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockSubscribeRequest.ProtoReflect.Descriptor instead.
func (*StockSubscribeRequest) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{0}
}

func (x *StockSubscribeRequest) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

func (x *StockSubscribeRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

// 股票价格更新消息
type StockPriceUpdate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Symbol        string                 `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`                                      // 股票代码
	CurrentPrice  float64                `protobuf:"fixed64,2,opt,name=current_price,json=currentPrice,proto3" json:"current_price,omitempty"`    // 当前价格
	ChangeAmount  float64                `protobuf:"fixed64,3,opt,name=change_amount,json=changeAmount,proto3" json:"change_amount,omitempty"`    // 变化金额
	ChangePercent float64                `protobuf:"fixed64,4,opt,name=change_percent,json=changePercent,proto3" json:"change_percent,omitempty"` // 变化百分比
	Timestamp     int64                  `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                               // 时间戳（Unix时间）
	Volume        int64                  `protobuf:"varint,6,opt,name=volume,proto3" json:"volume,omitempty"`                                     // 成交量
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StockPriceUpdate) Reset() {
	*x = StockPriceUpdate{}
	mi := &file_stock_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StockPriceUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockPriceUpdate) ProtoMessage() {}

func (x *StockPriceUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockPriceUpdate.ProtoReflect.Descriptor instead.
func (*StockPriceUpdate) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{1}
}

func (x *StockPriceUpdate) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *StockPriceUpdate) GetCurrentPrice() float64 {
	if x != nil {
		return x.CurrentPrice
	}
	return 0
}

func (x *StockPriceUpdate) GetChangeAmount() float64 {
	if x != nil {
		return x.ChangeAmount
	}
	return 0
}

func (x *StockPriceUpdate) GetChangePercent() float64 {
	if x != nil {
		return x.ChangePercent
	}
	return 0
}

func (x *StockPriceUpdate) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *StockPriceUpdate) GetVolume() int64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

var File_stock_proto protoreflect.FileDescriptor

var file_stock_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73,
	0x74, 0x6f, 0x63, 0x6b, 0x22, 0x4e, 0x0a, 0x15, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x22, 0xd1, 0x01, 0x0a, 0x10, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x72, 0x69,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x65, 0x72, 0x63, 0x65,
	0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x32, 0x60, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x63,
	0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x50, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12,
	0x1c, 0x2e, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x73, 0x74, 0x6f, 0x63, 0x6b, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x69, 0x6e, 0x32, 0x31, 0x31,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2d, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x3b, 0x73, 0x74, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stock_proto_rawDescOnce sync.Once
	file_stock_proto_rawDescData = file_stock_proto_rawDesc
)

func file_stock_proto_rawDescGZIP() []byte {
	file_stock_proto_rawDescOnce.Do(func() {
		file_stock_proto_rawDescData = protoimpl.X.CompressGZIP(file_stock_proto_rawDescData)
	})
	return file_stock_proto_rawDescData
}

var file_stock_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_stock_proto_goTypes = []any{
	(*StockSubscribeRequest)(nil), // 0: stock.StockSubscribeRequest
	(*StockPriceUpdate)(nil),      // 1: stock.StockPriceUpdate
}
var file_stock_proto_depIdxs = []int32{
	0, // 0: stock.StockService.SubscribeStockPrice:input_type -> stock.StockSubscribeRequest
	1, // 1: stock.StockService.SubscribeStockPrice:output_type -> stock.StockPriceUpdate
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_stock_proto_init() }
func file_stock_proto_init() {
	if File_stock_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stock_proto_goTypes,
		DependencyIndexes: file_stock_proto_depIdxs,
		MessageInfos:      file_stock_proto_msgTypes,
	}.Build()
	File_stock_proto = out.File
	file_stock_proto_rawDesc = nil
	file_stock_proto_goTypes = nil
	file_stock_proto_depIdxs = nil
}
