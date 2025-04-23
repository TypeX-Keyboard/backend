package secure

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/model"
	"log"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func initSetting(ctx context.Context) {
	// 初始化设置
	c := g.Config().MustGet(ctx, "setting")
	c.Struct(&model.SettingConfig)
	g.Log().Info(ctx, "init setting success")
}
func TestAes(t *testing.T) {
	initSetting(gctx.New())
	keystr := GetMD5("key")
	key := []byte(keystr)
	fmt.Println(keystr)
	data := "123"
	log.Println("明文：", data)
	res, err := ECBEncrypt([]byte(data), key)
	if err != nil {
		log.Println("加密失败：", err)
		return
	}
	e := base64.StdEncoding.EncodeToString(res)
	log.Println("密文结果：", e)

	d, _ := base64.StdEncoding.DecodeString(e)
	res1, _ := ECBDecrypt(d, key)
	r := string(res1)
	log.Println("解密结果：", r)
}

func removePKCS7Padding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	paddingLen := int(data[length-1])
	if paddingLen > length || paddingLen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:length-paddingLen], nil
}

func TestAes2(t *testing.T) {
	initSetting(gctx.New())
	// 生成随机的 AES-256 密钥
	// aesKey := make([]byte, 32) // AES-256 密钥长度为 32 字节
	// if _, err := rand.Read(aesKey); err != nil {
	// 	fmt.Println("Error generating AES key:", err)
	// 	return
	// }
	privateKeyPEM := []byte(consts.PrivateKey)
	publicKeyPEM := []byte(consts.PublicKey)
	// 要加密的数据
	data := []byte("81dc9bdb52d04dc20036dbd8313ed055")

	// 使用公钥加密数据
	encryptedData, err := Sha256Encrypt(data, publicKeyPEM)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	// 使用私钥解密数据
	decryptedData, err := Sha256Decrypt(encryptedData, privateKeyPEM)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Println("Original data:", len(data))
	fmt.Println("Original data:", string(data))
	fmt.Println("Original data:", base64.StdEncoding.EncodeToString(data))
	fmt.Println("Encrypted data:", base64.StdEncoding.EncodeToString(encryptedData))
	fmt.Println("Decrypted data:", base64.StdEncoding.EncodeToString(decryptedData))
}

func TestGenerateKeyPair(t *testing.T) {
	privateKeyB, publicKeyB, err := GenerateKeyPair(2048)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}
	g.Dump(string(privateKeyB))
	g.Dump(consts.PrivateKey)
	g.Dump(string(publicKeyB))
	g.Dump(consts.PublicKey)
	initSetting(gctx.New())
	privateKeyPEM := privateKeyB
	publicKeyPEM := publicKeyB
	// 要加密的数据
	data := []byte("81dc9bdb52d04dc20036dbd8313ed055")

	// 使用公钥加密数据
	encryptedData, err := Sha256Encrypt(data, publicKeyPEM)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	// 使用私钥解密数据
	decryptedData, err := Sha256Decrypt(encryptedData, privateKeyPEM)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Println("Original data:", len(data))
	fmt.Println("Original data:", string(data))
	fmt.Println("Original data:", base64.StdEncoding.EncodeToString(data))
	fmt.Println("Encrypted data:", base64.StdEncoding.EncodeToString(encryptedData))
	fmt.Println("Decrypted data:", base64.StdEncoding.EncodeToString(decryptedData))
}

func TestHMACSHA256(t *testing.T) {
	// 定义密钥长度（以字节为单位，例如 32 字节表示 256 位）
	keyLength := 32

	// 创建一个字节切片来存储密钥
	secretKey := make([]byte, keyLength)

	// 生成随机字节
	_, err := rand.Read(secretKey)
	if err != nil {
		fmt.Println("Error generating random key:", err)
		return
	}
	// 将密钥转换为十六进制字符串
	secretKeyHex := hex.EncodeToString(secretKey)

	fmt.Println("Generated Secret Key:", secretKeyHex)
	// 原始数据
	message := "your message here"

	// 创建 HMAC 使用 SHA-256
	h := hmac.New(sha256.New, []byte(secretKey))

	// 写入数据
	h.Write([]byte(message))

	// 计算 HMAC
	hmac := h.Sum(nil)

	// 转换为十六进制字符串
	hmacHex := hex.EncodeToString(hmac)

	fmt.Println("HMAC (SHA-256):", hmacHex)
}
