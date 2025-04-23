package secure

import (
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

func GetMD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func ECBEncrypt(pt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt) // pad last block of plaintext if block size less than block cipher size
	if err != nil {
		return nil, err
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct, nil
}

func ECBDecrypt(ct, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		return nil, err
	}
	return pt, nil
}

func Sha256Encrypt(pt, publicKeyPEM []byte) ([]byte, error) {
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKeyBlock == nil {
		return nil, errors.New("Failed to parse public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Failed to convert to RSA public key")
	}
	// 使用公钥加密数据
	encryptedData, err := encryptRSA(pt, rsaPublicKey)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

func Sha256Decrypt(encryptedData, privateKeyPEM []byte) ([]byte, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	if privateKeyBlock == nil {
		return nil, errors.New("Failed to parse private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	// 使用私钥解密数据
	decryptedData, err := decryptRSA(encryptedData, privateKey.(*rsa.PrivateKey))
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func encryptRSA(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func decryptRSA(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func GenerateKeyPair(bits int) ([]byte, []byte, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	b, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	// 编码私钥为 PEM 格式
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	})

	// 提取公钥
	publicKey := &privateKey.PublicKey

	// 编码公钥为 PEM 格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}

func GeneratedSecretKey(keyLength int) (string, error) {
	// 创建一个字节切片来存储密钥
	secretKey := make([]byte, keyLength)

	// 生成随机字节
	_, err := rand.Read(secretKey)
	if err != nil {
		fmt.Println("Error generating random key:", err)
		return "", err
	}
	// 将密钥转换为十六进制字符串
	secretKeyHex := hex.EncodeToString(secretKey)
	return secretKeyHex, nil
}
