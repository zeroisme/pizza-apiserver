package pizzatoppings

import (
	"context"
	"fmt"
	"io"

	"github.com/zeroisme/pizza-apiserver/pkg/admission/custominitializer"
	"github.com/zeroisme/pizza-apiserver/pkg/apis/restaurant"
	informers "github.com/zeroisme/pizza-apiserver/pkg/generated/informers/externalversions"
	listers "github.com/zeroisme/pizza-apiserver/pkg/generated/listers/restaurant/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register("PizzaToppings", func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

type PizzaToppingsPlugin struct {
	*admission.Handler
	toppingLister listers.ToppingLister
}

var _ = custominitializer.WantsRestaurantInformerFactory(&PizzaToppingsPlugin{})
var _ = admission.ValidationInterface(&PizzaToppingsPlugin{})

// Admit ensures that the object in-flight is of kind Pizza.
// In addition checks that the toppings are known.
func (d *PizzaToppingsPlugin) Validate(ctx context.Context, a admission.Attributes, _ admission.ObjectInterfaces) error {
	// we are only interested in pizzas
	if a.GetKind().GroupKind() != restaurant.Kind("Pizza") {
		return nil
	}

	if !d.WaitForReady() {
		return admission.NewForbidden(a, fmt.Errorf("not yet ready to handle request"))
	}

	obj := a.GetObject()
	pizza := obj.(*restaurant.Pizza)
	for _, top := range pizza.Spec.Toppings {
		if _, err := d.toppingLister.Get(top.Name); err != nil && errors.IsNotFound(err) {
			return admission.NewForbidden(
				a,
				fmt.Errorf("unknown topping: %s", top.Name),
			)
		}
	}
	return nil
}

// SetRestaurantInformerFactory gets Lister from SharedInformerFactory.
// The lister knows how to list toppings.
func (d *PizzaToppingsPlugin) SetRestaurantInformerFactory(f informers.SharedInformerFactory) {
	d.toppingLister = f.Restaurant().V1alpha1().Toppings().Lister()
	d.SetReadyFunc(f.Restaurant().V1alpha1().Toppings().Informer().HasSynced)
}

// ValidaValidateInitalization checks whether the plugin was correctly initialized.
func (d *PizzaToppingsPlugin) ValidateInitialization() error {
	if d.toppingLister == nil {
		return fmt.Errorf("missing topping lister")
	}
	return nil
}

// New creates a new instance of the PizzaToppingsPlugin.
func New() (*PizzaToppingsPlugin, error) {
	return &PizzaToppingsPlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}, nil
}
