package v1alpha1

import (
	"github.com/zeroisme/pizza-apiserver/pkg/apis/restaurant"
	"k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha1_PizzaSpec_To_restaurant_PizzaSpec(in *PizzaSpec, out *restaurant.PizzaSpec, s conversion.Scope) error {
	idx := map[string]int{}
	for _, top := range in.Toppings {
		if i, duplicate := idx[top]; duplicate {
			out.Toppings[i].Quantity++
			continue
		}
		idx[top] = len(out.Toppings)
		out.Toppings = append(out.Toppings, restaurant.PizzaTopping{
			Name:     top,
			Quantity: 1,
		})
	}
	return nil
}

func Convert_restaurant_PizzaSpec_To_v1alpha1_PizzaSpec(in *restaurant.PizzaSpec, out *PizzaSpec, s conversion.Scope) error {
	for i := range in.Toppings {
		for j := 0; j < in.Toppings[i].Quantity; j++ {
			out.Toppings = append(out.Toppings, in.Toppings[i].Name)
		}
	}
	return nil
}
