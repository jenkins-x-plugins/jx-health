package cmd

import (
	"github.com/jenkins-x-plugins/jx-health/pkg/cmd/status"
	"github.com/jenkins-x-plugins/jx-health/pkg/cmd/version"
	"github.com/jenkins-x-plugins/jx-health/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rootcmd.TopLevelCommand,
		Short: "Health commands",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Logger().Error(err.Error())
			}
		},
	}
	cmd.AddCommand(cobras.SplitCommand(status.NewCmdStatus()))
	cmd.AddCommand(cobras.SplitCommand(version.NewCmdVersion()))

	return cmd
}
