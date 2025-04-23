// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package bot

import (
	"context"

	"keyboard-api-go/api/bot/v1"
)

type IBotV1 interface {
	CreateBot(ctx context.Context, req *v1.CreateBotReq) (res *v1.CreateBotRes, err error)
	FindBotByAddress(ctx context.Context, req *v1.FindBotByAddressReq) (res *v1.FindBotByAddressRes, err error)
	Active(ctx context.Context, req *v1.ActiveReq) (res *v1.ActiveRes, err error)
	SubmitWork(ctx context.Context, req *v1.SubmitWorkReq) (res *v1.SubmitWorkRes, err error)
	Rank(ctx context.Context, req *v1.RankReq) (res *v1.RankRes, err error)
	SetAutoAcceptFriend(ctx context.Context, req *v1.SetAutoAcceptFriendReq) (res *v1.SetAutoAcceptFriendRes, err error)
	FriendList(ctx context.Context, req *v1.FriendListReq) (res *v1.FriendListRes, err error)
	AddFriend(ctx context.Context, req *v1.AddFriendReq) (res *v1.AddFriendRes, err error)
	DelFriend(ctx context.Context, req *v1.DelFriendReq) (res *v1.DelFriendRes, err error)
	FriendMsgList(ctx context.Context, req *v1.FriendMsgListReq) (res *v1.FriendMsgListRes, err error)
	HandleFriendMsg(ctx context.Context, req *v1.HandleFriendMsgReq) (res *v1.HandleFriendMsgRes, err error)
}
