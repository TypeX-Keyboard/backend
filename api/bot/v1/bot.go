package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/model"
)

type CreateBotReq struct {
	g.Meta  `path:"/create" tags:"bot" method:"post" summary:"create a bot"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
}

type CreateBotRes model.BotInfo

type FindBotByAddressReq struct {
	g.Meta  `path:"/find" tags:"bot" method:"get" summary:"Query bots based on address"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
}

type FindBotByAddressRes model.BotInfo

type ActiveReq struct {
	g.Meta  `path:"/active" tags:"bot" method:"get" summary:"Report active status"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
}

type ActiveRes struct {
}

type SubmitWorkReq struct {
	g.Meta    `path:"/submitWork" tags:"bot" method:"get" summary:"Report the workload"`
	Address   string `json:"address" dc:"walletAddress" required:"true"`
	TypeCount int    `json:"typeCount" dc:"The number of characters to be reported this time"`
}

type SubmitWorkRes struct {
}

type RankReq struct {
	g.Meta  `path:"/rank" tags:"bot" method:"get" summary:"league tables"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type RankRes struct {
	List     []model.RankBot `json:"list"`
	SelfRank int             `json:"selfRank"`
	Amount   int64           `json:"amount"`
}

type SetAutoAcceptFriendReq struct {
	g.Meta  `path:"/setAutoAcceptFriend" tags:"bot" method:"post" summary:"Set whether or not to automatically accept friends"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
	Auto    bool   `json:"auto" dc:"Whether or not to automatically accept friends"`
}

type SetAutoAcceptFriendRes struct {
}

type FriendListReq struct {
	g.Meta  `path:"/friendList" tags:"bot" method:"get" summary:"Get a list of friends"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
}

type FriendListRes []model.FriendInfo

type AddFriendReq struct {
	g.Meta        `path:"/addFriend" tags:"bot" method:"post" summary:"Add friends(UserNotFoundErr = 501;CanNotAddFriendSelfErr = 502；FriendLimitErr = 503; FriendAddedErr = 504)"`
	Address       string `json:"address" dc:"walletAddress" required:"true"`
	FriendAddress string `json:"friendAddress" dc:"Friend s wallet address" required:"true"`
}

type AddFriendRes struct {
	IsAutoAccept bool `json:"isAutoAccept"`
}

type DelFriendReq struct {
	g.Meta        `path:"/delFriend" tags:"bot" method:"post" summary:"Delete a friend"`
	Address       string `json:"address" dc:"walletAddress" required:"true"`
	FriendAddress string `json:"friendAddress" dc:"Friend's wallet address" required:"true"`
}

type DelFriendRes struct {
}

type FriendMsgListReq struct {
	g.Meta  `path:"/friendMsgList" tags:"bot" method:"get" summary:"Get a list of friends messages"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
	IsSelf  bool   `json:"isSelf" dc:"true=i add others，false=someone else added me"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type FriendMsgListRes struct {
	List  []model.FriendMsg `json:"list"`
	Total int               `json:"total"`
}

type HandleFriendMsgReq struct {
	g.Meta  `path:"/handleFriendMsg" tags:"bot" method:"post" summary:"Handle friend messages"`
	Address string `json:"address" dc:"walletAddress" required:"true"`
	Id      int64  `json:"id" required:"true"`
	Accept  bool   `json:"accept"`
}

type HandleFriendMsgRes struct {
}
