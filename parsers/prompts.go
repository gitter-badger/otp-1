package parse

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

// GetPassword gets a password from the terminal without local echo
func GetPassword() (passwd string, err error) {
	fmt.Print("Enter Password: ")
	pass, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	// make sure we arnt given back a 0 length password
	for len(string(pass)) < 3 {
		fmt.Println("Bad entry. Please input your password again.")
		fmt.Print("Enter Password: ")
		pass, err = terminal.ReadPassword(0)
		if err != nil {
			return "", err
		}
	}
	passwd = string(pass)

	return passwd, nil
}
