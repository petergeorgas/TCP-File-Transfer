package cmd

import (
	"file_trans/client/transfer"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var cfgFile string
var transferCmd = &cobra.Command{
	Use:   "tcp-transfer [file(s)]",
	Short: "tcp-transfer is a simple file transfer application",
	Long:  `tcp-transfer is a simple CLI application that allows you to transfer files from a client to a specified server using TCP.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Do stuff here
		//fmt.Println(args[1])

		var wg sync.WaitGroup
		wg.Add(len(args))

		for _, fileName := range args {
			go transfer.TransferFile("localhost:8080", fileName, &wg)
		}
		wg.Wait()
	},
}

func Execute() {
	if err := transferCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {

	transferCmd.AddCommand()
}
