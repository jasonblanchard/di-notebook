/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"fmt"
	"log"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// ReadEntryCmd represents the ReadEntry command
var ReadEntryCmd = &cobra.Command{
	Use:   "ReadEntry",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		id, _ := cmd.Flags().GetString("id")
		principalID, _ := cmd.Flags().GetString("principal")

		address := fmt.Sprintf("%s:%s", host, port)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := notebook.NewNotebookClient(conn)

		request := &notebook.ReadEntryGRPCRequest{
			Principal: &notebook.Principal{
				Id:   principalID,
				Type: notebook.Principal_USER,
			},
			Payload: &notebook.ReadEntryGRPCRequest_Payload{
				Id: id,
			},
		}

		ctx := context.TODO()

		response, err := client.ReadEntry(ctx, request)

		fmt.Println(response)

		return err
	},
}

func init() {
	grpcCmd.AddCommand(ReadEntryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ReadEntryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ReadEntryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ReadEntryCmd.Flags().StringP("host", "o", "localhost", "host")
	ReadEntryCmd.Flags().StringP("port", "p", "8080", "port")
	ReadEntryCmd.Flags().StringP("id", "i", "", "id")
	ReadEntryCmd.Flags().StringP("principal", "w", "", "principal id")
}
