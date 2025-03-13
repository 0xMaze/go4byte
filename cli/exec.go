package cli

import (
	"fbyte/cli/export"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func runCommand(cmd *cobra.Command, args []string) {
	if cfg.abiFlag || cfg.hashFlag {
		if cfg.sig.IsEmpty() {
			fmt.Println("Function signature was not provided")
			return
		}

		var initFns []export.ExportOptFunc
		if exp {
			initFns = append(initFns, export.WithExport)
		}
		if !(outStr == "") {
			initFns = append(initFns, export.WithExportPath(outStr))
		}

		exportOps := export.NewExportOpts(initFns...)

		output, err := executeTasks(exportOps)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(output)
	} else {
		fmt.Println("No flag provided. Use --help for usage information.")
	}
}

func executeTasks(exportOps export.ExportOptions) (string, error) {
	tasks := createTasks(exportOps)
	return processTasks(tasks)
}

type task struct {
	run    func() (string, error)
	header string
	errMsg string
}

func createTasks(exportOpts export.ExportOptions) []task {
	var tasks []task
	if cfg.abiFlag {
		tasks = append(tasks, task{
			run: func() (string, error) {
				fmt.Printf("export: %v\n", exportOpts.Path)
				return cfg.sig.GenerateABI(exportOpts)
			},
			header: "ABI:",
			errMsg: "Error generating ABI:",
		})
	}
	if cfg.hashFlag {
		tasks = append(tasks, task{
			run:    cfg.sig.FourBytes,
			header: "\n4-Byte Signature:",
			errMsg: "Error calculating four-byte hash:",
		})
	}

	return tasks
}

func processTasks(tasks []task) (string, error) {
	var output strings.Builder
	for _, t := range tasks {
		result, err := t.run()
		if err != nil {
			return "", fmt.Errorf("%s %w", t.errMsg, err)
		}
		output.WriteString(t.header)
		output.WriteString("\n")
		output.WriteString(result)
		output.WriteString("\n")
	}
	return strings.TrimSuffix(output.String(), "\n"), nil
}
