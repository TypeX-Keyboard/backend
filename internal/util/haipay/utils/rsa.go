package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	RSAKeySize = 2048
)

func parsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte("-----BEGIN PRIVATE KEY-----\n" + privateKey + "\n-----END PRIVATE KEY-----"))
	if block == nil {
		return nil, errors.New("failed to parse private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not RSA private key")
	}

	return priv, nil
}

func parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + publicKey + "\n-----END PUBLIC KEY-----"))
	if block == nil {
		return nil, errors.New("failed to parse public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func RSAEncrypt(data string, pubKey string) (string, error) {
	publicKey, err := parsePublicKey(pubKey)
	if err != nil {
		return "", err
	}

	partLen := publicKey.Size() - 11
	var buffer bytes.Buffer
	dataBytes := []byte(data)

	for i := 0; i < len(dataBytes); i += partLen {
		end := i + partLen
		if end > len(dataBytes) {
			end = len(dataBytes)
		}
		encryptedBlock, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, dataBytes[i:end])
		if err != nil {
			return "", err
		}
		buffer.Write(encryptedBlock)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func RSADecrypt(cipherText string, privKey string) (string, error) {
	privateKey, err := parsePrivateKey(privKey)
	if err != nil {
		return "", err
	}

	cipherData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	partLen := privateKey.PublicKey.Size()
	var buffer bytes.Buffer

	for i := 0; i < len(cipherData); i += partLen {
		end := i + partLen
		if end > len(cipherData) {
			end = len(cipherData)
		}
		decryptedBlock, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherData[i:end])
		if err != nil {
			return "", err
		}
		buffer.Write(decryptedBlock)
	}

	return buffer.String(), nil
}

func SignWithPrivateKey(data string, privKey string) (string, error) {
	privateKey, err := parsePrivateKey(privKey)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifyWithPublicKey(data string, pubKey string, sign string) (bool, error) {
	publicKey, err := parsePublicKey(pubKey)
	if err != nil {
		return false, err
	}

	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, signature)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetSign(data map[string]interface{}, secretKey string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if data[k] != nil && data[k] != "" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := fmt.Sprintf("%v", data[k])
		if v != "" {
			sb.WriteString(fmt.Sprintf("%s=%s&", k, v))
		}
	}
	sb.WriteString("key=" + secretKey)
	return sb.String()
}
