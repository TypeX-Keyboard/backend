package consts

const (
	SignParamKey         = "sign"
	DeviceIdParamKey     = "device_id"
	OS                   = "os"
	AesKeyCacheKeyPrefix = "AES:"
	LANGUAGE             = "language"
	TIMEZONE             = "timezone"
)

const ListenTokenPrices = "ws:listen_token_prices"

const RedisProxyUsKey = "proxy_us"

const (
	AccountProgram   = "11111111111111111111111111111111"
	TokenProgram     = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	Token2022Program = "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"
	SolAddress       = "So11111111111111111111111111111111111111112"
)

const (
	HockeyAppVersionPrefix = "APP.VERSION."
)

// bot
const (
	PowerLimit    = 120 // 满电120个单位
	EVEDecimal    = 6
	MiningRateOfH = 1
	BuffStep      = 3   // buff增加3个单位
	BuffLimit     = 400 // buff上限400个单位
	FriendLimit   = 10  // 好友上限10个
	PowerOfOne    = 60  // 每电60个单位
	TypeLimit     = 100 // 每电100个字符
	RankLimit     = 100 // 排行榜限制

	SubmitWorkLock = "SUBMIT_WORK_LOCK:%s" // 上报工作锁
	AddFriendLock  = "ADD_FRIEND_LOCK:%s"  // 添加好友锁
	PowerLock      = "POWER_LOCK"          // power锁

	RedisActivePrefix = "ACTIVE" // 活跃状态
)

// friend status
const (
	FriendPending = iota + 1
	FriendAccept
	FriendReject
)

// 钱包地址Eign4g779oTaexjVnQt6gcQh3SXa1t1sRca4ajrNbvFK
const FeeAccountATA = "D83BGpPKkAzg5YRj4hQ25gkgkK1sUGVPko5cb1wfj2BK"
const FeeAccount = "Eign4g779oTaexjVnQt6gcQh3SXa1t1sRca4ajrNbvFK"

const (
	OKXGetTXList    = "Okx:GetTxList"
	TxSuccessHandle = "TxSuccessHandle"

	PriceIn0 = "PriceIn0:%s"

	SolPrice = "SolPrice"
)
