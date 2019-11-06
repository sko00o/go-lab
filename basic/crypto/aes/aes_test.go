package aes

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAESGCM(t *testing.T) {
	assert := require.New(t)

	var key AES128Key
	keyBS, err := hex.DecodeString("3736db93cfc750d626504fb58be91575")
	assert.NoError(err)
	assert.Equal(len(key), len(keyBS))
	copy(key[:], keyBS)

	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	assert.NoError(err)

	oriData := []byte("secret data")

	// encode
	encData, err := Encrypt(key, oriData, nonce)
	assert.NoError(err)

	// decode
	decData, err := Decrypt(key, encData, nonce)
	assert.NoError(err)

	assert.Equal(oriData, decData)
}
