package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/pkg/errors"
)

//这里是AES加密解密的类

//秘钥可以是16， 24， 32位，分别对应AES-128，AES-192，AES-256加密方法
var PwdKey = []byte("DIS**#KKKDJJSKDI")

//PKCS7填充模式
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	//这里的padding指的是将cipherText用大小为blockSize的来分割后，剩下的一些值
	//例如1010和100,此时padding为90
	padding := blockSize - len(cipherText)%blockSize
	//获取一个字节数组，长度为padding
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	//将这个字节数组拼到cipherText后面，让其变成一个符合加密长度的字节数组(1010+90 此时是100的倍数)
	return append(cipherText, padText...)
}

//删除PKCS7填充的值
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取原始数据的长度
	l := len(origData)
	if l == 0 {
		return nil, errors.New("加密字符串长度错误")
	} else {
		//获取填充字符串的长度,这里为什么是减1呢，因为上面填充的时候，是以padding这个值连填充的
		//所以最后一个值的int值就是padding的长度
		unpadding := int(origData[l-1])
		return origData[:(l - unpadding)], nil
	}

}

//实现加密
func AesEcrypt(origData, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中的CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密，将加密后的数据写入到crypted中
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//实现解密
func AesDeCrypt(cypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建解密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充的字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

//base64加密
func EnPwdCode(source []byte) (string, error) {
	result, err := AesEcrypt(source, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

//base64解密
func DePwdCode(source string) ([]byte, error) {
	sourceByte, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		return nil, err
	}
	return AesDeCrypt(sourceByte, PwdKey)
}
