package cmd

import (
	"github.com/kryptonn36/opsctl/internal/backup"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(backupCmd)
	dbCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(dbCmd)
}

var dbCmd = &cobra.Command{
	Use: "db",
	Short: "Database operations",
}


var (
	backuphost string
	backupport string
	backupdbname string
	backupuser string
	backupoutput_dir string

	restorehost string
	restoreport string
	restoredbname string
	restoreuser string
)

func init() {
	backupCmd.Flags().StringVar(&backuphost, "host", "", "PostgreSQL host")
	backupCmd.Flags().StringVar(&backupport, "port", "", "Port")
	backupCmd.Flags().StringVar(&backupdbname, "dbname", "", "Database name")
	backupCmd.Flags().StringVar(&backupuser, "user", "", "User")
	backupCmd.Flags().StringVar(&backupoutput_dir, "output_dir", "", "Output directory")
}

var backupCmd = &cobra.Command{
	Use: "backup",
	Short: "backup of db",
	Long: `This command stores the database backup in a compressed zip file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := backup.Config{
			Host: backuphost,
			Port: backupport,
			User: backupuser,
			DBname: backupdbname,
			Output_dir: backupoutput_dir,
			
		}
		return backup.Run(cmd.Context(), cfg)
	},
}

func init() {
	backupCmd.Flags().StringVar(&restorehost, "host", "", "PostgreSQL host")
	backupCmd.Flags().StringVar(&restoreport, "port", "", "Port")
	backupCmd.Flags().StringVar(&restoredbname, "dbname", "", "Database name")
	backupCmd.Flags().StringVar(&restoreuser, "user", "", "User")
}

var restoreCmd = &cobra.Command{
	Use: "restore",
	Short: "restore the file that stored in zip folder",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := backup.Config{
			Host: restorehost,
			Port: restoreport,
			User: restoreuser,
			DBname: restoredbname,
		}
		return backup.Restore(cmd.Context(), cfg, args[0])
	},
}