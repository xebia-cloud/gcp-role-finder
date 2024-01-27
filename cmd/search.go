/*
Copyright © 2024 Xebia Netherlands B.V
*/
package cmd

import (
	"encoding/json"
	"os"
	"strings"
	"xebia-cloud/gcp-role-finder/internal/search/fulltext"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Args:  cobra.MinimumNArgs(1),
	Short: "for roles with the specified IAM permission",
	Long: `
searches for the roles with the desired IAM permission, the role with the least permissions first.

The command line arguments provide a fulltext search query. 

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		roles, err := roleRepository.GetRoles(cmd.Context())
		if err != nil {
			return err
		}
		excludedRoles, _ := cmd.Flags().GetStringSlice("excluded-roles")
		searcher, err := fulltext.NewSearcher(cmd.Context(), fulltext.WithExcludedRoles(excludedRoles))
		if err != nil {
			return err
		}
		if err = searcher.IndexRoles(cmd.Context(), roles); err != nil {
			return err
		}
		query := strings.Join(args, " ")
		matchingRoles, err := searcher.FindRoles(cmd.Context(), query)
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(matchingRoles)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().StringSlice("excluded-roles", []string{"roles/owner", "roles/editor", "roles/viewer"}, "from search result")
}
