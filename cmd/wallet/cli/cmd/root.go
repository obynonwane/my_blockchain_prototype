package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	accountName string
	accountPath string
)

const (
	keyExtenstion = ".ecdsa"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Your simple wallet",
}

// definining default values even before main is executed
func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&accountName, "account", "a", "private.ecdsa", "The account to use.")
	rootCmd.PersistentFlags().StringVarP(&accountPath, "account-path", "p", "cmd/zblock/accounts/", "Path to the directory")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getPrivateKeyPath() string {
	if !strings.HasSuffix(accountName, keyExtenstion) {
		accountName += keyExtenstion
	}

	return filepath.Join(accountPath, accountName)
}
