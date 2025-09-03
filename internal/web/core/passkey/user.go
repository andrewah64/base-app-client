package passkey

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	Id          []byte
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

func (user *User) WebAuthnID() []byte {
	return user.Id
}

func (user *User) WebAuthnName() string {
	return user.Name
}

func (user *User) WebAuthnDisplayName() string {
	return user.DisplayName
}

func (user *User) WebAuthnCredentials() []webauthn.Credential {
	return user.Credentials
}
