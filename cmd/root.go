/*
Copyright © 2020 Ben Tran <ben.btran@gmail.com>

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
	"fmt"
	"os"
	_ "reflect"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

//var cfgFile string

var (
	version, list             bool
	cluster, service, profile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecs-mon",
	Short: "AWS ECS Monitor",
	Long:  `AWS ECS Monitor`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		// if version {
		// 	fmt.Println(cluster)
		// 	fmt.Println(service)
		// 	return printVersion()

		// }
		// sess, err := session.NewSession(&aws.Config{
		// 	Region: aws.String("us-east-1"),
		// })
		// // Equivalent to session.New
		// sess, err := session.NewSessionWithOptions(session.Options{})

		// // Specify profile to load for the session's config
		// sess, err := session.NewSessionWithOptions(session.Options{
		// 	Profile: "profile_name",
		// })

		// // Specify profile for config and region for requests
		var sess *session.Session
		var err error
		if profile != "" {
			sess, err = session.NewSessionWithOptions(session.Options{
				Config:  aws.Config{Region: aws.String("us-east-1")},
				Profile: profile,
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if service != "" && cluster != "" {
			return getSVCInfo(sess, cluster, service)
		}
		if list {
			return listECSServices(sess, cluster)
		}
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	// rootCmd.Flags().BoolVarP(&version, "version", "v", false, "show current version of CLI")
	rootCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "ECS cluster")
	rootCmd.Flags().StringVarP(&service, "service", "s", "", "ECS service")
	rootCmd.Flags().BoolVarP(&list, "list", "l", false, "List ECS Services")
	rootCmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS profile name")
}
