package servicenow

import (
	"encoding/base64"
	"errors"
	"os"
)

// Credentials store the username and password to access ServiceNow
type Credentials struct {
	Username string
	Password string
}

// GetCredentialsFromEnvironment pulls ServiceNow credentials from SN_USER and SN_PASSWORD environment variables
func GetCredentialsFromEnvironment() (Credentials, error) {
	user, userOk := os.LookupEnv("SN_USER")
	encoded_password, passwordOk := os.LookupEnv("SN_PASSWORD_BASE64")

	byte_password, _ := base64.StdEncoding.DecodeString(encoded_password)
	password := string(byte_password)
	if !userOk || !passwordOk {
		return Credentials{"", ""}, errors.New("Please specify serviceNow username and password by setting environment variables: SN_USER & SN_PASSWORD")
	}

	return Credentials{user, password}, nil
}
