package status

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/jenkins-x/jx-helpers/v3/pkg/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"

	"k8s.io/client-go/kubernetes"

	healthopts "github.com/jenkins-x-plugins/jx-health/pkg/health"
	"github.com/jenkins-x-plugins/jx-health/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	cmdLong = templates.LongDesc(`
		Prints health statuses in a table
`)

	cmdExample = templates.Examples(`
		# prints all health statuses for the current namespace in a table
		%s get status

		# prints specific health check statuses by passing arguments
		%s get status default jenkins-x
	`)
)

// LabelOptions the options for the command
type Options struct {
	HealthOptions healthopts.Options
	Args          []string
	All           bool
	AllNamespaces bool
	Namespace     string

	KubeClient kubernetes.Interface
}

// NewCmdStatus creates a command object for the command
func NewCmdStatus() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"statuses"},
		Short:   "Gets health statuses",
		Long:    cmdLong,
		Example: fmt.Sprintf(cmdExample, rootcmd.BinaryName, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			o.Args = args
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "if present, list the requested object(s) across all namespaces.\nNamespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "namespace to get status checks, defaults to current namespace")

	return cmd, o
}

// Validate verifies settings
func (o *Options) Validate() error {
	err := o.HealthOptions.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate health options")
	}

	f := kubeclient.NewFactory()
	cfg, err := f.CreateKubeConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to get kubernetes config")
	}

	o.KubeClient, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return errors.Wrapf(err, "error building kubernetes clientset")
	}

	return nil
}

// Run transforms the YAML files
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate status options")
	}

	namespaces, err := o.getNamespaces()
	if err != nil {
		return errors.Wrapf(err, "failed to get which namespaces to look for health statuses")
	}

	resultTable := table.CreateTable(os.Stdout)

	// add table headers
	resultTable.AddRow("Name", "Namespace", "Status", "Error Message")

	for _, n := range namespaces {
		err := o.HealthOptions.GetJenkinsXTable(&resultTable, n)
		if err != nil {
			return errors.Wrapf(err, "failed to build health table, is Kuberhealthy installed?")
		}
	}

	if len(resultTable.Rows) < 2 { // first row is the column headers
		return fmt.Errorf("failed to find any health status rows for namespace %s", o.Namespace)
	}

	resultTable.Render()
	return nil
}

// decide the namespace to search for kuberhealthy states in this order
// 1. if --all-namespaces set then first lookup the namespaces user has permission to list
// 2. if --namespace is set then use the value provided
// 3. the default is to search for health statuses in the current namespace
func (o *Options) getNamespaces() ([]string, error) {
	var namespaces []string
	switch {
	case o.AllNamespaces:
		// get all namespaces
		nList, err := o.KubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "error listing namespaces")
		}
		for _, n := range nList.Items {
			namespaces = append(namespaces, n.Name)
		}
		sort.Strings(namespaces)

	case o.Namespace != "":
		namespaces = []string{o.Namespace}

	default:
		ns, err := kubeclient.CurrentNamespace()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find the current namespace")
		}
		namespaces = []string{ns}

	}
	return namespaces, nil
}