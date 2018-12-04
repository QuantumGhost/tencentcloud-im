package signature

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spacemonkeygo/openssl"
	"strconv"

	"io/ioutil"

	"encoding/base64"
	"time"
)

const (
	encodeSign = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789*-"
	bufSize    = 300
)

var signEncoding = base64.NewEncoding(encodeSign).WithPadding('_')

type tlsJSONPayload struct {
	AppIDAt3rd  string `json:"TLS.appid_at_3rd"`
	AccountType string `json:"TLS.account_type"`
	Identifier  string `json:"TLS.identifier"`
	SDKAppId    string `json:"TLS.sdk_appid"`
	// 创建时间戳，精确到秒
	Time string `json:"TLS.time"`
	// 过期时间，精确到秒
	ExpireAfter string `json:"TLS.expire_after"`
	Version     string `json:"TLS.version"`
	Signature   string `json:"TLS.sig"`
}

func (p *tlsJSONPayload) ToSignString() string {
	const (
		LINE_END  = '\n'
		SEPERATOR = ':'
	)

	buf := bytes.NewBuffer(make([]byte, 0, bufSize))
	pairs := []string{
		"TLS.appid_at_3rd", p.AppIDAt3rd,
		"TLS.account_type", p.AccountType,
		"TLS.identifier", p.Identifier,
		"TLS.sdk_appid", p.SDKAppId,
		"TLS.time", p.Time,
		"TLS.expire_after", p.ExpireAfter,
	}
	for index, value := range pairs {
		buf.WriteString(value)
		if index%2 == 0 {
			buf.WriteRune(SEPERATOR)
		} else {
			buf.WriteRune(LINE_END)
		}
	}
	return buf.String()
}

func newTLSJSONPayload() tlsJSONPayload {
	return tlsJSONPayload{Version: "201512300000"}
}

type TLSSignature struct {
	Valid      bool
	UrlSig     string
	ExpireTime time.Time
	InitTime   time.Time
}

type CheckTLSSignatureResult struct {
	Valid       bool
	ExpireAfter time.Duration
	InitTime    time.Time
}

func (s *Signer) Sign(identifier string, expireAfter time.Duration) (TLSSignature, error) {
	payload := newTLSJSONPayload()
	payload.AppIDAt3rd = "0"
	payload.AccountType = s.accountTypeStr
	payload.Identifier = identifier
	payload.SDKAppId = s.appIdStr
	now := s.timeFunc()
	expireAt := now.Add(expireAfter)
	payload.Time = strconv.FormatInt(s.timeFunc().Unix(), 10)
	payload.ExpireAfter = strconv.Itoa(int(expireAfter.Seconds()))
	signStr := payload.ToSignString()
	sign, err := s.privKey.SignPKCS1v15(openssl.SHA256_Method, []byte(signStr))
	if err != nil {
		return TLSSignature{}, errors.Wrap(err, "error while calc signature")
	}
	signResultStr := base64.StdEncoding.EncodeToString(sign)
	payload.Signature = signResultStr
	jsonBytes, err := json.Marshal(&payload)
	if err != nil {
		return TLSSignature{}, errors.Wrap(err, "error while marshalling json")
	}
	compressed, err := zlibCompress(jsonBytes)
	escaped := signEncoding.EncodeToString(compressed)
	return TLSSignature{
		UrlSig: escaped, ExpireTime: expireAt, InitTime: now,
		Valid: true,
	}, nil
}

func (s *Signer) Verify(urlSig string, identifier string) (CheckTLSSignatureResult, error) {
	var err error
	decoded, err := b64decode(signEncoding, []byte(urlSig))
	if err != nil {
		return CheckTLSSignatureResult{}, err
	}
	decompressed, err := zlibDecompress(decoded)
	if err != nil {
		return CheckTLSSignatureResult{}, err
	}
	payload := tlsJSONPayload{}
	err = json.Unmarshal(decompressed, &payload)
	if err != nil {
		return CheckTLSSignatureResult{}, errors.Wrap(err, "error while decode json")
	}
	if payload.Identifier != identifier || payload.SDKAppId != s.appIdStr {
		return CheckTLSSignatureResult{}, nil
	}
	verifyStr := payload.ToSignString()
	rawSig, err := b64decode(base64.StdEncoding, []byte(payload.Signature))
	if err != nil {
		return CheckTLSSignatureResult{}, err
	}
	err = s.privKey.VerifyPKCS1v15(openssl.SHA256_Method, []byte(verifyStr), rawSig)
	if err != nil {
		return CheckTLSSignatureResult{}, nil
	}
	// time.Unix(int64(payload.Time), 0)
	// time.Second * time.Duration(payload.ExpireAfter)
	// TODO(QuantumGhost): convert time and duration string back to native type.
	check := CheckTLSSignatureResult{Valid: true}
	initTimeStamp, err := strconv.ParseInt(payload.Time, 10, 64)
	if err != nil {
		return check, errors.Wrap(err, "cannot parse timestamp string TLS.time")
	}
	check.InitTime = time.Unix(initTimeStamp, 0)
	expireAfterSec, err := strconv.ParseInt(payload.ExpireAfter, 10, 64)
	if err != nil {
		return check, errors.Wrap(err, "cannot parse duration string TLS.expire_after")
	}
	check.ExpireAfter = time.Duration(expireAfterSec) * time.Second
	return check, nil
}

type Signer struct {
	appId          int
	appIdStr       string
	accountType    int
	accountTypeStr string
	privKey        openssl.PrivateKey
	// for testing purpose, defaults to time.Now
	timeFunc func() time.Time
}

func NewSigner(appId int, accountType int, privKeyStr string) (*Signer, error) {
	privKey, err := openssl.LoadPrivateKeyFromPEM([]byte(privKeyStr))
	if err != nil {
		return nil, err
	}
	return &Signer{
		appId: appId, appIdStr: strconv.Itoa(appId),
		accountType: accountType, accountTypeStr: strconv.Itoa(accountType),
		privKey: privKey, timeFunc: time.Now}, nil
}

func zlibCompress(src []byte) ([]byte, error) {
	var err error
	buf := new(bytes.Buffer)
	buf.Grow(bufSize)
	writer := zlib.NewWriter(buf)

	if _, err = writer.Write(src); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func zlibDecompress(src []byte) ([]byte, error) {
	gzipReader, err := zlib.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, errors.Wrap(err, "error while creating gzip reader")
	}
	d, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, errors.Wrap(err, "error while decompress gzip")
	}
	return d, nil
}

func b64encode(encoding *base64.Encoding, src []byte) []byte {
	dst := make([]byte, encoding.EncodedLen(len(src)))
	encoding.Encode(dst, src)
	return dst
}

func b64decode(encoding *base64.Encoding, src []byte) ([]byte, error) {
	dst := make([]byte, encoding.DecodedLen(len(src)))
	n, err := encoding.Decode(dst, src)
	if err != nil {
		return nil, errors.Wrap(err, "error while base64 decoding")
	}
	return dst[:n], nil
}
