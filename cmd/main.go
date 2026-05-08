package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/connordoman/cadence/internal/compiler"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:  "cadence",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		verboseFlag, _ := cmd.Flags().GetBool("verbose")
		humanReadableFlag, _ := cmd.Flags().GetBool("human-readable")

		if len(args) == 0 {
			repl(ReplOptions{
				Verbose:       verboseFlag,
				HumanReadable: humanReadableFlag,
			})

			return nil
		}

		input := os.Args[1]

		comp := compiler.NewCompiler(input)
		results, err := comp.Compile(verboseFlag)
		if err != nil {
			return err
		}

		for _, result := range results {
			fmt.Println(result.Format("2006-01-02"))
		}

		return nil
	},
}

func init() {
	RootCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	RootCmd.Flags().BoolP("human-readable", "r", false, "human readable output")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type ReplOptions struct {
	Verbose       bool
	HumanReadable bool
}

func repl(options ReplOptions) {
	if options.Verbose {
		log.Println("Running in verbose mode")
	}

	fmt.Println("Cadence REPL (type 'exit' to quit)")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		}
		comp := compiler.NewCompiler(line)
		results, err := comp.Compile(options.Verbose)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		for _, result := range results {
			if options.HumanReadable {
				fmt.Println(result.Format("Monday, 02 January 2006"))
			} else {
				fmt.Println(result.Format("2006-01-02"))
			}
		}
	}
	fmt.Println("Bye!")
}
