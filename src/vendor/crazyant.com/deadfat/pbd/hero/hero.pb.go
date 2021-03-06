// Code generated by protoc-gen-go.
// source: hero.proto
// DO NOT EDIT!

/*
Package hero is a generated protocol buffer package.

It is generated from these files:
	hero.proto

It has these top-level messages:
	Empty
	Error
	Login
	GuestRegister
	LoginResp
	VersionUpdateAlert
	GetAchievement
	SetAchievement
	GetAchievementResp
	GetArenaLearderboardRankResp
	GetLearderboardRange
	LearderboardInfo
	GetLearderboardRangeResp
	UpdateLearderBoardScore
	UpdateLearderBoardScoreResp
	UnLockAchievement
	UnLockAchievementResp
	GetUnLockAchieveDate
	GetUnLockAchieveDateResp
*/
package hero

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// 指令号定义
type CG int32

const (
	CG_ID_ZERO_CG                CG = 0
	CG_ID_GUEST_LOGIN            CG = 2001
	CG_ID_GUEST_REGISTER         CG = 2002
	CG_ID_GET_ACHIEVEMENT        CG = 2003
	CG_ID_SET_ACHIEVEMENT        CG = 2004
	CG_ID_GET_ARENA_RANK         CG = 2005
	CG_ID_GET_LEARDERBOARD_RANGE CG = 2006
	CG_ID_UPDATE_LEARD_SCORE     CG = 2007
	CG_ID_UNLOCK_ACHIEVEMENT     CG = 2008
	CG_ID_GET_UNLOCK_ACHIEVEDATE CG = 2009
)

var CG_name = map[int32]string{
	0:    "ID_ZERO_CG",
	2001: "ID_GUEST_LOGIN",
	2002: "ID_GUEST_REGISTER",
	2003: "ID_GET_ACHIEVEMENT",
	2004: "ID_SET_ACHIEVEMENT",
	2005: "ID_GET_ARENA_RANK",
	2006: "ID_GET_LEARDERBOARD_RANGE",
	2007: "ID_UPDATE_LEARD_SCORE",
	2008: "ID_UNLOCK_ACHIEVEMENT",
	2009: "ID_GET_UNLOCK_ACHIEVEDATE",
}
var CG_value = map[string]int32{
	"ID_ZERO_CG":                0,
	"ID_GUEST_LOGIN":            2001,
	"ID_GUEST_REGISTER":         2002,
	"ID_GET_ACHIEVEMENT":        2003,
	"ID_SET_ACHIEVEMENT":        2004,
	"ID_GET_ARENA_RANK":         2005,
	"ID_GET_LEARDERBOARD_RANGE": 2006,
	"ID_UPDATE_LEARD_SCORE":     2007,
	"ID_UNLOCK_ACHIEVEMENT":     2008,
	"ID_GET_UNLOCK_ACHIEVEDATE": 2009,
}

func (x CG) String() string {
	return proto.EnumName(CG_name, int32(x))
}
func (CG) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// 协议号定义
type GC int32

const (
	GC_ID_ZERO_GC                     GC = 0
	GC_ID_GUEST_LOGIN_RESP            GC = 2998
	GC_ID_GET_ACHIEVEMENT_RESP        GC = 2997
	GC_ID_SET_ACHIEVEMENT_RESP        GC = 2996
	GC_ID_GET_ARENA_RANK_RESP         GC = 2995
	GC_ID_GET_LEARDERBOARD_RANGE_RESP GC = 2994
	GC_ID_UPDATE_LEARD_SCOREE_RESP    GC = 2993
	GC_ID_UNLOCK_ACHIEVEMENT_RESP     GC = 2992
	GC_ID_GET_UNLOCK_ACHIEVEDATE_RESP GC = 2991
)

