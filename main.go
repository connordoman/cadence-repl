package main

import (
	"fmt"
	"os"

	"github.com/connordoman/cadence"
	"github.com/connordoman/cadence-repl/internal/repl"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:  "cadence",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		verboseFlag, _ := cmd.Flags().GetBool("verbose")
		humanReadableFlag, _ := cmd.Flags().GetBool("human-readable")

		if len(args) == 0 {
			r := repl.NewRepl(verboseFlag, humanReadableFlag)
			return r.Run()
		}

		input := os.Args[1]
		jsonFlag, _ := cmd.Flags().GetBool("json")
		if jsonFlag {
			json, err := cadence.CompileAsJSON(input)
			if err != nil {
				return err
			}
			fmt.Println(json)
			return nil
		}

		results, err := cadence.Compile(input)
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
	RootCmd.Flags().BoolP("json", "j", false, "output as JSON")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
