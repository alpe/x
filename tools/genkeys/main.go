package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	name := flag.String("name", "id_rsa", "names for the key pair created")
	flag.Parse()

	if err := genkeys(*name); err != nil {
		log.Fatal(err)
	}
}

func genkeys(name string) error {
	pk, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("could not generate private key: %s", err)
	}
	priv := x509.MarshalPKCS1PrivateKey(pk)
	err = ioutil.WriteFile(name, priv, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write to '%s': %s", name, err)
	}

	pub, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return fmt.Errorf("could not generate public key: %s", err)
	}
	err = ioutil.WriteFile(name+".pub", pub, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write to '%s': %s", name+".pub", err)
	}
	return nil
}
