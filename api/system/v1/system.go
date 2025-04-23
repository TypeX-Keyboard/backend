package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"keyboard-api-go/internal/model"
)

type SettingsReq struct {
	g.Meta `path:"/settings" tags:"系统" method:"get" summary:"系统设置"`
}

type SettingsRes struct {
	g.Meta             `mime:"application/json"`
	Version            model.Version `json:"version" dc:"版本号"`
	Timezone           string        `json:"timezone" dc:"时区"`
	WorkReportInterval int           `json:"work_report_interval" dc:"工作汇报间隔"`
	MiningRateOfH      float64       `json:"miningRateOfH" dc:"每小时挖矿效率"`
	PrivacyPolicyURL   string        `json:"privacyPolicyUrl" dc:"隐私政策URL"`
	UserAgreement      string        `json:"userAgreement" dc:"用户协议URL"`
}
