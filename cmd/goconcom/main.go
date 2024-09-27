package main

import (
	"flag"
	"fmt"
	"github.com/roemer/goconcom/internal/app/goconcom"
	"os"
	"slices"
)

// Command Holds information about a CLI command that can be executed
type Command struct {
	Name string
	Help string
	Run  func(args []string) error
}

// The list of CLI commands
var commands = []Command{
	{Name: "help", Help: "Prints this help", Run: goconcom.HelpCmd},
	{Name: "run", Help: "Runs the goconcom process", Run: goconcom.RunCmd},
}

func main() {
	// CLI flags
	flag.Usage = printUsage
	flag.Parse()

	// A command need to be passed
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Command and the command arguments
	subCmd := flag.Arg(0)
	subCmdArgs := flag.Args()[1:]

	// Run the command
	runCommand(subCmd, subCmdArgs)
}

// Prints the base usage
func printUsage() {
	fmt.Fprintf(os.Stderr, "goconcom v%s", goconcom.Version)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  goconcom [flags] <command> [command flags]")
	fmt.Fprintln(os.Stderr, "")

	fmt.Fprintln(os.Stderr, "Commands:")
	for _, cmd := range commands {
		fmt.Fprintf(os.Stderr, "  %-8s %s\n", cmd.Name, cmd.Help)
	}

	// Uncomment if there are flags
	//fmt.Fprintln(os.Stderr, "Flags:")
	//fmt.Fprintln(os.Stderr, "")
	//flag.PrintDefaults()

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Run `goconcom <command> -h` to get help for a specific command\n\n")
}

// Tries to run the given command
func runCommand(name string, args []string) {
	cmdIdx := slices.IndexFunc(commands, func(cmd Command) bool {
		return cmd.Name == name
	})

	if cmdIdx < 0 {
		fmt.Fprintf(os.Stderr, "command \"%s\" not found\n\n", name)
		flag.Usage()
		os.Exit(1)
	}

	if err := commands[cmdIdx].Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
