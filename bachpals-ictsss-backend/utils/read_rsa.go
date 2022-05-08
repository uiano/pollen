package utils

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "io/ioutil"
)

func ReadPrivateKey() interface{} {
    priv, err := ioutil.ReadFile("keypair/rsa")
    if err != nil {
        return "No RSA private key found"
    }

    privPem, _ := pem.Decode(priv)
    var privPemBytes []byte
    if privPem.Type != "RSA PRIVATE KEY" {
        return "RSA private key is of the wrong type"
    }

    rsaPrivateKeyPassword := ""

    if rsaPrivateKeyPassword != "" {
        privPemBytes, err = x509.DecryptPEMBlock(privPem, []byte(rsaPrivateKeyPassword))
    } else {
        privPemBytes = privPem.Bytes
    }

    var parsedKey interface{}
    if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
        if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
            return "Unable to parse RSA private key"
        }
    }

    privateKey, ok := parsedKey.(*rsa.PrivateKey)
    if !ok {
        return "Unable to parse RSA private key, generating a temp one"
    }

    return privateKey
}