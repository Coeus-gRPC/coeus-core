// Copyright © 2022 Yifan Huang

package coeus

import (
	"os"

	"github.com/Coeus-gRPC/coeus-core/internal/app"
	"github.com/spf13/cobra"
)

var (
	configFile    string
	config        = app.CoeusConfig{TotalCallNum: 1, Concurrent: 1, Insecure: false, Timeout: -1, TargetHost: "api.coeustool.dev:443"}
	runtimeConfig = app.CoeusRuntimeConfig{}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "root",
	Short:   "Coeus is a commandline tool to help you benchmark gRPC methods",
	Example: "coeus --config ./testdata/config/testconfig.json",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return app.LoadConfigFromFile(configFile, &config, &runtimeConfig)
	},
	Run: func(cmd *cobra.Command, args []string) {
		caller := app.Caller{Config: &config}

		err := caller.InitCaller(&runtimeConfig)
		if err != nil {
			panic(err)
		}

		err = caller.Run()
		if err != nil {
			panic(err)
		}

		return
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to the coeus config file")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
