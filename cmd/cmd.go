package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thiagozs/go-proxy-audit/cfg"
)

var (
	rootCmd = &cobra.Command{
		Use:   "paudit",
		Short: "Run proxy server for mysql",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.BindEnv("Mysql")
	viper.BindEnv("Proxy")
	viper.BindEnv("Debug")
	//viper.SetDefault("Mysql", "3306")
	//viper.SetDefault("Proxy", "33060")
	//viper.SetDefault("Debug", false)

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg.Config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

}
