package health

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sort"
	"strings"
	"syscall"

	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	khcheckcrd "github.com/kuberhealthy/kuberhealthy/v2/pkg/apis/khcheck/v1"

	"github.com/liggitt/tabwriter"

	"github.com/jenkins-x/jx-helpers/v3/pkg/knative_pkg/duck"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/jenkins-x-plugins/jx-health/pkg/health/lookup"

	khstatecrd "github.com/kuberhealthy/kuberhealthy/v2/pkg/apis/khstate/v1"

	"github.com/jenkins-x-plugins/jx-health/pkg/options"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const resourceStates = "khstates"
const resourceChecks = "khchecks"

// Options common CLI arguments for working with health
type Options struct {
	options.KHCheckOptions
	Info     bool
	InfoData lookup.LoopkupData
}

// WriteStatusTable returns a table containing statuses
func (o Options) WriteStatusTable(table *tabwriter.Writer, ns string) error {

	err := o.KHCheckOptions.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate KHCheckOptions: %w", err)
	}

	// get a list of all Kuberhealthy states
	states, err := o.KHCheckOptions.StateClient.KuberhealthyStates(ns).List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list resource %s in namespace %s: %w", resourceStates, ns, err)
	}

	// get a list of all Kuberhealthy checks
	checks, err := o.KHCheckOptions.CheckClient.KuberhealthyChecks(ns).List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list resource %s in namespace %s: %w", resourceChecks, ns, err)
	}

	o.annotateStates(&states, &checks)

	rows := o.populateTable(&states)
	for _, row := range rows {
		_, err = fmt.Fprintln(table, strings.Join(row, "\t"))
		if err != nil {
			log.Logger().Infof("error formatting row: %v", err)
		}
	}
	table.Flush()
	return nil
}

func (o Options) populateTable(checks *khstatecrd.KuberhealthyStateList) [][]string {

	sort.Slice(checks.Items, func(i, j int) bool {
		return checks.Items[i].Name < checks.Items[j].Name
	})

	var rows [][]string
	// add Kuberhealthy check results to the table
	for _, check := range checks.Items {
		rows = append(rows, o.populateRow(check)...)
	}
	return rows
}

func (o Options) populateRow(check khstatecrd.KuberhealthyState) [][]string {
	var rows [][]string

	status := "ERROR"
	if check.Spec.OK {
		status = termcolor.ColorInfo("OK")
	}

	// get matching information link
	informationDetail := ""
	if check.Annotations != nil {
		informationDetail = check.Annotations["docs.jenkins-x.io"]
	}
	if informationDetail == "" {
		informationDetail = o.InfoData.Info[check.Name]
	}

	// depending on if there are errors or how many there are we want to format the table to it is easy to consume
	// Name | Namespace | Status | Error Message        | Info (optional)
	// foo    jx          ok
	// bar    jx          error    first error for bar
	//                             second error for bar
	// cheese jx          ok
	rowEntries := []string{check.Name, check.Namespace, status}
	if len(check.Spec.Errors) == 0 {
		rowEntries = append(rowEntries, "")
		if o.Info {
			rowEntries = append(rowEntries, informationDetail)
		}
		rows = append(rows, rowEntries)
	} else {
		rowEntries = append(rowEntries, check.Spec.Errors[0])
		if o.Info {
			rowEntries = append(rowEntries, informationDetail)
		}
		rows = append(rows, rowEntries)

		// if we have multiple errors lets format the table so all errors appear underneath in the column
		if len(check.Spec.Errors) > 1 {
			for i := 1; i < len(check.Spec.Errors); i++ {
				rowEntries = []string{"", "", "", check.Spec.Errors[i]}
				rows = append(rows, rowEntries)
			}
		}
	}
	return rows
}

// WatchStates will watch for new and updates to Kuberhealthy checks and write rows when changes happen
func (o Options) WatchStates(table *tabwriter.Writer, cfg *rest.Config, namespace string) error {

	// Grab a dynamic interface that we can create informers from
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("could not generate dynamic client for config: %w", err)
	}
	// Create a factory object that we can say "hey, I need to watch this resource"
	// and it will give us back an informer for it
	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, namespace, nil)
	// Retrieve a "GroupVersionResource" type that we need when generating our informer from our dynamic factory
	gvr, _ := schema.ParseResourceArg("khstates.v1.comcast.github.io")
	// Finally, create our informer for deployments!
	i := f.ForResource(*gvr)

	stopCh := make(chan struct{})
	go o.startWatching(stopCh, i.Informer(), table)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	close(stopCh)

	return nil
}

func (o Options) startWatching(stopCh <-chan struct{}, s cache.SharedIndexInformer, table *tabwriter.Writer) {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newUnstructured := obj.(*unstructured.Unstructured)
			newState := &khstatecrd.KuberhealthyState{}

			err := duck.FromUnstructured(newUnstructured, newState)
			if err != nil {
				log.Logger().Infof("error converting unstructured object %s into KuberhealthyState: %v", newUnstructured.GetName(), err)
				return
			}
			o.writeRow(newState, table)
		},
		UpdateFunc: func(oldObj, obj interface{}) {

			newUnstructured := obj.(*unstructured.Unstructured)
			newState := &khstatecrd.KuberhealthyState{}

			err := duck.FromUnstructured(newUnstructured, newState)
			if err != nil {
				log.Logger().Infof("error converting unstructured object %s into KuberhealthyState: %v", newUnstructured.GetName(), err)
				return
			}

			oldUnstructured := oldObj.(*unstructured.Unstructured)
			oldState := &khstatecrd.KuberhealthyState{}

			err = duck.FromUnstructured(oldUnstructured, oldState)
			if err != nil {
				log.Logger().Infof("error converting unstructured object %s into KuberhealthyState: %v", oldUnstructured.GetName(), err)
				return
			}

			if newState.Spec.OK != oldState.Spec.OK || !reflect.DeepEqual(newState.Spec.Errors, oldState.Spec.Errors) {
				o.writeRow(newState, table)
			}
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
}

func (o Options) writeRow(state *khstatecrd.KuberhealthyState, table *tabwriter.Writer) {

	rows := o.populateRow(*state)
	for _, row := range rows {
		_, err := fmt.Fprintln(table, strings.Join(row, "\t"))
		if err != nil {
			log.Logger().Infof("error formatting row: %v", err)
		}
	}
	err := table.Flush()
	if err != nil {
		log.Logger().Infof("error printing row: %v", err)
	}
}

func (o Options) annotateStates(stateList *khstatecrd.KuberhealthyStateList, checkList *khcheckcrd.KuberhealthyCheckList) {
	m := map[string]*khcheckcrd.KuberhealthyCheck{}
	for i := range checkList.Items {
		check := &checkList.Items[i]
		m[check.Name] = check
	}

	for i := range stateList.Items {
		state := &stateList.Items[i]
		check := m[state.Name]
		if check != nil && check.Annotations != nil {
			if state.Annotations == nil {
				state.Annotations = map[string]string{}
			}
			for k, v := range check.Annotations {
				if state.Annotations[k] == "" {
					state.Annotations[k] = v
				}
			}
		}
	}
}
