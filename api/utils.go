package api

import (
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"

	"github.com/russross/blackfriday"
)

var wxCrypter = &WXBizDataCrypt{
	AppID:  APPID,
	Secret: SECRET,
}

// WXBizDataCrypt is wx biz data crypt and decrypt util
type WXBizDataCrypt struct {
	AppID, Secret string
}

// DecryptData is used for decrypt
func (c *WXBizDataCrypt) DecryptData(encryptedData, iv string) (string, error) {
	keyB, err := b64.StdEncoding.DecodeString(c.Secret)
	ivB, err := b64.StdEncoding.DecodeString(iv)
	encryptedDataB, _ := b64.StdEncoding.DecodeString(encryptedData)
	cip, err := aes.NewCipher(keyB)
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCBCDecrypter(cip, ivB)
	data := make([]byte, len(encryptedDataB))
	copy(data, encryptedDataB)
	decrypter.CryptBlocks(data, data)
	return string(data), nil
}

// MarkdownIt is ...
func MarkdownIt(md string) string {
	out := blackfriday.MarkdownCommon([]byte(md))
	return string(out)
}
