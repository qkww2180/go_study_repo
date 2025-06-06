// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.3
// source: student_service.proto

package grpc_service

import (
	context "context"
	model "dqq/go/basic/grpc/idl/model"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type QueryStudentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	PageSize int32  `protobuf:"varint,3,opt,name=PageSize,proto3" json:"PageSize,omitempty"`
	PageNo   int32  `protobuf:"varint,4,opt,name=PageNo,proto3" json:"PageNo,omitempty"`
}

func (x *QueryStudentRequest) Reset() {
	*x = QueryStudentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_student_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryStudentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryStudentRequest) ProtoMessage() {}

func (x *QueryStudentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_student_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryStudentRequest.ProtoReflect.Descriptor instead.
func (*QueryStudentRequest) Descriptor() ([]byte, []int) {
	return file_student_service_proto_rawDescGZIP(), []int{0}
}

func (x *QueryStudentRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *QueryStudentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *QueryStudentRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *QueryStudentRequest) GetPageNo() int32 {
	if x != nil {
		return x.PageNo
	}
	return 0
}

type QueryStudentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Students []*model.Student `protobuf:"bytes,1,rep,name=Students,proto3" json:"Students,omitempty"`
}

func (x *QueryStudentResponse) Reset() {
	*x = QueryStudentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_student_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryStudentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryStudentResponse) ProtoMessage() {}

func (x *QueryStudentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_student_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryStudentResponse.ProtoReflect.Descriptor instead.
func (*QueryStudentResponse) Descriptor() ([]byte, []int) {
	return file_student_service_proto_rawDescGZIP(), []int{1}
}

func (x *QueryStudentResponse) GetStudents() []*model.Student {
	if x != nil {
		return x.Students
	}
	return nil
}

type StudentIds struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []int64 `protobuf:"varint,1,rep,packed,name=Ids,proto3" json:"Ids,omitempty"`
}

func (x *StudentIds) Reset() {
	*x = StudentIds{}
	if protoimpl.UnsafeEnabled {
		mi := &file_student_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StudentIds) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StudentIds) ProtoMessage() {}

func (x *StudentIds) ProtoReflect() protoreflect.Message {
	mi := &file_student_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StudentIds.ProtoReflect.Descriptor instead.
func (*StudentIds) Descriptor() ([]byte, []int) {
	return file_student_service_proto_rawDescGZIP(), []int{2}
}

func (x *StudentIds) GetIds() []int64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type StudentId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *StudentId) Reset() {
	*x = StudentId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_student_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StudentId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StudentId) ProtoMessage() {}

func (x *StudentId) ProtoReflect() protoreflect.Message {
	mi := &file_student_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StudentId.ProtoReflect.Descriptor instead.
func (*StudentId) Descriptor() ([]byte, []int) {
	return file_student_service_proto_rawDescGZIP(), []int{3}
}

func (x *StudentId) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_student_service_proto protoreflect.FileDescriptor

var file_student_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x64, 0x6c, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x1a, 0x13, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x73, 0x74, 0x75, 0x64,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x13, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x50, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x50, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x22, 0x46, 0x0a, 0x14, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2e, 0x0a, 0x08, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x73,
	0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73,
	0x22, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x12, 0x10,
	0x0a, 0x03, 0x49, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x03, 0x49, 0x64, 0x73,
	0x22, 0x1b, 0x0a, 0x09, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x32, 0xfe, 0x02,
	0x0a, 0x07, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x53, 0x0a, 0x0c, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x2e, 0x69, 0x64, 0x6c, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x69, 0x64,
	0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53,
	0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c,
	0x0a, 0x0e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x31,
	0x12, 0x17, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53,
	0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x1a, 0x21, 0x2e, 0x69, 0x64, 0x6c, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x32, 0x12, 0x17,
	0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x1a, 0x12, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x12, 0x4d, 0x0a,
	0x0e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x33, 0x12,
	0x16, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74,
	0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x1a, 0x21, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x12, 0x40, 0x0a, 0x0e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x34, 0x12, 0x16,
	0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x1a, 0x12, 0x2e, 0x69, 0x64, 0x6c, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x2e, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x28, 0x01, 0x30, 0x01, 0x42, 0x1a,
	0x5a, 0x18, 0x69, 0x64, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3b, 0x67, 0x72,
	0x70, 0x63, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_student_service_proto_rawDescOnce sync.Once
	file_student_service_proto_rawDescData = file_student_service_proto_rawDesc
)

func file_student_service_proto_rawDescGZIP() []byte {
	file_student_service_proto_rawDescOnce.Do(func() {
		file_student_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_student_service_proto_rawDescData)
	})
	return file_student_service_proto_rawDescData
}

