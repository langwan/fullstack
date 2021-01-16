package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)



func SymmetricEncryptSB(key string, data []byte) []byte {
	hexKey := []byte(key)
	return SymmetricEncryptBB(hexKey, data)
}

func SymmetricEncryptBB(hexKey []byte, data []byte) []byte {
	block, err := aes.NewCipher(hexKey)
	if err != nil {
		panic(err)
	}
	plaintext, _ := pkcs7Pad(data, block.BlockSize())
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(cipherText[aes.BlockSize:], plaintext)
	return cipherText
}

func SymmetricDecryptSB(key string, data []byte) ([]byte, error) {
	// Key
	hexKey := []byte(key)
	return SymmetricDecryptBB(hexKey, data)
}

func SymmetricDecryptBB(hexKey []byte, data []byte) ([]byte, error) {

	// Create the AES cipher
	block, err := aes.NewCipher(hexKey)
	if err != nil {
		panic(err)
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]

	if len(cipherText)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")
	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")
	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

func pkcs7Pad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blockSize - (len(b) % blockSize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

func pkcs7Unpad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blockSize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

func AsymmetricEncryptSB(pubKey string, data []byte) ([]byte, error) {
	pubData, _ := base64.StdEncoding.DecodeString(pubKey)
	pub, _ := x509.ParsePKCS1PublicKey(pubData)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func AsymmetricDecryptSB(priKey string, data []byte) ([]byte, error) {
	priData, _ := base64.StdEncoding.DecodeString(priKey)
	pri, _ := x509.ParsePKCS1PrivateKey(priData)
	return rsa.DecryptPKCS1v15(rand.Reader, pri, data)
}

func H2S(hexString string) string {
	hex, _ := hex.DecodeString(hexString)
	return string(hex)
}

func S2H(string string) string {
	return hex.EncodeToString([]byte(string))
}


func RSACreate() (pub string, pri string) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	bytePri := x509.MarshalPKCS1PrivateKey(privateKey)
	pri = base64.StdEncoding.EncodeToString(bytePri)
	bytePub := x509.MarshalPKCS1PublicKey(publicKey)
	pub = base64.StdEncoding.EncodeToString(bytePub)
	return pub, pri
}

func Hash(str string) string {
	mac := hmac.New(sha256.New, []byte(str))
	b := mac.Sum(nil)
	return hex.EncodeToString(b)
}

const (
	LETTERS_LETTER  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LETTERS_NUMBEER = "0123456789"
	LETTERS_SYMBOL  = "~`!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"
)

func RandString(n int, letters ...string) (string, error) {

	lettersDefaultValue := LETTERS_LETTER + LETTERS_NUMBEER + LETTERS_SYMBOL

	if len(letters) > 0 {
		lettersDefaultValue = letters[0]
	}

	bytes := make([]byte, n)

	_, err := rand.Read(bytes)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = lettersDefaultValue[b%byte(len(lettersDefaultValue))]
	}

	return string(bytes), nil
}