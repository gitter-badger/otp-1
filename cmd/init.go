// Copyright © 2016 Buck Brady <bbrady@mythic.tech>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mythic-tech/otp/logging"
	parse "github.com/mythic-tech/otp/parsers"
	"github.com/spf13/cobra"
	// init for sqlcipher
	_ "github.com/xeodou/go-sqlcipher"
)

const projectdir = ".otp/"

// InitData passwords yo
type InitData struct {
	Passwd string
}

var (
	force bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalize OTP database.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		logging.Debug("Initialization Requested...", DEBUG)
		logging.Debug("Requesting password...", DEBUG)
		fmt.Println("Please enter a password for encrypting the database.")
		pass, err := parse.GetPassword()
		fmt.Println("")
		if err != nil {
			logging.Debug("Failed to get password", DEBUG)
			log.Fatalf("Failed to get password: %s", err)
		}
		data := InitData{Passwd: pass}
		status := initOtp(&data)
		if status == false {
			fmt.Println("Failed to initalize otp. Use --debug to troubleshoot.")
		} else {
			fmt.Println("OTP database initalized. Use 'otp add' to add a key.")
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().BoolVar(&force, "force", false, "Force init. This will overwrite the current database.")

}

// Init initializes the db and program folder
func initOtp(data *InitData) (status bool) {
	var DBPath = ProjectPath + "otp.db"
	logging.Debug("Init() starting", DEBUG)

	logging.Debug("Calling createdir()", DEBUG)
	// Create ~/.otp/
	err := createdir()
	if err != nil {
		logging.Debug("createdir() failed. Exiting", DEBUG)
		log.Fatalf("Failed to create project dir (~/.otp): %s", err)
	}
	// Check if database exists
	logging.Debug("Checking if databases exists...", DEBUG)
	info, err := os.Stat(DBPath)
	logging.Debug("Stat info:", DEBUG)
	if DEBUG {
		log.Println(info)
		log.Println(err)
	}
	if err == nil {
		return false
	}

	logging.Debug("Creating Database...", DEBUG)
	// Create database
	err = createdb(data)
	if err != nil {
		logging.Debug("Failed to create database", DEBUG)
		log.Fatalf("Failed to create database (~/.otp/otp.db): %s", err)
	}

	return true
}

// create proper
func createdir() (err error) {
	logging.Debug("Getting user home dir...", DEBUG)

	if err != nil {
		logging.Debug("There was an error getting current user.", DEBUG)
		return err
	}
	logging.Debug("Create Project Dir: "+ProjectPath, DEBUG)
	err = os.MkdirAll(ProjectPath, 0700)
	if err != nil {
		logging.Debug("There was an error creating the project dir!", DEBUG)
		return err
	}
	logging.Debug("Project Dir Created!", DEBUG)

	return nil
}

// create database file and initialize password
func createdb(data *InitData) (err error) {
	dbpath := ProjectPath + "otp.db"
	logging.Debug("Database path: "+dbpath, DEBUG)
	logging.Debug("Encrypting Database with password: "+data.Passwd, DEBUG)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		logging.Debug("Error opening Database", DEBUG)
		return err
	}
	defer db.Close()

	q := fmt.Sprintf("PRAGMA key = '%s';", data.Passwd)
	_, err = db.Exec(q)
	if err != nil {
		logging.Debug("Failed to execute statement #1", DEBUG)
		return err
	}

	c := "CREATE TABLE `accounts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `nick` TEXT, `issuer` TEXT, `key` TEXT);"
	res, err := db.Exec(c)
	fmt.Sprintln(res)
	if err != nil {
		logging.Debug("Failed to execute statement #2", DEBUG)
		return err
	}

	logging.Debug("Fixing Permissions on otp.db", DEBUG)
	err = os.Chmod(dbpath, 0600)
	if err != nil {
		fmt.Printf("unable to fix permissions on %s\n", dbpath)
		fmt.Printf("to fix permission please run chmod 0600 %s", dbpath)
	}

	return nil
}
