package pizza

import (
	"context"
	"fmt"

	"github.com/zeroisme/pizza-apiserver/pkg/apis/restaurant"
	"github.com/zeroisme/pizza-apiserver/pkg/apis/restaurant/validation"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

// NewStrategy creates and returns a pizzaStrategy instance
func NewStrategy(typer runtime.ObjectTyper) pizzaStrategy {
	return pizzaStrategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, the presence of Initializers if any
// and error in case the given runtime.Object is not a Pizza
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*restaurant.Pizza)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Pizza")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchPizza is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields
func MatchPizza(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *restaurant.Pizza) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type pizzaStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (pizzaStrategy) NamespaceScoped() bool {
	return true
}

func (pizzaStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (pizzaStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (pizzaStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	pizza := obj.(*restaurant.Pizza)
	return validation.ValidatePizza(pizza)
}

func (pizzaStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (pizzaStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (pizzaStrategy) Canonicalize(obj runtime.Object) {
}

func (pizzaStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (pizzaStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return []string{}
}

func (pizzaStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return []string{}
}
