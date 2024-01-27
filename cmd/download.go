/*
Copyright © 2024 Xebia Netherlands B.V
*/
package cmd

import (
	"xebia-cloud/gcp-role-finder/internal/storage/fs"

	"github.com/spf13/cobra"
)

// downloadCmd represents the refresh command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "downloads all GCP IAM roles",
	Long: `the latest IAM role definitions into a file.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		roles, err := roleRepository.GetRoles(cmd.Context())
		if err != nil {
			return err
		}
		filename, _ := cmd.Flags().GetString("data-file")
		repository, err := fs.NewRepository(cmd.Context(), fs.WithFile(filename))
		if err != nil {
			return err
		}
		return repository.SaveRoles(cmd.Context(), roles)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
