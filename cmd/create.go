package cmd

import (
	"fmt"
	"pm/internal/archive"
	"pm/internal/config"
	"pm/internal/ssh"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [config.json]",
	Short: "Create and upload package to server",
	Long:  `Creates an archive from files matching patterns defined in the configuration file and uploads it to the remote server via SSH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var packet config.Packet

		if len(args) != 1 {
			return fmt.Errorf("invalid count of flags")
		}

		if err := config.LoadPacketConfig(args[0], &packet); err != nil {

			return err
		}

		archiveName := fmt.Sprintf("%s_%s.tar.gz", packet.Name, packet.Ver)

		if err := archive.CreateArchive(archiveName, packet.Targets); err != nil {
			return err
		}

		client, err := ssh.New("localhost", "root", "./vault/test-key/test_key")
		if err != nil {
			return err
		}
		defer client.Close()

		remotePath := ssh.RemotePath + archiveName
		if err := client.Upload(archiveName, remotePath); err != nil {
			return err
		}

		fmt.Printf("✅ Загружен на сервер: %s\n", remotePath)
		return nil

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
