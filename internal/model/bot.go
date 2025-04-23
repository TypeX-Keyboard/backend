package model

type BotInfo struct {
	Address        string  `json:"address"    orm:"address"     description:"publicKeyAddress"`
	Amount         int64   `json:"amount"     orm:"amount"      description:"Number of EVEs"`
	Decimals       int     `json:"decimals" dc:"EVE decimal"`
	TypeRate       int     `json:"typeRate"   orm:"type_rate"   description:"Typing progress"`
	TypeLimit      int64   `json:"typeLimit"   description:"Progression caps"`
	TypeCount      int64   `json:"typeCount"  orm:"type_count"  description:"Total number of typing"`
	Active         int     `json:"active"     orm:"active"      description:"1 = active, 0 = inactive"`
	Power          float64 `json:"power"      orm:"power"       description:"Current power"`
	PowerLimit     int     `json:"powerLimit" orm:"power_limit" description:"Power cap"`
	AutoAccept     int     `json:"autoAccept" orm:"auto_accept" description:"1=Automatically accepted"`
	Buff           int     `json:"buff"       orm:"buff"        description:"Gain (in %)"`
	BuffLimit      int     `json:"buffLimit" dc:"Maximum Gain (Unit%)"`
	SelfRank       int     `json:"selfRank"`
	MiningRateOfH  float64 `json:"miningRateOfH" dc:"Mining efficiency per hour"`
	ConsumptionOfH float64 `json:"consumptionOfH" dc:"Electricity consumption per hour"`
}

type RankBot struct {
	Rank    int    `json:"rank"`
	Address string `json:"address"`
	Amount  int64  `json:"amount" dc:"Number of EVEs"`
}

type FriendInfo struct {
	Address string `json:"address"`
	Active  bool   `json:"active"`
}

type FriendMsg struct {
	Id      int64  `json:"id"`
	Address string `json:"address"`
	Status  int    `json:"status" dc:"3=Rejected, 2=Approved, 1=Application"`
}

type ActiveFriends struct {
	Address     string `json:"address"`
	FriendCount int    `json:"friendCount"`
	ActiveCount int    `json:"activeCount"`
}
