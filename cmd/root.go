// Copyright © 2016 Zhang Peihao

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/BPing/aliyun-live-go-sdk/client"
	"github.com/BPing/aliyun-live-go-sdk/device/cdn"
)

var cfgFile string
var cfgAccessKeyId, cfgAccessKeySecret string
var cfgUrl string
var cfgDomainName string
var cfgDebug bool

type RefreshTResponse struct {
	RefreshTaskId string
	RequestId string
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "alicdn",
	Short: "阿里云CDN刷新程序",
	Long: `直接通过命令刷新CDN.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(cfgAccessKeySecret) == 0 || len(cfgAccessKeyId) == 0 || len(cfgDomainName) == 0 || len(cfgUrl) == 0 {
			return fmt.Errorf("Options miss")
		}
		cert := client.NewCredentials(cfgAccessKeyId, cfgAccessKeySecret)
		cdnM := cdn.NewCDN(cert).SetDebug(cfgDebug)
		// resp := make(map[string]interface{})
		var resp RefreshTResponse
		err := cdnM.RefreshObjectCaches(cfgUrl, cdn.DirectoryRefreshType, &resp)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		fmt.Println("RefreshTaskId: ", resp.RefreshTaskId, ", RequestId: ", resp.RequestId)
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alicdn.yaml)")
	RootCmd.PersistentFlags().StringVar(&cfgAccessKeyId, "access-key-id", "", "access key id")
	RootCmd.PersistentFlags().StringVar(&cfgAccessKeySecret, "access-key-secret", "", "access key secret")
	RootCmd.PersistentFlags().StringVar(&cfgDomainName, "domain-name", "", "CDN domain name")
	RootCmd.PersistentFlags().StringVar(&cfgUrl, "url", "", "The url to refresh")
	RootCmd.PersistentFlags().BoolVarP(&cfgDebug, "debug", "d", false, "debug mode.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".alicdn") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
