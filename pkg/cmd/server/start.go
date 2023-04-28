package server

import (
	"fmt"
	"io"
	"net"

	"github.com/spf13/cobra"
	"github.com/zeroisme/pizza-apiserver/pkg/admission/custominitializer"
	"github.com/zeroisme/pizza-apiserver/pkg/admission/plugin/pizzatoppings"
	"github.com/zeroisme/pizza-apiserver/pkg/apis/restaurant/v1alpha1"
	"github.com/zeroisme/pizza-apiserver/pkg/apiserver"
	clientset "github.com/zeroisme/pizza-apiserver/pkg/generated/clientset/versioned"
	informers "github.com/zeroisme/pizza-apiserver/pkg/generated/informers/externalversions"
	restaurantopenapi "github.com/zeroisme/pizza-apiserver/pkg/generated/openapi"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
)

const defaultEtcdPathPrefix = "/registry/pizza-apiserver.programming-kubernetes.info"

type CustomServerOptions struct {
	RecommendedOptions    *genericoptions.RecommendedOptions
	SharedInformerFactory informers.SharedInformerFactory
}

func NewCustomServerOptions(out, errOut io.Writer) *CustomServerOptions {
	o := &CustomServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			apiserver.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion),
		),
	}

	return o
}

func (o *CustomServerOptions) Config() (*apiserver.Config, error) {
	err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost",
		nil,
		[]net.IP{net.ParseIP("127.0.0.1")},
	)

	if err != nil {
		return nil, fmt.Errorf("error creating self-signed cert: %v", err)
	}

	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)

	o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
		client, err := clientset.NewForConfig(c.LoopbackClientConfig)
		if err != nil {
			return nil, err
		}
		informerFactory := informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
		o.SharedInformerFactory = informerFactory
		return []admission.PluginInitializer{custominitializer.New(informerFactory)}, nil
	}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(restaurantopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(apiserver.Scheme))
	serverConfig.OpenAPIConfig.Info.Title = "Pizza API"
	serverConfig.OpenAPIConfig.Info.Version = "1.0"

	if utilfeature.DefaultFeatureGate.Enabled(features.OpenAPIV3) {
		serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(restaurantopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(apiserver.Scheme))
		serverConfig.OpenAPIV3Config.Info.Title = "Pizza API"
		serverConfig.OpenAPIV3Config.Info.Version = "1.0"
	}
	err = o.RecommendedOptions.ApplyTo(serverConfig)
	if err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
	}
	return config, nil
}

func (o CustomServerOptions) Run(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	server.GenericAPIServer.AddPostStartHook("start-pizza-apiserver-informers",
		func(context genericapiserver.PostStartHookContext) error {
			config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
			o.SharedInformerFactory.Start(context.StopCh)
			return nil
		},
	)

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}

// Validate validates CustomServerOptions.
func (o CustomServerOptions) Validate() error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *CustomServerOptions) Complete() error {
	// register admission plugins
	pizzatoppings.Register(o.RecommendedOptions.Admission.Plugins)

	// add admission plugins to the RecommendedPluginOrder
	oldOrder := o.RecommendedOptions.Admission.RecommendedPluginOrder
	o.RecommendedOptions.Admission.RecommendedPluginOrder =
		append(oldOrder, "PizzaToppings")
	return nil
}

// NewCommandStartCustomServer provides a CLI handler for 'start master' command
// with a default CustomServerOptions.
func NewCommandStartCustomServer(defaults *CustomServerOptions, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch a custom API server",
		Long:  "Launch a custom API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(stopCh); err != nil {
				return err
			}
			return nil
		},
	}
	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)

	return cmd
}
