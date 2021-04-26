// +build amd64,linux amd64,darwin arm64,darwin

package machine

import (
	"github.com/containers/common/pkg/completion"
	"github.com/containers/podman/v3/cmd/podman/registry"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/machine"
	"github.com/containers/podman/v3/pkg/machine/qemu"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:               "create [options] [NAME]",
		Short:             "Create a vm",
		Long:              "Create a virtual machine for Podman to run on. Virtual machines are used to run Podman.",
		RunE:              create,
		Args:              cobra.MaximumNArgs(1),
		Example:           `podman machine create myvm`,
		ValidArgsFunction: completion.AutocompleteNone,
	}
)

type CreateCLIOptions struct {
	CPUS         uint64
	Memory       uint64
	Devices      []string
	ImagePath    string
	IgnitionPath string
	Name         string
}

var (
	createOpts                = CreateCLIOptions{}
	defaultMachineName string = "podman-machine-default"
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: createCmd,
		Parent:  machineCmd,
	})
	flags := createCmd.Flags()

	cpusFlagName := "cpus"
	flags.Uint64Var(
		&createOpts.CPUS,
		cpusFlagName, 1,
		"Number of CPUs. The default is 0.000 which means no limit",
	)
	_ = createCmd.RegisterFlagCompletionFunc(cpusFlagName, completion.AutocompleteNone)

	memoryFlagName := "memory"
	flags.Uint64VarP(
		&createOpts.Memory,
		memoryFlagName, "m", 2048,
		"Memory (in MB)",
	)
	_ = createCmd.RegisterFlagCompletionFunc(memoryFlagName, completion.AutocompleteNone)

	ImagePathFlagName := "image-path"
	flags.StringVar(&createOpts.ImagePath, ImagePathFlagName, "", "Path to qcow image")
	_ = createCmd.RegisterFlagCompletionFunc(ImagePathFlagName, completion.AutocompleteDefault)

	IgnitionPathFlagName := "ignition-path"
	flags.StringVar(&createOpts.IgnitionPath, IgnitionPathFlagName, "", "Path to ignition file")
	_ = createCmd.RegisterFlagCompletionFunc(IgnitionPathFlagName, completion.AutocompleteDefault)
}

// TODO should we allow for a users to append to the qemu cmdline?
func create(cmd *cobra.Command, args []string) error {
	createOpts.Name = defaultMachineName
	if len(args) > 0 {
		createOpts.Name = args[0]
	}
	vmOpts := machine.CreateOptions{
		CPUS:         createOpts.CPUS,
		Memory:       createOpts.Memory,
		IgnitionPath: createOpts.IgnitionPath,
		ImagePath:    createOpts.ImagePath,
		Name:         createOpts.Name,
	}
	var (
		vm     machine.VM
		vmType string
		err    error
	)
	switch vmType {
	default: // qemu is the default
		vm, err = qemu.NewMachine(vmOpts)
	}
	if err != nil {
		return err
	}
	return vm.Create(vmOpts)
}
