package signature

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"testing"
	"testing/quick"
)

const (
	testingPrivKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIMlqlc7hHlhjN/AjcSNg6z8zS2oL8icinZpCWE1qaWf4oAoGCCqGSM49
AwEHoUQDQgAElmfLMxM6p6Puo4Sr445c8IOmPgsO/Obndw3ue4cb7JsDxPeGi8rk
FmdVRXewk+3sqiemPn3bHVP6TNUnWUo+RA==
-----END EC PRIVATE KEY-----`
	testIter = 1000
)

var testRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Test_zlibDecompress(t *testing.T) {
	const expectedStr = "Hello，世界😂"
	var compressed = []byte{120, 156, 243, 72, 205, 201, 201, 127, 191, 167, 231, 201, 142, 105, 207, 167, 246, 124, 152, 63, 163, 9, 0, 97, 222, 11, 15}
	decompressed, err := zlibDecompress(compressed)
	assert.NoError(t, err)
	assert.Equal(t, expectedStr, string(decompressed))
}

func Test_compressAndDecompress(t *testing.T) {
	for i := 0; i < testIter; i++ {
		value, ok := quick.Value(reflect.TypeOf(""), testRand)
		str := value.Interface().(string)
		assert.True(t, ok)
		compressed, err := zlibCompress([]byte(str))
		assert.NoError(t, err)
		decompressed, err := zlibDecompress(compressed)
		assert.NoError(t, err)
		assert.Equal(t, str, string(decompressed))
	}
}

func Test_signEncoding(t *testing.T) {
	var err error
	for i := 0; i < testIter; i++ {
		src := make([]byte, 100)
		_, err = rand.Read(src)
		assert.NoError(t, err)
		replacer := strings.NewReplacer("+", "*", "/", "-", "=", "_")
		b64Str := signEncoding.EncodeToString(src)
		assert.NoError(t, err)
		stdB64Str := base64.StdEncoding.EncodeToString(src)
		assert.Equal(t, replacer.Replace(stdB64Str), b64Str)
	}
}

func TestSigner_SignAndVerify(t *testing.T) {
	signer, err := NewSigner(1, 2, testingPrivKey)
	assert.NoError(t, err)
	for i := 0; i < testIter; i++ {
		value, ok := quick.Value(reflect.TypeOf(""), testRand)
		assert.True(t, ok)
		str := value.Interface().(string)
		sig, err := signer.Sign(str, 720 * time.Hour)
		assert.NoError(t, err)
		assert.True(t, sig.Valid)
		verify, err := signer.Verify(sig.UrlSig, str)
		assert.NoError(t, err)
		assert.True(t, verify.Valid)
	}
}

func Test_base64EncodeAndDecode(t *testing.T) {
	for i := 0; i < testIter; i++ {
		value, ok := quick.Value(reflect.TypeOf(""), testRand)
		assert.True(t, ok)
		str := value.Interface().(string)
		encoded := b64encode(base64.StdEncoding, []byte(str))
		decoded, err := b64decode(base64.StdEncoding, encoded)
		assert.NoError(t, err)
		assert.Equal(t, str, string(decoded))
	}
}

func TestSigner_Verify(t *testing.T) {
	const (
		testingAccountType    = 1532909569
		testingAppId          = 1400167151
		testingIdentifier     = "tencentcloud-im-test-user"
		expectedInitTimestamp = 1543885094
		expectedExpireAfter   = 3600
		testingSig = `eJw1jl1PgzAUhv8Lt4oW2jJq4gVhm1mEZHPDTGJCgJZ5lI8GCroZ-7sdY*fuPM978p5fYxds79I8b-paJeoohfFgWBTbDDHqMON29MBFraAA0WqrRJ3rNS*bnptQmUp0yuw77S7hVErgSaoS3HIdRxPu*FcyqnMBQchyZha1Jil*JLQiSQs1VmAHXc8G0XbQ1BraSOdtjM4zSQXV5V*CXZciRq5dcNA4XET*ytt7c29eFlsrZu-3zlAccB5G2cd*pl5f*Gd2s3Or02ZNWVw*e7DwkGTQMCGXJAj84LShZeMc8be-Uu6yJHHcZSQa3PVb*IQejb9-g2hf9A__`
		testingSig2 = `eJw1jl1vgkAQRf-Lvra0CwvINumDaBO1YmvQ*JEmZGEH3CofhYUoTf97V8R5m3Pu5M4vWs39JxZFeZ3JQF4KQC9It4hBMbVsih47LzhkUsQCSmUlZJFao1Nec02kmoRKanWl3C3MikLwgMmAlFzFcY8rfgw6dS0wMdbtgW7pvYRzIUoIWCy7CmLj*1kDZSXyTEEDq7xB8HV6KUV6*9ckjmNhat67RKKw97YeTZej2t89jA80bIfbjQAnTNLPhhaLr2db5z*bldu2bFLNgEfb9XKazJrjOKVm85GfHWYa3*nEduOQHNzLnoPFT7LZe*27P1jMvVf09w-dGGHE`
	)
	signer, err := NewSigner(testingAppId, testingAccountType, testingPrivKey)
	assert.NoError(t, err)
	for i, urlSig := range []string{testingSig, testingSig2} {
		check, err := signer.Verify(urlSig, testingIdentifier)
		assert.NoError(t, err, "verify error for case %d", i)
		assert.True(t, check.Valid, "sig not valid for case %d", i)
		assert.Equal(t, int64(expectedInitTimestamp), check.InitTime.Unix(),
			"InitTime incorrect for case %d", i)
		assert.Equal(t, time.Duration(expectedExpireAfter) * time.Second, check.ExpireAfter,
			"ExpireAfter incorrect for case %d", i)
	}
}
