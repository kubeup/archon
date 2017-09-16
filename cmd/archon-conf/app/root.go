package app

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "archon-conf",
	Short: "Archon resource definition generator",
	Long:  `archon-conf embeds jsonnet tool and definitions into one easy to use program. It will generate Archon resource definition based user provided arguments.`,
}
