package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRSAOAEP(t *testing.T) {
	assert := require.New(t)
	rng := rand.Reader

	// key generate
	key, err := rsa.GenerateKey(rng, 2048)
	assert.NoError(err)
	privateKey := key
	publicKey := &key.PublicKey

	hash := sha256.New()
	oriData := []byte("secret data")
	label := []byte("orders")

	ciphertext, err := rsa.EncryptOAEP(hash, rng, publicKey, oriData, label)
	assert.NoError(err)

	plaintext, err := rsa.DecryptOAEP(hash, rng, privateKey, ciphertext, label)
	assert.NoError(err)

	assert.Equal(oriData, plaintext)
}

func TestRSAPKCS1v15(t *testing.T) {
	assert := require.New(t)
	rng := rand.Reader

	// key generate
	key, err := rsa.GenerateKey(rng, 2048)
	assert.NoError(err)
	privateKey := key
	publicKey := &key.PublicKey

	hashed := sha256.Sum256([]byte("small message to be signed"))

	signature, err := rsa.SignPKCS1v15(rng, privateKey, crypto.SHA256, hashed[:])
	assert.NoError(err)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	assert.NoError(err)
}
