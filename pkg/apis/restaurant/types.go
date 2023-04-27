package restaurant

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Pizza specifies an offered pizza with toppings.
type Pizza struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   PizzaSpec
	Status PizzaStatus
}

type PizzaSpec struct {
	// Toppings is a list of toppings to be added to the pizza.
	Toppings []PizzaTopping
}

type PizzaTopping struct {
	Name     string
	Quantity int
}

type PizzaStatus struct {
	// cost is the cost of the whole pizza including all toppings.
	Cost float64
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PizzaList is a list of Pizza objects.
type PizzaList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Pizza
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Topping is a topping put onto a pizza
type Topping struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec ToppingSpec
}

type ToppingSpec struct {
	Cost float64
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ToppingList is a list of Topping objects.
type ToppingList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Topping
}