var GC_name = map[int32]string{
	0:    "ID_ZERO_GC",
	2998: "ID_GUEST_LOGIN_RESP",
	2997: "ID_GET_ACHIEVEMENT_RESP",
	2996: "ID_SET_ACHIEVEMENT_RESP",
	2995: "ID_GET_ARENA_RANK_RESP",
	2994: "ID_GET_LEARDERBOARD_RANGE_RESP",
	2993: "ID_UPDATE_LEARD_SCOREE_RESP",
	2992: "ID_UNLOCK_ACHIEVEMENT_RESP",
	2991: "ID_GET_UNLOCK_ACHIEVEDATE_RESP",
}
var GC_value = map[string]int32{
	"ID_ZERO_GC":                     0,
	"ID_GUEST_LOGIN_RESP":            2998,
	"ID_GET_ACHIEVEMENT_RESP":        2997,
	"ID_SET_ACHIEVEMENT_RESP":        2996,
	"ID_GET_ARENA_RANK_RESP":         2995,
	"ID_GET_LEARDERBOARD_RANGE_RESP": 2994,
	"ID_UPDATE_LEARD_SCOREE_RESP":    2993,
	"ID_UNLOCK_ACHIEVEMENT_RESP":     2992,
	"ID_GET_UNLOCK_ACHIEVEDATE_RESP": 2991,
}

func (x GC) String() string {
	return proto.EnumName(GC_name, int32(x))
}
func (GC) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Error struct {
	Code int32    `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Args []string `protobuf:"bytes,2,rep,name=args" json:"args,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Login struct {
	Imei          string `protobuf:"bytes,1,opt,name=imei" json:"imei,omitempty"`
	ClientVersion string `protobuf:"bytes,2,opt,name=client_version,json=clientVersion" json:"client_version,omitempty"`
	Channel       string `protobuf:"bytes,3,opt,name=channel" json:"channel,omitempty"`
	FysdkToken    string `protobuf:"bytes,4,opt,name=fysdk_token,json=fysdkToken" json:"fysdk_token,omitempty"`
	FysdkUuid     string `protobuf:"bytes,5,opt,name=fysdk_uuid,json=fysdkUuid" json:"fysdk_uuid,omitempty"`
	NickName      string `protobuf:"bytes,6,opt,name=nickName" json:"nickName,omitempty"`
}

func (m *Login) Reset()                    { *m = Login{} }
func (m *Login) String() string            { return proto.CompactTextString(m) }
func (*Login) ProtoMessage()               {}
func (*Login) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type GuestRegister struct {
	Account string `protobuf:"bytes,1,opt,name=account" json:"account,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *GuestRegister) Reset()                    { *m = GuestRegister{} }
func (m *GuestRegister) String() string            { return proto.CompactTextString(m) }
func (*GuestRegister) ProtoMessage()               {}
func (*GuestRegister) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type LoginResp struct {
	Err            *Error              `protobuf:"bytes,1,opt,name=err" json:"err,omitempty"`
	Uid            uint32              `protobuf:"varint,2,opt,name=uid" json:"uid,omitempty"`
	Token          []byte              `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	ArenaScore     uint32              `protobuf:"varint,4,opt,name=arenaScore" json:"arenaScore,omitempty"`
	ArenaRank      uint32              `protobuf:"varint,5,opt,name=arenaRank" json:"arenaRank,omitempty"`
	ChallengeScore uint64              `protobuf:"varint,6,opt,name=challengeScore" json:"challengeScore,omitempty"`
	ChallengeRank  uint32              `protobuf:"varint,7,opt,name=challengeRank" json:"challengeRank,omitempty"`
	Update         *VersionUpdateAlert `protobuf:"bytes,8,opt,name=update" json:"update,omitempty"`
}

func (m *LoginResp) Reset()                    { *m = LoginResp{} }
func (m *LoginResp) String() string            { return proto.CompactTextString(m) }
func (*LoginResp) ProtoMessage()               {}
func (*LoginResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *LoginResp) GetErr() *Error {
	if m != nil {
		return m.Err
	}
	return nil
}

func (m *LoginResp) GetUpdate() *VersionUpdateAlert {
	if m != nil {
		return m.Update
	}
	return nil
}

