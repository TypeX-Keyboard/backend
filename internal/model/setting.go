package model

type Setting struct {
	PublicKey                string `json:"PublicKey"`
	PrivateKey               string `json:"PrivateKey"`
	Timezone                 string `json:"Timezone"`
	WorkReportInterval       int    `json:"WorkReportInterval"`
	BasePassportNO           int    `json:"BasePassportNO"`
	PortraitGenerateMaxTimes int    `json:"PortraitGenerateMaxTimes"`
	RecentEarningsType       string `json:"RecentEarningsType"`
	RecentEarningsNum        int    `json:"RecentEarningsNum"`
	LowestRankToCalc         int    `json:"LowestRankToCalc"`
	MaxLostReportNumAllowed  int    `json:"MaxLostReportNumAllowed"`
	ComputingPointsPerMinute int    `json:"ComputingPointsPerMinute"`
	TypingPointsPerWord      int    `json:"TypingPointsPerWord"`
	UpPointsPerWord          int    `json:"UpPointsPerWord"`
}

var SettingConfig *Setting
