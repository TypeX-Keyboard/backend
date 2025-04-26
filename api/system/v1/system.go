package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/model"
)

type SettingsReq struct {
	g.Meta `path:"/settings" tags:"system" method:"get" summary:"System settings"`
}

type SettingsRes struct {
	g.Meta             `mime:"application/json"`
	Version            model.Version `json:"version" dc:"Version number"`
	Timezone           string        `json:"timezone" dc:"time zone"`
	WorkReportInterval int           `json:"work_report_interval" dc:"Debriefing intervals"`
	MiningRateOfH      float64       `json:"miningRateOfH" dc:"Mining efficiency per hour"`
	PrivacyPolicyURL   string        `json:"privacyPolicyUrl" dc:"Privacy Policy URL"`
	UserAgreement      string        `json:"userAgreement" dc:"User Agreement URL"`
}
