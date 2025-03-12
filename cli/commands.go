package cli

import (
	"fmt"
	fnprocessor "fourbyte/fnprocessor"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fbyte",
		Short: "fbyte is a tool to generate function ABI and calculate 4-byte signatures",
		Run: func(cmd *cobra.Command, args []string) {
			if abiFlag || hashFlag {
				if sig.IsEmpty() {
					fmt.Println("Function signature was not provided")
					return
				}

				type task struct {
					run    func() (string, error)
					header string
					errMsg string
				}

				var tasks []task
				if abiFlag {
					tasks = append(tasks, task{
						run:    sig.GenerateABI,
						header: "ABI:",
						errMsg: "Error generating ABI:",
					})
				}
				if hashFlag {
					tasks = append(tasks, task{
						run:    sig.FourBytes,
						header: "\n4-Byte Signature:",
						errMsg: "Error calculating four-byte hash:",
					})
				}

				var output strings.Builder
				for _, t := range tasks {
					result, err := t.run()
					if err != nil {
						fmt.Printf("%s %v\n", t.errMsg, err)
						os.Exit(1)
					}
					output.WriteString(t.header)
					output.WriteString("\n")
					output.WriteString(result)
					output.WriteString("\n")
				}

				fmt.Println(strings.TrimSuffix(output.String(), "\n"))
			} else {
				fmt.Println("No flag provided. Use --help for usage information.")
			}
		},
	}

	abiFlag, hashFlag bool
	sig               fnprocessor.FnSig
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&abiFlag, "abi", false, "Generate function's ABI")
	rootCmd.Flags().BoolVar(&hashFlag, "hash", false, "Calculate 4-byte signature of a function")
	rootCmd.Flags().VarP(&sig, "sig", "s", "Solidity function signature (e.g., \"function approve(address spender, uint256 amount) external returns (bool)\")")
}
