package custominitializer

import (
	"k8s.io/apiserver/pkg/admission"

	informers "github.com/zeroisme/pizza-apiserver/pkg/generated/informers/externalversions"
)

// WantsRestaurantInformerFactory defines an interface for objects that want a RestaurantInformerFactory.
type WantsRestaurantInformerFactory interface {
	SetRestaurantInformerFactory(informers.SharedInformerFactory)
	admission.InitializationValidator
}