var file_student_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_student_service_proto_goTypes = []interface{}{
	(*QueryStudentRequest)(nil),  // 0: idl.service.QueryStudentRequest
	(*QueryStudentResponse)(nil), // 1: idl.service.QueryStudentResponse
	(*StudentIds)(nil),           // 2: idl.service.StudentIds
	(*StudentId)(nil),            // 3: idl.service.StudentId
	(*model.Student)(nil),        // 4: idl.model.student
}
var file_student_service_proto_depIdxs = []int32{
	4, // 0: idl.service.QueryStudentResponse.Students:type_name -> idl.model.student
	0, // 1: idl.service.student.QueryStudent:input_type -> idl.service.QueryStudentRequest
	2, // 2: idl.service.student.QueryStudents1:input_type -> idl.service.StudentIds
	2, // 3: idl.service.student.QueryStudents2:input_type -> idl.service.StudentIds
	3, // 4: idl.service.student.QueryStudents3:input_type -> idl.service.StudentId
	3, // 5: idl.service.student.QueryStudents4:input_type -> idl.service.StudentId
	1, // 6: idl.service.student.QueryStudent:output_type -> idl.service.QueryStudentResponse
	1, // 7: idl.service.student.QueryStudents1:output_type -> idl.service.QueryStudentResponse
	4, // 8: idl.service.student.QueryStudents2:output_type -> idl.model.student
	1, // 9: idl.service.student.QueryStudents3:output_type -> idl.service.QueryStudentResponse
	4, // 10: idl.service.student.QueryStudents4:output_type -> idl.model.student
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_student_service_proto_init() }
func file_student_service_proto_init() {
	if File_student_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_student_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryStudentRequest); i {
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
		file_student_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryStudentResponse); i {
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
		file_student_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StudentIds); i {
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
		file_student_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StudentId); i {
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
			RawDescriptor: file_student_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_student_service_proto_goTypes,
		DependencyIndexes: file_student_service_proto_depIdxs,
		MessageInfos:      file_student_service_proto_msgTypes,
	}.Build()
	File_student_service_proto = out.File
	file_student_service_proto_rawDesc = nil
	file_student_service_proto_goTypes = nil
	file_student_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StudentClient is the client API for Student service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StudentClient interface {
	// Unary RPC
	QueryStudent(ctx context.Context, in *QueryStudentRequest, opts ...grpc.CallOption) (*QueryStudentResponse, error)
	QueryStudents1(ctx context.Context, in *StudentIds, opts ...grpc.CallOption) (*QueryStudentResponse, error)
	// Server streaming RPC
	QueryStudents2(ctx context.Context, in *StudentIds, opts ...grpc.CallOption) (Student_QueryStudents2Client, error)
	// Client streaming RPC
	QueryStudents3(ctx context.Context, opts ...grpc.CallOption) (Student_QueryStudents3Client, error)
	// Bidirectional streaming RPC
	QueryStudents4(ctx context.Context, opts ...grpc.CallOption) (Student_QueryStudents4Client, error)
}

type studentClient struct {
	cc grpc.ClientConnInterface
}

func NewStudentClient(cc grpc.ClientConnInterface) StudentClient {
	return &studentClient{cc}
}

