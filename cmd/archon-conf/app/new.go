package app

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"kubeup.com/archon/pkg/jsonnet"
)

var (
	profile string
	configs []string
)

var cmdNew = &cobra.Command{
	Use:   "new [resouce type] [resource name]",
	Short: "Generate resource definition with given type and name",
	Long:  "Generate resource definition with given type and name",
	RunE:  newRun,
}

func newRun(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Not enough args")
	}
	vm, err := jsonnet.Make(profile)
	if err != nil {
		return err
	}
	for i := 0; i < len(configs); i++ {
		kv := strings.SplitN(configs[i], "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("Error in --config argument: Expected arg to be 'key=value': %q", configs[i])
		}
		vm.Config(kv[0], kv[1])
	}
	json, err := vm.New(args[0], args[1])
	if err != nil {
		return err
	}
	fmt.Print(json)
	vm.Destroy()
	return nil
}

func init() {
	RootCmd.AddCommand(cmdNew)
	cmdNew.Flags().StringVarP(&profile, "profile", "p", "", "Archon profile to use")
	cmdNew.Flags().StringArrayVarP(&configs, "config", "c", []string{}, "Set config variable for the resource, repeat this flag for multiple items")
}
