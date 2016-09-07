// tokens provide generating and validaing JSON Web Tokens
package tokens

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/optiopay/go-oidc/jose"
)

type Signer struct {
	rsa *jose.SignerRSA
}

// NewSigner creates a new Signer which can be used to create JWT.
//
// It takes a unique id for a private-public key pair.
func NewSigner(id string, privateKey io.Reader) (*Signer, error) {
	raw, err := ioutil.ReadAll(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot read private key: %s", err)
	}
	key, err := x509.ParsePKCS1PrivateKey(raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse private key: %s", err)
	}
	rsa := jose.NewSignerRSA(id, *key)
	return &Signer{rsa}, nil
}

// Generate creates a Encoded JSON Web Token for the given claims.
func (s *Signer) Generate(claims map[string]interface{}) (token string, err error) {
	jwt, err := jose.NewSignedJWT(jose.Claims(claims), s.rsa)
	if err != nil {
		return "", fmt.Errorf("could not create JWT: %s", err)
	}
	return jwt.Encode(), nil
}

type Verifier struct {
	rsa *jose.VerifierRSA
}

// NewVerifier creates a new Verifier.
//
// It takes a unique id for the given private-public key pair. It should
// be the same ID used when for signing tokens using the Signer.
func NewVerifier(id string, publicKey io.Reader) (*Verifier, error) {
	raw, err := ioutil.ReadAll(publicKey)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %s", err)
	}
	interm, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %s", err)
	}
	key, ok := interm.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot handle key type %T", interm)
	}
	rsa := &jose.VerifierRSA{
		KeyID:     id,
		PublicKey: *key,
		Hash:      crypto.SHA256,
	}
	return &Verifier{rsa}, nil
}

// Parse validates the provided token and decodes the claims encoded in it.
func (v *Verifier) Parse(token string) (claims map[string]interface{}, err error) {
	jwt, err := jose.ParseJWT(token)
	if err != nil {
		return nil, fmt.Errorf("could not parse JWT: %s", err)
	}
	err = v.rsa.Verify(jwt.Signature, []byte(jwt.Data()))
	if err != nil {
		return nil, fmt.Errorf("could not verify signature: %s", err)
	}
	claims, err = jwt.Claims()
	if err != nil {
		return nil, fmt.Errorf("could not decode claims: %s", err)
	}
	return claims, nil
}
