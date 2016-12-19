package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/adminfromhell/otp/logging"
	// init for sqlcipher
	_ "github.com/xeodou/go-sqlcipher"
)

const projectdir = ".otp/"

// DEBUG status
var DEBUG bool

// ProjectPath is injected from main
var ProjectPath string

// InitData is used to pass correct data Init
type InitData struct {
	Passwd string
}

// Init initializes the db and program folder
func Init(data *InitData) (status bool) {
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `accounts` (`id` INTEGER PRIMARY KEY, `accountName` char, `issuer` char, `key` chart)")
	if err != nil {
		logging.Debug("Failed to execute statement #2", DEBUG)
		return err
	}

	return nil
}
