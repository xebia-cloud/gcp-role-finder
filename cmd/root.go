package cmd

import (
	"context"
	"os"

	"xebia-cloud/gcp-role-finder/internal"
	"xebia-cloud/gcp-role-finder/internal/storage/fs"
	"xebia-cloud/gcp-role-finder/internal/storage/gcp"

	"github.com/binxio/gcloudconfig"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
)

var roleRepository internal.RoleRepository

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gcp-role-finder",
	Short: "Explore the appropriate Google Cloud Platform IAM roles",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		parents, _ := cmd.PersistentFlags().GetStringSlice("parents")
		if err = internal.ValidateRoleParents(parents); err != nil {
			return err
		}
		roleRepository, err = MakeRepository(cmd.Context(), cmd)
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// MakeRepository makes a role repository depending on the command flags.
func MakeRepository(ctx context.Context, cmd *cobra.Command) (internal.RoleRepository, error) {
	fromFile, _ := cmd.Flags().GetBool("from-file")
	filename, _ := cmd.Flags().GetString("data-file")
	useDefaultCredentials, _ := cmd.Flags().GetBool("use-default-credentials")
	parents, _ := cmd.Flags().GetStringSlice("parent")

	if cmd.Flags().Changed("parent") && len(parents) == 0 {
		parents = []string{""}
	}

	var err error
	if fromFile {
		return fs.NewRepository(ctx, fs.WithFile(filename))
	}

	var credentials *google.Credentials
	if useDefaultCredentials || !gcloudconfig.IsGCloudOnPath() {
		credentials, err = google.FindDefaultCredentials(ctx)
	} else {
		credentials, err = gcloudconfig.GetCredentials("")
	}
	if err != nil {
		return nil, err
	}

	return gcp.NewRepository(ctx, gcp.WithCredentials(credentials), gcp.WithParents(parents))
}

func init() {
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().Bool("from-file", false, "read the roles from file")
	rootCmd.PersistentFlags().String("data-file", "data/roles.json", "to read the roles from")
	rootCmd.PersistentFlags().Bool("use-default-credentials", false, "use the Google default credentials")
	rootCmd.PersistentFlags().StringSlice("parent", []string{}, "to read the roles from")
}
