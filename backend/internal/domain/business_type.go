package domain

// BusinessType represents the type of business a company operates
type BusinessType string

const (
	// BusinessTypeRestaurant represents a traditional restaurant
	BusinessTypeRestaurant BusinessType = "restaurant"
	// BusinessTypeBakery represents a bakery
	BusinessTypeBakery BusinessType = "bakery"
	// BusinessTypeConfectionery represents a confectionery/sweet shop
	BusinessTypeConfectionery BusinessType = "confectionery"
	// BusinessTypeCoffeeShop represents a coffee shop/cafe
	BusinessTypeCoffeeShop BusinessType = "coffee_shop"
	// BusinessTypePizzeria represents a pizzeria
	BusinessTypePizzeria BusinessType = "pizzeria"
	// BusinessTypeBurger represents a burger joint
	BusinessTypeBurger BusinessType = "burger"
	// BusinessTypeIceCream represents an ice cream shop
	BusinessTypeIceCream BusinessType = "ice_cream"
	// BusinessTypeAcai represents an acai shop
	BusinessTypeAcai BusinessType = "acai"
	// BusinessTypeFoodTruck represents a food truck
	BusinessTypeFoodTruck BusinessType = "food_truck"
	// BusinessTypeDarkKitchen represents a dark kitchen/ghost kitchen
	BusinessTypeDarkKitchen BusinessType = "dark_kitchen"
	// BusinessTypeGeneric represents a generic food business
	BusinessTypeGeneric BusinessType = "generic"
)

// IsValid checks if a BusinessType is valid
func (bt BusinessType) IsValid() bool {
	switch bt {
	case BusinessTypeRestaurant,
		BusinessTypeBakery,
		BusinessTypeConfectionery,
		BusinessTypeCoffeeShop,
		BusinessTypePizzeria,
		BusinessTypeBurger,
		BusinessTypeIceCream,
		BusinessTypeAcai,
		BusinessTypeFoodTruck,
		BusinessTypeDarkKitchen,
		BusinessTypeGeneric:
		return true
	default:
		return false
	}
}

// String returns the string representation of BusinessType
func (bt BusinessType) String() string {
	return string(bt)
}

// DisplayName returns a human-readable display name for the BusinessType
func (bt BusinessType) DisplayName() string {
	switch bt {
	case BusinessTypeRestaurant:
		return "Restaurante"
	case BusinessTypeBakery:
		return "Padaria"
	case BusinessTypeConfectionery:
		return "Confeitaria"
	case BusinessTypeCoffeeShop:
		return "Cafeteria"
	case BusinessTypePizzeria:
		return "Pizzaria"
	case BusinessTypeBurger:
		return "Hamburgueria"
	case BusinessTypeIceCream:
		return "Sorveteria"
	case BusinessTypeAcai:
		return "Açaí"
	case BusinessTypeFoodTruck:
		return "Food Truck"
	case BusinessTypeDarkKitchen:
		return "Dark Kitchen"
	case BusinessTypeGeneric:
		return "Genérico"
	default:
		return "Desconhecido"
	}
}

// AllBusinessTypes returns a list of all valid BusinessTypes
func AllBusinessTypes() []BusinessType {
	return []BusinessType{
		BusinessTypeRestaurant,
		BusinessTypeBakery,
		BusinessTypeConfectionery,
		BusinessTypeCoffeeShop,
		BusinessTypePizzeria,
		BusinessTypeBurger,
		BusinessTypeIceCream,
		BusinessTypeAcai,
		BusinessTypeFoodTruck,
		BusinessTypeDarkKitchen,
		BusinessTypeGeneric,
	}
}
