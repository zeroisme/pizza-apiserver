package custominitializer

import (
	"k8s.io/apiserver/pkg/admission"

	informers "github.com/zeroisme/pizza-apiserver/pkg/generated/informers/externalversions"
)

type restaurantInformerPluginInitializer struct {
	informers informers.SharedInformerFactory
}

var _ admission.PluginInitializer = restaurantInformerPluginInitializer{}

// New creates an instance of the restaurantInformerPluginInitializer.
func New(informers informers.SharedInformerFactory) restaurantInformerPluginInitializer {
	return restaurantInformerPluginInitializer{
		informers: informers,
	}
}

// Initialize checks the initialization interfaces implemented by a plugin
// and provide the appropriate initialization data
func (i restaurantInformerPluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsRestaurantInformerFactory); ok {
		wants.SetRestaurantInformerFactory(i.informers)
	}
}
