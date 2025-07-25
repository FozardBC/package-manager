/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pm/internal/archive"
	"pm/internal/config"
	"pm/internal/ssh"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [packages.json]",
	Short: "Download and install packages",
	Long: `Downloads specified packages from the remote server via SSH and extracts them locally. (./installed)
The list of packages and version constraints is read from a JSON or YAML configuration file.
Supports version specifiers like >=1.0, <=2.0, or exact versions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var packages config.PackagesFile

		if len(args) != 1 {
			return fmt.Errorf("invalid count of flags")
		}

		cfg := config.MustLoad(".env")

		if err := config.LoadPacketConfig(args[0], &packages); err != nil {
			return err
		}

		client, err := ssh.New(cfg.SshHost, cfg.SshUser, cfg.SshKeyPath)
		if err != nil {
			return err
		}
		defer client.Close()

		cache := ".pm_cache"
		if err := os.MkdirAll(cache, 0755); err != nil {
			return err
		}

		for _, p := range packages.Packages {

			filename, err := client.Search(p.Name, p.Ver)
			if err != nil {
				if errors.Is(err, ssh.ErrNotFound) {
					fmt.Println(err.Error())
					continue
				}
				return err
			}

			remotePath := filepath.Join(ssh.RemotePath, filename)
			localPath := filepath.Join(cache, filename)

			if err := client.Download(remotePath, localPath); err != nil {
				return err
			}

			destDir := filepath.Join("installed", p.Name)
			fmt.Printf("📂 Exctracting into %s\n", destDir)
			if err := archive.ExtractArchive(localPath, destDir); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
