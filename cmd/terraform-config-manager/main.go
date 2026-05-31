package main

import (
	"flag"
	"fmt"
	"os"

	"terraform-config-manager/internal/scaffold"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: terraform-config-manager <command> [flags]")
		fmt.Fprintln(os.Stderr, "Commands: scaffold")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scaffold":
		scaffoldCmd := flag.NewFlagSet("scaffold", flag.ExitOnError)
		project := scaffoldCmd.String("project", "", "GCP project ID (required)")
		outputDir := scaffoldCmd.String("output-dir", ".", "directory to create the project folder in")
		scaffoldCmd.Parse(os.Args[2:])

		if *project == "" {
			fmt.Fprintln(os.Stderr, "Error: -project is required")
			scaffoldCmd.Usage()
			os.Exit(1)
		}

		if err := scaffold.Run(*project, *outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
