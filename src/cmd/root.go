/*
 Copyright 2023 adamswanglin

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cmd

import (
	"flag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"mail2dingrobot/internal"
	"os"
)

var (
	config string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mail2dingrobot",
	Short: "邮箱中转服务",
	Long:  `虚拟邮件server，将接收到的邮件转发到钉钉机器人，另外可配置同时转发到邮箱`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		internal.LoadConfiguration(config)
		internal.InitServer()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mail2dingrobot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	rootCmd.Flags().StringVarP(&config, "config", "c", "", "配置文件地址")
}
