// Copyright © 2022 Yifan Huang

package root

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "Coeus",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:`,
	}
}

func UnaryCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "unary",
		Short:   "Send a unary gRPC call to target server",
		Example: `coeus unary localhost:8080`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			println(args[0])
		},
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
//func Execute() {
//
//	err := rootCmd.Execute()
//	if err != nil {
//		os.Exit(1)
//	}
//}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Coeus.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
