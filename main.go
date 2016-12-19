package main

import (
	"fmt"
	"log"

	"github.com/adminfromhell/otp/cmd"
	"github.com/adminfromhell/otp/logging"
	parse "github.com/adminfromhell/otp/parsers"
	flag "github.com/ogier/pflag"
)

// DEBUG status
var DEBUG bool

// Init status
var Init bool

func init() {
	flag.BoolVar(&Init, "init", false, "Initialize OTP database")
	flag.BoolVar(&DEBUG, "debug", false, "Turn on debugging messages")
	flag.Parse()
	logging.Debug("====   Running in Debug Mode   ====", DEBUG)
	logging.Debug("Flags Parsed", DEBUG)
}

func main() {
	// Pass debugging status to cmd's
	cmd.DEBUG = DEBUG
	homedir, _ := parse.GetHomeDir()
	cmd.ProjectPath = homedir + "/.otp/"

	// check if we need to init database
	if Init {

		logging.Debug("Initialization Requested...", DEBUG)
		logging.Debug("Requesting password...", DEBUG)
		fmt.Println("Please enter a password for encrypting the database.")
		pass, err := parse.GetPassword()
		if err != nil {
			logging.Debug("Failed to get password", DEBUG)
			log.Fatalf("Failed to get password: %s", err)
		}
		data := cmd.InitData{Passwd: pass}
		cmd.Init(&data)
	}

}
