package v1beta1

func init() {
	localSchemeBuilder.Register(RegisterDefaults)
}

func SetDefaults_PizzaSpec(obj *PizzaSpec) {
	if len(obj.Toppings) == 0 {
		obj.Toppings = []PizzaTopping{
			{"salami", 1},
			{"mozzarella", 1},
			{"tomato", 1},
		}
	}

	for i := range obj.Toppings {
		if obj.Toppings[i].Quantity == 0 {
			obj.Toppings[i].Quantity = 1
		}
	}
}
