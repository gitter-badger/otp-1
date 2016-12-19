package parse

import "os/user"

// GetHomeDir return the home dir of the current User
func GetHomeDir() (path string, err error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil

}
