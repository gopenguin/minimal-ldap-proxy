// Copyright © 2018 gopenguin
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

package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"database/sql"
	"github.com/gopenguin/minimal-ldap-proxy/pkg"
	"github.com/gopenguin/minimal-ldap-proxy/types"
	"github.com/gopenguin/minimal-ldap-proxy/util"
	"os/signal"
	"strings"
	"syscall"
	"crypto/tls"
)

var (
	cfgFile string

	cmdConfig types.CmdConfig
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "minimal-ldap-proxy",
	Short: "Proxy ldap authentication requests to a database backend",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRunE: func(cmd *cobra.Command, args []string) error {
		jww.SetStdoutThreshold(jww.LevelInfo)

		err := loadConfig()
		if err != nil {
			return err
		}

		if !util.ContainsString(sql.Drivers(), cmdConfig.Driver) {
			return fmt.Errorf("%s is not one of the supported drivers: %s", cmdConfig.Driver, strings.Join(sql.Drivers(), ", "))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		backend, err := pkg.NewBackend(cmdConfig.Driver, cmdConfig.Conn, cmdConfig.AuthQuery, cmdConfig.SearchQuery)
		if err != nil {
			jww.ERROR.Fatalf("Error configuring backend: %v", err)
		}

		cert, err := tls.LoadX509KeyPair(cmdConfig.Cert, cmdConfig.Key)
		if err != nil {
			jww.ERROR.Fatalf("Error loading tls certificate: %v", err)
		}

		frontend := pkg.NewFrontend(cmdConfig.ServerAddress, cert, cmdConfig.BaseDn, cmdConfig.Rdn, cmdConfig.Attributes, backend)

		frontend.Serve()

		// When CTRL+C, SIGINT and SIGTERM signal occurs
		// Then stop server gracefully
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		close(ch)

		frontend.Stop()
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
	cobra.OnInitialize(initConfig)

	RootCmd.Flags().String("serverAddress", "127.0.0.1:1636", "the address to listen on")
	RootCmd.Flags().String("cert", "", "a pem encoded certificate")
	RootCmd.Flags().String("key", "", "a pem encoded certificate key")

	RootCmd.Flags().String("driver", "", fmt.Sprintf("the sql driver to use (%s)", strings.Join(sql.Drivers(), ", ")))
	RootCmd.Flags().String("conn", "", "the connection string")
	RootCmd.Flags().String("authQuery", "", "a sql query to retrieve the password by the username. The username is passed a the first parameter. The query must return one field, the password")
	RootCmd.Flags().String("searchQuery", "", "a sql query to retrieve the user attributes. This string should contain one %s for the projection and one ? for the selection")
	RootCmd.Flags().String("rdn", "", "the rdn of the user")
	RootCmd.Flags().String("baseDn", "", "the base dn for users")
	RootCmd.Flags().StringSlice("attributes", nil, "the attributes supported by the query provided to the backend backend (format: 'attr1,attr2,attr3,...')")

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/minimal-ldap-proxy.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("minimal-ldap-proxy") // name of config file (without extension)
	viper.AddConfigPath(".")                  // adding home directory as first search path
	viper.SetEnvPrefix("LDAP_PROXY")          // set the prefix for environment variables
	viper.AutomaticEnv()                      // read in environment variables that match

	flags := []string{
		"serverAddress",
		"driver",
		"conn",
		"authQuery",
		"searchQuery",
		"rdn",
		"baseDn",
		"attributes",
		"cert",
		"key",
	}

	for _, flag := range flags {
		viper.BindPFlag(flag, RootCmd.Flags().Lookup(flag))
	}
}

func loadConfig() error {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		jww.INFO.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cmdConfig); err != nil {
		return fmt.Errorf("unable to unmarshal config: %v", err)
	}

	jww.INFO.Printf("Configuration: %+v", cmdConfig)

	return nil
}
