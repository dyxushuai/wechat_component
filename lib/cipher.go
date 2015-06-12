package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 用于微信消息的加解密
// messageCrypter 封装了生成签名和消息加解密的方法
type messageCrypter struct {
	token string
	appID string
	key   []byte
	iv    []byte
}

// NewmessageCrypter 方法用于创建 messageCrypter 实例
//
// token 为开发者在微信开放平台上设置的 Token，
// encodingAESKey 为开发者在微信开放平台上设置的 EncodingAESKey，
// appID 为企业号的 CorpId 或者 AppId
func newmessageCrypter(token, encodingAESKey, appID string) (messageCrypter, error) {
	var key []byte
	var err error

	if key, err = aesKeyDecode(encodingAESKey); err != nil {
		return messageCrypter{}, err
	}

	iv := key[:16]

	return messageCrypter{
		token,
		appID,
		key,
		iv,
	}, nil
}

func aesKeyDecode(encodedAESKey string) (key []byte, err error) {
	if len(encodedAESKey) != 43 {
		err = errors.New("the length of encodedAESKey must be equal to 43")
		return
	}
	key, err = base64.StdEncoding.DecodeString(encodedAESKey + "=")
	if err != nil {
		return
	}
	if len(key) != 32 {
		err = errors.New("encodingAESKey invalid")
		return
	}
	return
}

// getSignature 方法用于返回签名
func (w messageCrypter) getSignature(timestamp, nonce, msgEncrypt string) string {
	sl := []string{w.token, timestamp, nonce, msgEncrypt}
	sort.Strings(sl)

	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))

	return fmt.Sprintf("%x", s.Sum(nil))
}

// Decrypt 方法用于对密文进行解密
//
// 返回解密后的消息，CropId/AppId, 或者错误信息
func (w messageCrypter) decrypt(text string) ([]byte, string, error) {
	var msgDecrypt []byte
	var id string

	deciphered, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, "", err
	}

	c, err := aes.NewCipher(w.key)
	if err != nil {
		return nil, "", err
	}

	cbc := cipher.NewCBCDecrypter(c, w.iv)
	cbc.CryptBlocks(deciphered, deciphered)

	decoded := delDecode(deciphered)

	buf := bytes.NewBuffer(decoded[16:20])

	var msgLen int32
	binary.Read(buf, binary.BigEndian, &msgLen)

	msgDecrypt = decoded[20 : 20+msgLen]
	id = string(decoded[20+msgLen:])

	return msgDecrypt, id, nil
}

// Encrypt 方法用于对明文进行加密
func (w messageCrypter) encrypt(text string) (string, error) {
	message := []byte(text)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, int32(len(message))); err != nil {
		return "", err
	}

	msgLen := buf.Bytes()

	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}

	messageBytes := bytes.Join([][]byte{randBytes, msgLen, message, []byte(w.appID)}, nil)

	encoded := fillEncode(messageBytes)

	c, err := aes.NewCipher(w.key)
	if err != nil {
		return "", err
	}

	cbc := cipher.NewCBCEncrypter(c, w.iv)
	cbc.CryptBlocks(encoded, encoded)

	return base64.StdEncoding.EncodeToString(encoded), nil
}

// delDecode 方法用于删除解密后明文的补位字符
func delDecode(text []byte) []byte {
	pad := int(text[len(text)-1])

	if pad < 1 || pad > 32 {
		pad = 0
	}

	return text[:len(text)-pad]
}

// fillEncode 方法用于对需要加密的明文进行填充补位
func fillEncode(text []byte) []byte {
	const BlockSize = 32

	amountToPad := BlockSize - len(text)%BlockSize

	for i := 0; i < amountToPad; i++ {
		text = append(text, byte(amountToPad))
	}

	return text
}

// 包装读写接口 使其能读写加解密
type IOCipher interface {
	Encrypt(w io.Writer, b []byte) (err error)
	Decrypt(r io.Reader) (b []byte, err error)
}

// 用于管道加密
type Cipher struct {
	messageCrypter
	token string
}

func NewCipher(token, encodingAESKey, appID string) (IOCipher, error) {
	mc, err := newmessageCrypter(token, encodingAESKey, appID)
	if err != nil {
		return nil, err
	}

	return &Cipher{
		messageCrypter: mc,
		token:          token,
	}, nil
}

// xml cdata
type charData struct {
	Text []byte `xml:",innerxml"`
}

func newCharData(s string) charData {
	return charData{[]byte("<![CDATA[" + s + "]]>")}
}

type cipherToWX struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      charData
	MsgSignature charData
	TimeStamp    string
}

// 将b加密写入w
func (c *Cipher) Encrypt(w io.Writer, b []byte) (err error) {
	result, err := c.messageCrypter.encrypt(string(b))
	if err != nil {
		return
	}

	to := &cipherToWX{
		Encrypt:   newCharData(result),
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
	}
	// timestamp, nonce, encryptedMsg string)
	to.MsgSignature = newCharData(MsgSign(c.token, to.TimeStamp, createNonceStr(16), result))
	err = xml.NewEncoder(w).Encode(to)
	return
}

// 来自微信的密文结构
type cipherFromWX struct {
	ToUserName string
	Encrypt    string
}

// 从r读取并解密
func (c *Cipher) Decrypt(r io.Reader) (b []byte, err error) {
	from := &cipherFromWX{}
	err = xml.NewDecoder(r).Decode(from)
	if err != nil {
		return
	}

	b, _, err = c.messageCrypter.decrypt(from.Encrypt)
	return
}
