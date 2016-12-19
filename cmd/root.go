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
	"strings"
	"time"

	otp "github.com/hgfischer/go-otp"
	"github.com/mythic-tech/otp/logging"
	parse "github.com/mythic-tech/otp/parsers"
	"github.com/spf13/cobra"
	// blank import
	_ "github.com/xeodou/go-sqlcipher"
)

// var cfgFile string
// arg variables
var (
	DEBUG bool
)

// Global vars
var (
	LOG *log.Logger
)

// ProjectPath is global path
var ProjectPath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "otp",
	Short: "Generate OTP codes for stored accounts",
	// Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if password == "" {
			password, _ = parse.GetPassword()
			fmt.Println("")
		}
		err := getOTP(password)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	homedir, _ := parse.GetHomeDir()
	ProjectPath = homedir + "/.otp/"
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.otp.yaml)")
	RootCmd.PersistentFlags().StringVar(&password, "password", "", "Optional way to pass the password to unlock the database. WARNING: INSECURE!")
	RootCmd.PersistentFlags().BoolVar(&DEBUG, "debug", false, "Turn on debugging messages")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// setup the logger for DEBUG
	// LOG = log.New(os.Stdout, "", log.Ltime)
}

// LEAVING COMMENTED OUT UNTIL CONFIG IS NEEDED
// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" { // enable ability to specify config file via flag
// 		viper.SetConfigFile(cfgFile)
// 	}
//
// 	viper.SetConfigName(".otp")  // name of config file (without extension)
// 	viper.AddConfigPath("$HOME") // adding home directory as first search path
// 	viper.AutomaticEnv()         // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }

func getOTP(pass string) (err error) {
	dbpath := ProjectPath + "otp.db"
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		logging.Debug("Failed to Open db", DEBUG)
		return err
	}
	defer db.Close()

	q := fmt.Sprintf("PRAGMA key = '%s';", password)
	// fmt.Println(q)
	_, err = db.Exec(q)
	if err != nil {
		logging.Debug("Failed to execute statement #1", DEBUG)
		return err
	}

	q = "SELECT nick,issuer,key FROM accounts ORDER BY nick asc, issuer asc;"
	rows, err := db.Query(q)
	if err != nil {
		logging.Debug("Failed to execute statement #2", DEBUG)
		return err
	}
	defer rows.Close()

	fmt.Printf("%s\t%s\t%s\n", "Nickname", "Valid For", "OTP Code")
	fmt.Println(strings.Repeat("-", 40))

	for rows.Next() {
		var (
			nick   string
			issuer string
			key    string
		)
		rows.Scan(&nick, &issuer, &key)
		genCode(nick, issuer, key)
	}
	fmt.Println("")
	// fmt.Println("")
	return nil
}

func genCode(nick string, issuer string, key string) {
	key = strings.ToUpper(key)
	totp := &otp.TOTP{Secret: key, IsBase32Secret: true}
	token := totp.Get()
	// issuer = strings.Replace(issuer, "\n", "", -1)
	fmt.Printf("%s\t%02d seconds\t%-15s\n", nick, (30 - time.Now().Unix()%30), token)
}