// 版本更新提醒
type VersionUpdateAlert struct {
	Force   bool   `protobuf:"varint,1,opt,name=force" json:"force,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	Link    string `protobuf:"bytes,3,opt,name=link" json:"link,omitempty"`
}

func (m *VersionUpdateAlert) Reset()                    { *m = VersionUpdateAlert{} }
func (m *VersionUpdateAlert) String() string            { return proto.CompactTextString(m) }
func (*VersionUpdateAlert) ProtoMessage()               {}
func (*VersionUpdateAlert) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type GetAchievement struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *GetAchievement) Reset()                    { *m = GetAchievement{} }
func (m *GetAchievement) String() string            { return proto.CompactTextString(m) }
func (*GetAchievement) ProtoMessage()               {}
func (*GetAchievement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type SetAchievement struct {
	Name     string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	NowValue int32  `protobuf:"varint,2,opt,name=nowValue" json:"nowValue,omitempty"`
}

func (m *SetAchievement) Reset()                    { *m = SetAchievement{} }
func (m *SetAchievement) String() string            { return proto.CompactTextString(m) }
func (*SetAchievement) ProtoMessage()               {}
func (*SetAchievement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type GetAchievementResp struct {
	Value int32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *GetAchievementResp) Reset()                    { *m = GetAchievementResp{} }
func (m *GetAchievementResp) String() string            { return proto.CompactTextString(m) }
func (*GetAchievementResp) ProtoMessage()               {}
func (*GetAchievementResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type GetArenaLearderboardRankResp struct {
	Rank uint32 `protobuf:"varint,1,opt,name=rank" json:"rank,omitempty"`
}

func (m *GetArenaLearderboardRankResp) Reset()                    { *m = GetArenaLearderboardRankResp{} }
func (m *GetArenaLearderboardRankResp) String() string            { return proto.CompactTextString(m) }
func (*GetArenaLearderboardRankResp) ProtoMessage()               {}
func (*GetArenaLearderboardRankResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type GetLearderboardRange struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Start uint32 `protobuf:"varint,2,opt,name=start" json:"start,omitempty"`
	End   uint32 `protobuf:"varint,3,opt,name=end" json:"end,omitempty"`
}

func (m *GetLearderboardRange) Reset()                    { *m = GetLearderboardRange{} }
func (m *GetLearderboardRange) String() string            { return proto.CompactTextString(m) }
func (*GetLearderboardRange) ProtoMessage()               {}
func (*GetLearderboardRange) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type LearderboardInfo struct {
	Score int32  `protobuf:"varint,1,opt,name=score" json:"score,omitempty"`
	Rank  uint32 `protobuf:"varint,2,opt,name=rank" json:"rank,omitempty"`
	Name  string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *LearderboardInfo) Reset()                    { *m = LearderboardInfo{} }
func (m *LearderboardInfo) String() string            { return proto.CompactTextString(m) }
func (*LearderboardInfo) ProtoMessage()               {}
func (*LearderboardInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type GetLearderboardRangeResp struct {
	Count uint32              `protobuf:"varint,1,opt,name=count" json:"count,omitempty"`
	Infos []*LearderboardInfo `protobuf:"bytes,2,rep,name=infos" json:"infos,omitempty"`
}

func (m *GetLearderboardRangeResp) Reset()                    { *m = GetLearderboardRangeResp{} }
func (m *GetLearderboardRangeResp) String() string            { return proto.CompactTextString(m) }
func (*GetLearderboardRangeResp) ProtoMessage()               {}
func (*GetLearderboardRangeResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *GetLearderboardRangeResp) GetInfos() []*LearderboardInfo {
	if m != nil {
		return m.Infos
	}
	return nil
}

type UpdateLearderBoardScore struct {
	Name  string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Score uint64 `protobuf:"varint,1,opt,name=score" json:"score,omitempty"`
}

func (m *UpdateLearderBoardScore) Reset()                    { *m = UpdateLearderBoardScore{} }
func (m *UpdateLearderBoardScore) String() string            { return proto.CompactTextString(m) }
func (*UpdateLearderBoardScore) ProtoMessage()               {}
func (*UpdateLearderBoardScore) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type UpdateLearderBoardScoreResp struct {
	NewRank uint32 `protobuf:"varint,1,opt,name=newRank" json:"newRank,omitempty"`
	Success bool   `protobuf:"varint,2,opt,name=success" json:"success,omitempty"`
}

func (m *UpdateLearderBoardScoreResp) Reset()                    { *m = UpdateLearderBoardScoreResp{} }
func (m *UpdateLearderBoardScoreResp) String() string            { return proto.CompactTextString(m) }
func (*UpdateLearderBoardScoreResp) ProtoMessage()               {}
func (*UpdateLearderBoardScoreResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

type UnLockAchievement struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *UnLockAchievement) Reset()                    { *m = UnLockAchievement{} }
func (m *UnLockAchievement) String() string            { return proto.CompactTextString(m) }
func (*UnLockAchievement) ProtoMessage()               {}
func (*UnLockAchievement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type UnLockAchievementResp struct {
	Date string `protobuf:"bytes,1,opt,name=date" json:"date,omitempty"`
}

func (m *UnLockAchievementResp) Reset()                    { *m = UnLockAchievementResp{} }
func (m *UnLockAchievementResp) String() string            { return proto.CompactTextString(m) }
func (*UnLockAchievementResp) ProtoMessage()               {}
func (*UnLockAchievementResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

type GetUnLockAchieveDate struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *GetUnLockAchieveDate) Reset()                    { *m = GetUnLockAchieveDate{} }
func (m *GetUnLockAchieveDate) String() string            { return proto.CompactTextString(m) }
func (*GetUnLockAchieveDate) ProtoMessage()               {}
func (*GetUnLockAchieveDate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

type GetUnLockAchieveDateResp struct {
	Date string `protobuf:"bytes,1,opt,name=date" json:"date,omitempty"`
}

func (m *GetUnLockAchieveDateResp) Reset()                    { *m = GetUnLockAchieveDateResp{} }
func (m *GetUnLockAchieveDateResp) String() string            { return proto.CompactTextString(m) }
func (*GetUnLockAchieveDateResp) ProtoMessage()               {}
func (*GetUnLockAchieveDateResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func init() {
	proto.RegisterType((*Empty)(nil), "hero.Empty")
	proto.RegisterType((*Error)(nil), "hero.Error")
	proto.RegisterType((*Login)(nil), "hero.Login")
	proto.RegisterType((*GuestRegister)(nil), "hero.GuestRegister")
	proto.RegisterType((*LoginResp)(nil), "hero.LoginResp")
	proto.RegisterType((*VersionUpdateAlert)(nil), "hero.VersionUpdateAlert")
	proto.RegisterType((*GetAchievement)(nil), "hero.GetAchievement")
	proto.RegisterType((*SetAchievement)(nil), "hero.SetAchievement")
	proto.RegisterType((*GetAchievementResp)(nil), "hero.GetAchievementResp")
	proto.RegisterType((*GetArenaLearderboardRankResp)(nil), "hero.GetArenaLearderboardRankResp")
	proto.RegisterType((*GetLearderboardRange)(nil), "hero.GetLearderboardRange")
	proto.RegisterType((*LearderboardInfo)(nil), "hero.LearderboardInfo")
	proto.RegisterType((*GetLearderboardRangeResp)(nil), "hero.GetLearderboardRangeResp")
	proto.RegisterType((*UpdateLearderBoardScore)(nil), "hero.UpdateLearderBoardScore")
	proto.RegisterType((*UpdateLearderBoardScoreResp)(nil), "hero.UpdateLearderBoardScoreResp")
	proto.RegisterType((*UnLockAchievement)(nil), "hero.UnLockAchievement")
	proto.RegisterType((*UnLockAchievementResp)(nil), "hero.UnLockAchievementResp")
	proto.RegisterType((*GetUnLockAchieveDate)(nil), "hero.GetUnLockAchieveDate")
	proto.RegisterType((*GetUnLockAchieveDateResp)(nil), "hero.GetUnLockAchieveDateResp")
	proto.RegisterEnum("hero.CG", CG_name, CG_value)
	proto.RegisterEnum("hero.GC", GC_name, GC_value)
}

func init() { proto.RegisterFile("hero.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 941 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0xdb, 0x6e, 0xdb, 0x46,
	0x10, 0x8d, 0x2e, 0xb4, 0xad, 0x71, 0x25, 0x31, 0x1b, 0xc7, 0x62, 0x6d, 0xc7, 0x31, 0xd8, 0xb4,
	0x0d, 0xdc, 0xc2, 0x2d, 0xdc, 0xe7, 0x02, 0x95, 0x25, 0x82, 0x11, 0xa2, 0x4a, 0xee, 0x4a, 0x32,
	0x8a, 0x3c, 0x94, 0xa0, 0xa9, 0xb5, 0xcc, 0x4a, 0xda, 0x35, 0x96, 0x94, 0x83, 0x3c, 0xf7, 0x07,
	0xfa, 0xd0, 0xaf, 0xe8, 0x4b, 0xef, 0x45, 0x7f, 0xa3, 0xf7, 0xf6, 0x6f, 0x82, 0x9d, 0x25, 0x65,
	0xdd, 0x0c, 0xbf, 0xcd, 0xcc, 0x39, 0x3b, 0x3b, 0x73, 0x66, 0x96, 0x04, 0xb8, 0x64, 0x52, 0x1c,
	0x5d, 0x49, 0x11, 0x0b, 0x92, 0x57, 0xb6, 0xbd, 0x0e, 0x86, 0x33, 0xbe, 0x8a, 0x5f, 0xd9, 0x1f,
	0x80, 0xe1, 0x48, 0x29, 0x24, 0x21, 0x90, 0x0f, 0x44, 0x9f, 0x59, 0x99, 0x83, 0xcc, 0x53, 0x83,
	0xa2, 0xad, 0x62, 0xbe, 0x1c, 0x44, 0x56, 0xf6, 0x20, 0xf7, 0xb4, 0x40, 0xd1, 0xb6, 0x7f, 0xcb,
	0x80, 0xd1, 0x14, 0x83, 0x90, 0x2b, 0x34, 0x1c, 0xb3, 0x10, 0x4f, 0x14, 0x28, 0xda, 0xe4, 0x6d,
	0x28, 0x05, 0xa3, 0x90, 0xf1, 0xd8, 0xbb, 0x66, 0x32, 0x0a, 0x05, 0xb7, 0xb2, 0x88, 0x16, 0x75,
	0xf4, 0x4c, 0x07, 0x89, 0x05, 0xeb, 0xc1, 0xa5, 0xcf, 0x39, 0x1b, 0x59, 0x39, 0xc4, 0x53, 0x97,
	0x3c, 0x86, 0xcd, 0x8b, 0x57, 0x51, 0x7f, 0xe8, 0xc5, 0x62, 0xc8, 0xb8, 0x95, 0x47, 0x14, 0x30,
	0xd4, 0x55, 0x11, 0xf2, 0x08, 0xb4, 0xe7, 0x4d, 0x26, 0x61, 0xdf, 0x32, 0x10, 0x2f, 0x60, 0xa4,
	0x37, 0x09, 0xfb, 0x64, 0x07, 0x36, 0x78, 0x18, 0x0c, 0x5b, 0xfe, 0x98, 0x59, 0x6b, 0x08, 0x4e,
	0x7d, 0xfb, 0x63, 0x28, 0xba, 0x13, 0x16, 0xc5, 0x94, 0x0d, 0xc2, 0x28, 0x66, 0x52, 0x95, 0xe1,
	0x07, 0x81, 0x98, 0xf0, 0x38, 0x69, 0x22, 0x75, 0x55, 0x6f, 0x5c, 0xa5, 0xd0, 0xd5, 0xa3, 0x6d,
	0x7f, 0x9d, 0x85, 0x02, 0x76, 0x4e, 0x59, 0x74, 0x45, 0x1e, 0x41, 0x8e, 0x49, 0x89, 0xe7, 0x36,
	0x8f, 0x37, 0x8f, 0x50, 0x61, 0x54, 0x92, 0xaa, 0x38, 0x31, 0x21, 0xa7, 0xea, 0x53, 0xe7, 0x8b,
	0x54, 0x99, 0x64, 0x0b, 0x0c, 0xdd, 0x93, 0xea, 0xf8, 0x0d, 0xaa, 0x1d, 0xb2, 0x0f, 0xe0, 0x4b,
	0xc6, 0xfd, 0x4e, 0x20, 0x24, 0xc3, 0x76, 0x8b, 0x74, 0x26, 0x42, 0xf6, 0xa0, 0x80, 0x1e, 0xf5,
	0xf9, 0x10, 0xbb, 0x2d, 0xd2, 0x9b, 0x00, 0x79, 0x07, 0x4a, 0xc1, 0xa5, 0x3f, 0x1a, 0x31, 0x3e,
	0x60, 0x3a, 0x83, 0xea, 0x39, 0x4f, 0x17, 0xa2, 0xe4, 0x09, 0x14, 0xa7, 0x11, 0xcc, 0xb4, 0x8e,
	0x99, 0xe6, 0x83, 0xe4, 0x43, 0x58, 0x9b, 0x5c, 0xf5, 0xfd, 0x98, 0x59, 0x1b, 0xd8, 0x95, 0xa5,
	0xbb, 0x4a, 0x86, 0xd6, 0x43, 0xa8, 0x3a, 0x62, 0x32, 0xa6, 0x09, 0xcf, 0xfe, 0x1c, 0xc8, 0x32,
	0xaa, 0x3a, 0xbd, 0x10, 0x32, 0xd0, 0xbb, 0xb4, 0x41, 0xb5, 0xa3, 0xc4, 0x9e, 0xdf, 0x89, 0xd4,
	0x55, 0x62, 0x8f, 0x42, 0x3e, 0x4c, 0x56, 0x01, 0x6d, 0xfb, 0x09, 0x94, 0x5c, 0x16, 0x57, 0x83,
	0xcb, 0x90, 0x5d, 0xb3, 0x31, 0x9b, 0x19, 0x49, 0x66, 0x66, 0x24, 0x9f, 0x40, 0xa9, 0x73, 0x27,
	0x0b, 0x77, 0x42, 0xbc, 0x3c, 0xf3, 0x47, 0x13, 0x3d, 0x50, 0x83, 0x4e, 0x7d, 0xfb, 0x10, 0xc8,
	0xfc, 0x3d, 0x38, 0xdc, 0x2d, 0x30, 0xae, 0x91, 0xae, 0x5f, 0x83, 0x76, 0xec, 0x63, 0xd8, 0x53,
	0x5c, 0xa5, 0x7e, 0x93, 0xf9, 0xb2, 0xcf, 0xe4, 0xb9, 0xf0, 0x65, 0x5f, 0x69, 0x87, 0xa7, 0x08,
	0xe4, 0xa5, 0x12, 0x37, 0x83, 0xe2, 0xa2, 0x6d, 0x53, 0xd8, 0x72, 0x59, 0xbc, 0x40, 0x1f, 0xb0,
	0x95, 0x75, 0x6e, 0x81, 0x11, 0xc5, 0xbe, 0x8c, 0x93, 0xad, 0xd1, 0x8e, 0xda, 0x24, 0xc6, 0xfb,
	0x28, 0x4e, 0x91, 0x2a, 0xd3, 0x3e, 0x05, 0x73, 0x36, 0x61, 0x83, 0x5f, 0x08, 0x3c, 0x8b, 0x0b,
	0x90, 0x54, 0x8c, 0xce, 0xb4, 0xa2, 0xec, 0x4d, 0x45, 0xd3, 0x9b, 0x73, 0x33, 0x3a, 0x7e, 0x01,
	0xd6, 0xaa, 0x2a, 0x53, 0x2d, 0x6e, 0x9e, 0x48, 0x91, 0x6a, 0x87, 0xbc, 0x0f, 0x46, 0xc8, 0x2f,
	0x84, 0xfe, 0x36, 0x6c, 0x1e, 0x6f, 0xeb, 0x55, 0x59, 0x2c, 0x8b, 0x6a, 0x92, 0x5d, 0x83, 0x8a,
	0x5e, 0x90, 0x84, 0x70, 0xa2, 0x08, 0x9d, 0xb4, 0xc4, 0xc5, 0x97, 0x36, 0xdf, 0x4c, 0x3e, 0x69,
	0xc6, 0xfe, 0x0c, 0x76, 0x6f, 0x49, 0x82, 0x75, 0x5a, 0xb0, 0xce, 0xd9, 0x4b, 0x7a, 0x33, 0x80,
	0xd4, 0x55, 0x48, 0x34, 0x09, 0x02, 0x16, 0x45, 0x78, 0xcb, 0x06, 0x4d, 0x5d, 0xfb, 0x5d, 0xb8,
	0xdf, 0xe3, 0x4d, 0x11, 0x0c, 0xef, 0x5a, 0xb4, 0xf7, 0xe0, 0xe1, 0x12, 0x31, 0x9d, 0x39, 0xbe,
	0x98, 0x84, 0x8c, 0xaf, 0xe2, 0x10, 0x67, 0x3e, 0xc7, 0xaf, 0xfb, 0xf1, 0xca, 0x99, 0xdb, 0x47,
	0xa8, 0xfc, 0x12, 0xf7, 0xb6, 0xdc, 0x87, 0x5f, 0x65, 0x21, 0x5b, 0x73, 0x49, 0x09, 0xa0, 0x51,
	0xf7, 0x5e, 0x38, 0xb4, 0xed, 0xd5, 0x5c, 0xf3, 0x1e, 0x79, 0x00, 0xa5, 0x46, 0xdd, 0x73, 0x7b,
	0x4e, 0xa7, 0xeb, 0x35, 0xdb, 0x6e, 0xa3, 0x65, 0xfe, 0x5e, 0x26, 0xdb, 0x70, 0x7f, 0x1a, 0xa4,
	0x8e, 0xdb, 0xe8, 0x74, 0x1d, 0x6a, 0xfe, 0x51, 0x26, 0x15, 0x20, 0x2a, 0xee, 0x74, 0xbd, 0x6a,
	0xed, 0x59, 0xc3, 0x39, 0x73, 0x3e, 0x75, 0x5a, 0x5d, 0xf3, 0xcf, 0x14, 0xe8, 0x2c, 0x00, 0x7f,
	0x4d, 0x33, 0x29, 0x80, 0x3a, 0xad, 0xaa, 0x47, 0xab, 0xad, 0xe7, 0xe6, 0xdf, 0x65, 0xb2, 0x0f,
	0x6f, 0x26, 0xf1, 0xa6, 0x53, 0xa5, 0x75, 0x87, 0x9e, 0xb4, 0xab, 0xb4, 0xae, 0x60, 0xd7, 0x31,
	0xff, 0x29, 0x93, 0x1d, 0x78, 0xd8, 0xa8, 0x7b, 0xbd, 0xd3, 0x7a, 0xb5, 0xeb, 0x68, 0x8a, 0xd7,
	0xa9, 0xb5, 0xa9, 0x63, 0xfe, 0x3b, 0xc5, 0x5a, 0xcd, 0x76, 0xed, 0xf9, 0xdc, 0x7d, 0xff, 0xcd,
	0xe6, 0x9d, 0xc7, 0x55, 0x1e, 0xf3, 0xff, 0xf2, 0xe1, 0x37, 0x59, 0xc8, 0xba, 0xb5, 0x59, 0x15,
	0xdc, 0x9a, 0x79, 0x8f, 0x58, 0xf0, 0x60, 0x5e, 0x05, 0x8f, 0x3a, 0x9d, 0x53, 0xf3, 0xd7, 0x0a,
	0xd9, 0x83, 0xca, 0x72, 0xcb, 0x1a, 0xfd, 0x25, 0x45, 0x3b, 0xab, 0xd0, 0x9f, 0x2b, 0x64, 0x17,
	0xb6, 0x97, 0x9a, 0xd7, 0xe0, 0x4f, 0x15, 0xf2, 0x16, 0xec, 0xdf, 0xaa, 0x80, 0x26, 0xfd, 0x58,
	0x21, 0x07, 0xb0, 0xbb, 0x52, 0x86, 0x84, 0xf1, 0x43, 0x85, 0x3c, 0x86, 0x9d, 0x95, 0x62, 0x68,
	0xc2, 0xf7, 0xb3, 0xf7, 0x2c, 0x2b, 0xa2, 0x49, 0xdf, 0x55, 0x4e, 0xca, 0xcf, 0x72, 0xdf, 0x66,
	0xe1, 0x54, 0x8a, 0x2f, 0x5f, 0x88, 0xf1, 0x79, 0xc8, 0xce, 0xd7, 0xf0, 0x9f, 0xff, 0xd1, 0xeb,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xb2, 0xe7, 0xf0, 0xb1, 0x01, 0x08, 0x00, 0x00,
}
