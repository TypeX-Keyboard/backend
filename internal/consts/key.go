package consts

const PublicKey = `-----BEGIN PUBLIC KEY-----

-----END PUBLIC KEY-----`

const PrivateKey = `-----BEGIN PRIVATE KEY-----

-----END PRIVATE KEY-----`

const (
	RedisClientPrivateKey = "Client:PrivateKey"
	RedisClientPublicKey  = "Client:PublicKey"
	RedisClientHMACKey    = "Client:HMACKey"

	ExpireH = 24
)
