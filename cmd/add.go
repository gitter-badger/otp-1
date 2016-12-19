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
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/mythic-tech/otp/logging"
	parse "github.com/mythic-tech/otp/parsers"
	"github.com/spf13/cobra"
	// sql cipher
	_ "github.com/xeodou/go-sqlcipher"
)

var (
	nickname string
	issuer   string
	secret   string
	password string
	update   bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new account to the database",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		if nickname == "" {
			fmt.Println("No Nickname specified. Will use issuers name.")
		}
		if issuer == "" {
			for len(strings.TrimSpace(issuer)) < 1 {
				fmt.Print("Enter an issuer name: ")
				issuer, _ = reader.ReadString('\n')
				fmt.Println("")
			}
			logging.Debug(issuer, DEBUG)
		}
		if secret == "" {
			for len(strings.TrimSpace(secret)) < 1 {
				fmt.Print("Enter secret: ")
				secret, _ = reader.ReadString('\n')
				fmt.Println("")
			}
			logging.Debug(secret, DEBUG)
		}
		if password == "" {
			password, _ = parse.GetPassword()
		}
		err := addKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("OTP account secuessfully added!")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringVar(&nickname, "nickname", "", "Give OTP account a nickname")
	addCmd.Flags().StringVar(&issuer, "issuer", "", "Note what provider the key is from.")
	addCmd.Flags().StringVar(&secret, "secret", "", "OTP secret from provider.")
	// addCmd.Flags().BoolVar(&update, "update", false, "Overide an already existing key")
}

func addKey() (err error) {
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
	if nickname == "" {
		nickname = "-"
	}
	q = fmt.Sprintf("INSERT INTO `accounts` (nick, issuer, key) values('%s', '%s', '%s');", nickname, issuer, secret)
	// fmt.Println(q)
	_, err = db.Exec(q)
	if err != nil {
		logging.Debug("Failed to insert", DEBUG)
		return err
	}
	return nil
}
