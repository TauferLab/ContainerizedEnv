package main

import (
	"fmt"

	"github.com/apptainer/apptainer/pkg/cmdline"
	"github.com/spf13/cobra"
)

func callbackRegisterCmd(manager *cmdline.CommandManager) {
	var createFlag *bool
	var runFlag *bool

	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Use the 'workflow'subcommand to create and run workflows",
		Long:  `Use the 'workflow'subcommand to create and run workflows`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowEntryPoint(*createFlag, *runFlag, args)
		},
	}

	createFlag = workflowCmd.Flags().BoolP("create", "c", false, "Pass in a workflow description file as an argument to build a workflow from a JSON description, or use this flag without any arguments to create a workflow from the web interface")
	runFlag = workflowCmd.Flags().BoolP("run", "r", false, "Pass in a workflow description file as an argument to run that workflow")

	manager.RegisterCmd(workflowCmd)
}

func workflowEntryPoint(createFlag, runFlag bool, args []string) error {
	if createFlag && runFlag {
		return fmt.Errorf("Cannot create a workflow and run it at the same time")
	} else if createFlag {
		if len(args) > 0 {
			if err := workflowCreateJSON(args[0]); err != nil {
				return fmt.Errorf("Unable to create workflow from JSON file: %s: %v", args[0], err)
			}
		} else {
			if err := workflowCreateWeb(); err != nil {
				return fmt.Errorf("Unable to create workflow from web interface: %v", err)
			}
		}
	} else if runFlag {
		if len(args) == 0 {
			return fmt.Errorf("Cannot execute workflow without workflow description")
		}
		if err := execWorkflow(args[0]); err != nil {
			return fmt.Errorf("Could not execute workflow: %v", err)
		}
	} else {
		return fmt.Errorf("Must use at least one flag.\nrun `apptainer workflow --help` for more info")
	}

	return nil
}
