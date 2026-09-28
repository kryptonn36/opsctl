package cmd

import (
	"github.com/kryptonn36/opsctl/internal/backup"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(dbCmd)
}

var dbCmd = &cobra.Command{
	Use: "db",
	Short: "Database operations",
}


var (
	host string
	port string
	dbname string
	user string
	output_dir string
)

func init() {
	backupCmd.Flags().StringVar(&host, "host", "", "PostgreSQL host")
	backupCmd.Flags().StringVar(&port, "port", "", "Port")
	backupCmd.Flags().StringVar(&dbname, "dbname", "", "Database name")
	backupCmd.Flags().StringVar(&user, "user", "", "User")
	backupCmd.Flags().StringVar(&output_dir, "output_dir", "", "Output directory")
}

var backupCmd = &cobra.Command{
	Use: "backup",
	Short: "backup of db",
	Long: `This command stores the database backup in a compressed zip file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := backup.Config{
			Host: host,
			Port: port,
			User: user,
			DBname: dbname,
			Output_dir: output_dir,
			
		}
		return backup.Run(cmd.Context(), cfg)
	},
}