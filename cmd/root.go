package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "opsctl",
	Short : "opsctl is a devops tool to moniter",
	Long : `It is used to scan local environment for required tools, verifies daemon
	socket and checks network connectivity also help in interactive cleanup`,

	Run : func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command starting")
	},
}

func Excute(){
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1) 
	} 
}