func (c *studentClient) QueryStudent(ctx context.Context, in *QueryStudentRequest, opts ...grpc.CallOption) (*QueryStudentResponse, error) {
	out := new(QueryStudentResponse)
	err := c.cc.Invoke(ctx, "/idl.service.student/QueryStudent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentClient) QueryStudents1(ctx context.Context, in *StudentIds, opts ...grpc.CallOption) (*QueryStudentResponse, error) {
	out := new(QueryStudentResponse)
	err := c.cc.Invoke(ctx, "/idl.service.student/QueryStudents1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentClient) QueryStudents2(ctx context.Context, in *StudentIds, opts ...grpc.CallOption) (Student_QueryStudents2Client, error) {
	stream, err := c.cc.NewStream(ctx, &_Student_serviceDesc.Streams[0], "/idl.service.student/QueryStudents2", opts...)
	if err != nil {
		return nil, err
	}
	x := &studentQueryStudents2Client{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Student_QueryStudents2Client interface {
	Recv() (*model.Student, error)
	grpc.ClientStream
}

type studentQueryStudents2Client struct {
	grpc.ClientStream
}

func (x *studentQueryStudents2Client) Recv() (*model.Student, error) {
	m := new(model.Student)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *studentClient) QueryStudents3(ctx context.Context, opts ...grpc.CallOption) (Student_QueryStudents3Client, error) {
	stream, err := c.cc.NewStream(ctx, &_Student_serviceDesc.Streams[1], "/idl.service.student/QueryStudents3", opts...)
	if err != nil {
		return nil, err
	}
	x := &studentQueryStudents3Client{stream}
	return x, nil
}

type Student_QueryStudents3Client interface {
	Send(*StudentId) error
	CloseAndRecv() (*QueryStudentResponse, error)
	grpc.ClientStream
}

type studentQueryStudents3Client struct {
	grpc.ClientStream
}

func (x *studentQueryStudents3Client) Send(m *StudentId) error {
	return x.ClientStream.SendMsg(m)
}

func (x *studentQueryStudents3Client) CloseAndRecv() (*QueryStudentResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(QueryStudentResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *studentClient) QueryStudents4(ctx context.Context, opts ...grpc.CallOption) (Student_QueryStudents4Client, error) {
	stream, err := c.cc.NewStream(ctx, &_Student_serviceDesc.Streams[2], "/idl.service.student/QueryStudents4", opts...)
	if err != nil {
		return nil, err
	}
	x := &studentQueryStudents4Client{stream}
	return x, nil
}

type Student_QueryStudents4Client interface {
	Send(*StudentId) error
	Recv() (*model.Student, error)
	grpc.ClientStream
}

type studentQueryStudents4Client struct {
	grpc.ClientStream
}

func (x *studentQueryStudents4Client) Send(m *StudentId) error {
	return x.ClientStream.SendMsg(m)
}

func (x *studentQueryStudents4Client) Recv() (*model.Student, error) {
	m := new(model.Student)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StudentServer is the server API for Student service.
type StudentServer interface {
	// Unary RPC
	QueryStudent(context.Context, *QueryStudentRequest) (*QueryStudentResponse, error)
	QueryStudents1(context.Context, *StudentIds) (*QueryStudentResponse, error)
	// Server streaming RPC
	QueryStudents2(*StudentIds, Student_QueryStudents2Server) error
	// Client streaming RPC
	QueryStudents3(Student_QueryStudents3Server) error
	// Bidirectional streaming RPC
	QueryStudents4(Student_QueryStudents4Server) error
}

// UnimplementedStudentServer can be embedded to have forward compatible implementations.
type UnimplementedStudentServer struct {
}

func (*UnimplementedStudentServer) QueryStudent(context.Context, *QueryStudentRequest) (*QueryStudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryStudent not implemented")
}
func (*UnimplementedStudentServer) QueryStudents1(context.Context, *StudentIds) (*QueryStudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryStudents1 not implemented")
}
func (*UnimplementedStudentServer) QueryStudents2(*StudentIds, Student_QueryStudents2Server) error {
	return status.Errorf(codes.Unimplemented, "method QueryStudents2 not implemented")
}
func (*UnimplementedStudentServer) QueryStudents3(Student_QueryStudents3Server) error {
	return status.Errorf(codes.Unimplemented, "method QueryStudents3 not implemented")
}
func (*UnimplementedStudentServer) QueryStudents4(Student_QueryStudents4Server) error {
	return status.Errorf(codes.Unimplemented, "method QueryStudents4 not implemented")
}

func RegisterStudentServer(s *grpc.Server, srv StudentServer) {
	s.RegisterService(&_Student_serviceDesc, srv)
}

func _Student_QueryStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStudentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServer).QueryStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/idl.service.student/QueryStudent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServer).QueryStudent(ctx, req.(*QueryStudentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Student_QueryStudents1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StudentIds)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentServer).QueryStudents1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/idl.service.student/QueryStudents1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentServer).QueryStudents1(ctx, req.(*StudentIds))
	}
	return interceptor(ctx, in, info, handler)
}

func _Student_QueryStudents2_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StudentIds)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StudentServer).QueryStudents2(m, &studentQueryStudents2Server{stream})
}

type Student_QueryStudents2Server interface {
	Send(*model.Student) error
	grpc.ServerStream
}

type studentQueryStudents2Server struct {
	grpc.ServerStream
}

func (x *studentQueryStudents2Server) Send(m *model.Student) error {
	return x.ServerStream.SendMsg(m)
}

func _Student_QueryStudents3_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StudentServer).QueryStudents3(&studentQueryStudents3Server{stream})
}

type Student_QueryStudents3Server interface {
	SendAndClose(*QueryStudentResponse) error
	Recv() (*StudentId, error)
	grpc.ServerStream
}

type studentQueryStudents3Server struct {
	grpc.ServerStream
}

func (x *studentQueryStudents3Server) SendAndClose(m *QueryStudentResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *studentQueryStudents3Server) Recv() (*StudentId, error) {
	m := new(StudentId)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Student_QueryStudents4_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StudentServer).QueryStudents4(&studentQueryStudents4Server{stream})
}

type Student_QueryStudents4Server interface {
	Send(*model.Student) error
	Recv() (*StudentId, error)
	grpc.ServerStream
}

type studentQueryStudents4Server struct {
	grpc.ServerStream
}

func (x *studentQueryStudents4Server) Send(m *model.Student) error {
	return x.ServerStream.SendMsg(m)
}

func (x *studentQueryStudents4Server) Recv() (*StudentId, error) {
	m := new(StudentId)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Student_serviceDesc = grpc.ServiceDesc{
	ServiceName: "idl.service.student",
	HandlerType: (*StudentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryStudent",
			Handler:    _Student_QueryStudent_Handler,
		},
		{
			MethodName: "QueryStudents1",
			Handler:    _Student_QueryStudents1_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "QueryStudents2",
			Handler:       _Student_QueryStudents2_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "QueryStudents3",
			Handler:       _Student_QueryStudents3_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "QueryStudents4",
			Handler:       _Student_QueryStudents4_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "student_service.proto",
}
