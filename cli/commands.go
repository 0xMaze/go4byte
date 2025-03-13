package cli

import (
	fnprocessor "fourbyte/fnprocessor"

	"github.com/spf13/cobra"
)

type commandConfig struct {
	sig      fnprocessor.FnSig
	abiFlag  bool
	hashFlag bool
}

var (
	rootCmd = &cobra.Command{
		Use:   "fbyte",
		Short: "fbyte is a tool to generate function ABI and calculate 4-byte signatures",
		Run:   runCommand,
	}

	cfg commandConfig
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&cfg.abiFlag, "abi", false, "Generate function's ABI")
	rootCmd.Flags().BoolVar(&cfg.hashFlag, "hash", false, "Calculate 4-byte signature of a function")
	rootCmd.Flags().VarP(&cfg.sig, "sig", "s", "Solidity function signature")
}
