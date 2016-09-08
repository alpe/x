// tokens provide generating and validaing JSON Web Signatures
package tokens

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/square/go-jose.v1"
)

// Signer helps create JSON Web Signatures for any payload.
type Signer struct {
	rsa jose.Signer
}

// NewSigner creates a Signer.
//
// The provided private key must be generated using RSA.
func NewSigner(priv io.Reader) (*Signer, error) {
	raw, err := ioutil.ReadAll(priv)
	if err != nil {
		return nil, fmt.Errorf("cannot read private key: %s", err)
	}
	pk, err := jose.LoadPrivateKey(raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse private key: %s", err)
	}
	s, err := jose.NewSigner(jose.PS256, pk)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %s", err)
	}
	return &Signer{s}, nil
}

// Generate creates an Encoded JSON Web Signature for the given payload.
//
// The payload must be a pointer to the interface and must be JSON Marshallable.
func (s *Signer) Generate(payload interface{}) (token string, err error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("could not marshal payload")
	}
	jws, err := s.rsa.Sign(raw)
	if err != nil {
		return "", fmt.Errorf("could not sign payload: %s", err)
	}
	token, err = jws.CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("could not serialize: %s", err)
	}
	return token, nil
}

// Verifier helps decode signed tokens.
type Verifier struct {
	pub *rsa.PublicKey
}

// NewVerifier creates a new Verifier given the public key.
//
// The public key must be of a RSA private-public key pair.
func NewVerifier(pub io.Reader) (*Verifier, error) {
	raw, err := ioutil.ReadAll(pub)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %s", err)
	}
	interm, err := jose.LoadPublicKey(raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %s", err)
	}
	key, ok := interm.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot handle key type %T", interm)
	}
	return &Verifier{key}, nil
}

// Parse validates the token and extracts the payload from it.
//
// The payload must be a pointer to the interface and must be JSON Unmarshallable.
func (v *Verifier) Parse(token string, payload interface{}) (err error) {
	jws, err := jose.ParseSigned(token)
	if err != nil {
		return fmt.Errorf("cannot parse token: %s", err)
	}
	raw, err := jws.Verify(v.pub)
	if err != nil {
		return fmt.Errorf("could not verify token: %s", err)
	}
	err = json.Unmarshal(raw, payload)
	if err != nil {
		return fmt.Errorf("could not decode payload: %s", err)
	}
	return nil
}
