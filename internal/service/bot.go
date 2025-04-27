// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"keyboard-api-go/internal/model"
	"keyboard-api-go/internal/model/entity"
)

type (
	IBot interface {
		Active(ctx context.Context, address string) (err error)
		Entity2BotInfo(ctx context.Context, bot *entity.Bot) *model.BotInfo
		FindByAddress(ctx context.Context, address string) (bot *entity.Bot, err error)
		CreateBot(ctx context.Context, address string, deviceId string) (bot *entity.Bot, err error)
		UpdateAutoAccept(ctx context.Context, address string, autoAccept bool) (err error)
		FriendList(ctx context.Context, address string) (res []model.FriendInfo, err error)
		AddFriend(ctx context.Context, address string, friendAddress string) (bool, error)
		FriendMsgList(ctx context.Context, address string, isSelf bool, page int, size int) (res []model.FriendMsg, total int, err error)
		HandleFriendMsg(ctx context.Context, id int64, address string, accept bool) (err error)
		DelFriend(ctx context.Context, address string, friendAddress string) error
		ActiveTask(ctx context.Context) error
		ActiveFriendCounts(ctx context.Context) ([]model.ActiveFriends, error)
		MiningTask(ctx context.Context) error
		EveOf1M() int64
		Rank(ctx context.Context, address string) (rank []model.RankBot, selfRank int, selfAmount int64, err error)
		SubmitWork(ctx context.Context, address string, typeCount int) (err error)
	}
)

var (
	localBot IBot
)

func Bot() IBot {
	if localBot == nil {
		panic("implement not found for interface IBot, forgot register?")
	}
	return localBot
}

func RegisterBot(i IBot) {
	localBot = i
}
