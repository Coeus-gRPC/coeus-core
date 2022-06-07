// Copyright © 2022 Yifan Huang

package main

import (
	"Coeus/cmd/root"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmd := &cobra.Command{Use: "coeus"}
	cmd.AddCommand(
		root.UnaryCmd(),
		root.RootCmd(),
	)

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
