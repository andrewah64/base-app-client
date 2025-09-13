package saml2

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"time"
)

func GenCert(spcCnNm string, spcOrgNm [] string, ku x509.KeyUsage, fromTs time.Time, spcExpTs time.Time) ([] byte, [] byte, error) {
	rsaKey, rsaKeyErr := rsa.GenerateKey(rand.Reader, 4096)
	if rsaKeyErr != nil {
		return nil, nil, rsaKeyErr
	}

	tmpl := &x509.Certificate {
		Subject : pkix.Name {
			CommonName   : spcCnNm,
			Organization : spcOrgNm,
        	},
		NotBefore : fromTs,
		NotAfter  : spcExpTs,
		KeyUsage  : ku,
		ExtKeyUsage: []x509.ExtKeyUsage {
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
	}

	derCert, derCertErr := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &rsaKey.PublicKey, rsaKey)
	if derCertErr != nil {
		return nil, nil, derCertErr
	}

	pemCert := pem.EncodeToMemory(
		&pem.Block{
			Type  : "CERTIFICATE",
			Bytes : derCert,
		},
	)

	pemCertBytes, _ := pem.Decode(pemCert)

	derRsaKey, derRsaKeyErr := x509.MarshalPKCS8PrivateKey(rsaKey)
	if derRsaKeyErr != nil {
		return nil, nil, derRsaKeyErr
	}

	pemRsaKey := pem.EncodeToMemory(
		&pem.Block{
			Type  : "RSA PRIVATE KEY",
			Bytes : derRsaKey,
		},
	)

	pemRsaKeyBytes, _ := pem.Decode(pemRsaKey)

	return pemCertBytes.Bytes, pemRsaKeyBytes.Bytes, nil
}
