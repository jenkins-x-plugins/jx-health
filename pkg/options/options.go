package options

import (
	"fmt"

	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	khcheckcrd "github.com/kuberhealthy/kuberhealthy/v2/pkg/apis/khcheck/v1"
	khstatecrd "github.com/kuberhealthy/kuberhealthy/v2/pkg/apis/khstate/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// constants for using the kuberhealthy check CRD
const checkCRDGroup = "comcast.github.io"
const checkCRDVersion = "v1"

// KHCheckOptions common CLI arguments for working with health
type KHCheckOptions struct {
	StateClient *khstatecrd.KHStateV1Client
	CheckClient *khcheckcrd.KHCheckV1Client
}

// Validate validates the options and returns the KuberhealthyCheckClient
func (o *KHCheckOptions) Validate() error {

	f := kubeclient.NewFactory()
	cfg, err := f.CreateKubeConfig()
	if err != nil {
		log.Logger().Fatalf("failed to get kubernetes config: %v", err)
	}

	if o.StateClient == nil {
		o.StateClient, err = ClientStateClient(cfg, checkCRDGroup, checkCRDVersion)
		if err != nil {
			return fmt.Errorf("failed to create kuberhealthy check client: %w", err)
		}
	}

	if o.CheckClient == nil {
		o.CheckClient, err = ClientCheckClient(cfg, checkCRDGroup, checkCRDVersion)
		if err != nil {
			return fmt.Errorf("failed to create kuberhealthy check client: %w", err)
		}
	}

	return nil
}

// ClientStateClient creates a rest client to use for interacting with Kuberhealthy state CRDs
func ClientStateClient(cfg *rest.Config, groupName string, groupVersion string) (*khstatecrd.KHStateV1Client, error) {
	var err error

	err = khcheckcrd.ConfigureScheme(groupName, groupVersion)
	if err != nil {
		return &khstatecrd.KHStateV1Client{}, err
	}

	config := *cfg
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: groupName, Version: groupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	return khstatecrd.NewForConfig(&config)
}

// ClientStateClient creates a rest client to use for interacting with Kuberhealthy check CRDs
func ClientCheckClient(cfg *rest.Config, groupName string, groupVersion string) (*khcheckcrd.KHCheckV1Client, error) {
	var err error

	err = khcheckcrd.ConfigureScheme(groupName, groupVersion)
	if err != nil {
		return &khcheckcrd.KHCheckV1Client{}, err
	}

	config := *cfg
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: groupName, Version: groupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	return khcheckcrd.NewForConfig(&config)
}
