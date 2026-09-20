package cmd

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use: "doctor",
	Short: "It will check that all things are fine or not in local development",
	Long: `It verifies that the host machine have all the necessary toolchain and running services before starting the work or if something now 
	going according to our will then it told where is the fault`,

	RunE: func(cmd *cobra.Command, args []string) error {
		binaries := []string{"git", "docker", "kubectl", "psql"}
		for _, bin := range binaries{
			err := CheckBinaries(bin)
			if err == exec.ErrNotFound{
				fmt.Println(bin, "not found")
			} else if err != nil{
				return fmt.Errorf("error in running command try again: %v", err)
			}else{
				fmt.Println(bin, "found")
			}
		}
		if CheckDockerDeamon(){
			fmt.Println("docker deamon is running")
		}else{
			fmt.Println("docker deamon is not running")
		}
		return nil
	},
}

func CheckBinaries(bin string) error{

	_, err := exec.LookPath(bin)
	if err != nil{
		return err
	}
	return nil
}

func CheckDockerDeamon() bool{
	cli, err  := client.New(client.FromEnv)
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Ping(context.Background(), client.PingOptions{})
	return err == nil
}