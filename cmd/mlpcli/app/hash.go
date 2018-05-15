// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"github.com/gopenguin/minimal-ldap-proxy/pkg/password"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"os"
)

// passwordCmd represents the password command
var passwordCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash a user password",
	Long: `Hashes a user password from STDIN and outputs to STDOUT.
It uses the argon2 algorithm by default`,
	Run: func(cmd *cobra.Command, args []string) {
		pw, _ := gopass.GetPasswdPrompt("Enter password: ", false, os.Stdin, os.Stderr)

		hash, err := password.Hash(string(pw))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to hash the password: %v", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "%s\n", hash)
	},
}

func init() {
	RootCmd.AddCommand(passwordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// passwordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// passwordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
