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

	"github.com/mythic-tech/otp/logging"
	parse "github.com/mythic-tech/otp/parsers"
	"github.com/spf13/cobra"
	// security
	_ "github.com/xeodou/go-sqlcipher"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list info about each account in the database",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if password == "" {
			password, _ = parse.GetPassword()
		}
		err := listKeys()
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listKeys() (err error) {
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
	r := "SELECT nick,issuer FROM accounts order by nick asc, issuer asc;"
	rows, err := db.Query(r)
	if err != nil {
		logging.Debug("Failed to execute statement #2", DEBUG)
		return err
	}
	defer rows.Close()

	fmt.Println("")
	fmt.Println("--- Account List ---")
	fmt.Printf("%-15s %s\n", "Nickname", "Issuer")
	for rows.Next() {
		var (
			nick   string
			issuer string
		)
		rows.Scan(&nick, &issuer)
		fmt.Printf("%-20s %s\n", nick, issuer)
	}
	fmt.Println("")

	return nil
}
