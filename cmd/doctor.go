package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

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
			fmt.Println("Docker deamon is not running")
		}
		if CheckPostgres(){
			fmt.Println("Postgres is running")
		}else{
			fmt.Println("Postgres is not connected")
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
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	if pingOk(ctx, client.FromEnv){
		return true
	}

	home, err := os.UserHomeDir()
	if err == nil{
		socketHost := "unix:///"+ home + "/.docker/desktop/docker.sock"
		if pingOk(ctx, client.WithHost(socketHost)){
			return true
		}
	}
	return false
}

func pingOk(ctx context.Context, opt client.Opt) bool {
	cli, err := client.New(opt)
	if err != nil{
		return false
	}

	_, err = cli.Ping(ctx, client.PingOptions{})
	return err == nil
}

func CheckPostgres() bool{
	target := "localhost:5433"

	conn, err := net.DialTimeout("tcp", target, 2*time.Second)
	if err != nil {
		return false
	}else{
		_ = conn.Close()
		return true
	}

}