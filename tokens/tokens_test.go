package tokens

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
	mrand "math/rand"
	"strings"
	"testing"
	"time"
)

func TestTokenGenerator(t *testing.T) {
	priv, pub, err := generateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSigner("testkey", priv)
	if err != nil {
		t.Fatalf("could not create Signer: %s", err)
	}
	v, err := NewVerifier("testKey", pub)
	if err != nil {
		t.Fatalf("could not create Verifier: %s", err)
	}
	for i, test := range []map[string]interface{}{
		{
			"key1": "value1",
			"key2": "value2",
		},
		{
			"key1": true,
			"key2": "value2",
		},
		{
			"key1": float64(time.Now().UnixNano()), // json uses float for all numbers
			"key2": "value2",
		},
	} {
		token, err := s.Generate(test)
		if err != nil {
			t.Errorf("test %d failed: %s", i, err)
		}
		claims, err := v.Parse(token)
		if err != nil {
			t.Errorf("test %d failed to decoce claims: %s", i, err)
		}
		for key, value := range test {
			claim, ok := claims[key]
			if !ok {
				t.Errorf("test %d failed to find claim: %s", i, key)
			}
			if claim != value {
				t.Errorf("test %d failed, claims not same: expected %s, got %s", i, value, claim)
			}
		}
	}
}

func TestModifiedToken(t *testing.T) {
	priv, pub, err := generateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSigner("testkey", priv)
	if err != nil {
		t.Fatalf("could not create Signer: %s", err)
	}
	v, err := NewVerifier("testKey", pub)
	if err != nil {
		t.Fatalf("could not create Verifier: %s", err)
	}
	for i, test := range []map[string]interface{}{
		{
			"key1": "value1",
			"key2": "value2",
		},
		{
			"key1": true,
			"key2": "value2",
		},
		{
			"key1": float64(time.Now().UnixNano()), // json uses float for all numbers
			"key2": "value2",
		},
	} {
		token, err := s.Generate(test)
		if err != nil {
			t.Errorf("test %d failed: %s", i, err)
		}
		token = replaceRandomRune(token)
		claims, err := v.Parse(token)
		if err == nil {
			t.Errorf("expected test %d to fail, %v", i, claims)
		}
	}
}

func generateKeyPair() (priv, pub io.Reader, err error) {
	pk, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate private key: %s", err)
	}
	priv = bytes.NewReader(x509.MarshalPKCS1PrivateKey(pk))
	pubBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate public key: %s", err)
	}
	pub = bytes.NewReader(pubBytes)
	return priv, pub, err
}

func replaceRandomRune(in string) string {
	// the jws token is serialized as header.payload.signature
	// so modifying the payload would be a good place to check
	payloadStart := strings.Index(in, ".")
	if payloadStart == -1 {
		panic("incorrect token")
	}
	payloadEnd := strings.LastIndex(in, ".")
	if payloadEnd == -1 {
		panic("incorrect token")
	}
	// we need a random position in the payload
	n := mrand.Intn(payloadEnd-payloadStart) + payloadStart + 1
	out := []rune(in)
	out[n] = out[n] + 5
	return string(out)
}
