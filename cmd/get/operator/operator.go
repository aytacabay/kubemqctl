package operator

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/cmd/get/operator/logs"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

type GetOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get Kubemq operators list 
	kubemqctl get operators  
	# Get Kubemq operator pod logs 
	kubemqctl get operators log
`
var getLong = `Get Kubemq Operators List`
var getShort = `Get Kubemq Operators List`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &GetOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "operator",
		Aliases: []string{"operators", "op", "o"},
		Short:   getShort,
		Long:    getLong,
		Example: getExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.AddCommand(logs.NewCmdLogs(ctx, cfg))
	return cmd
}

func (o *GetOptions) Complete(args []string) error {
	return nil
}

func (o *GetOptions) Validate() error {
	return nil
}

func (o *GetOptions) Run(ctx context.Context) error {
	newClient, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	operatorManager, err := operator.NewManager(newClient)
	if err != nil {
		return err
	}

	utils.Println("Getting Kubemq Operators List...")
	operators, err := operatorManager.GetKubemqOperatorsDeployments()
	if err != nil {
		return err
	}
	if len(operators) == 0 {
		return fmt.Errorf("no Kubemq operators were found in the cluster")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tNAMSPACE\n")
	for _, item := range operators {
		fmt.Fprintf(w, "%s\t%s\n",
			item.Name,
			item.Namespace,
		)
	}
	_ = w.Flush()
	return nil
}
