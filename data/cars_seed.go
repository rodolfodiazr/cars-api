package data

import (
	"cars/models"
	u "cars/pkg/utils"
)

func Cars() map[string]models.Car {
	return map[string]models.Car{
		"JHK290XJ": {
			ID:   "JHK290XJ",
			Make: "Ford", Model: "F10", Package: u.Ptr("Base"),
			Color: "Silver", Year: 2010, Category: "Truck",
			Mileage: u.Ptr(int64(120123)), Price: u.Ptr(int64(1999900)),
		},
		"FWL37LA": {
			ID:   "FWL37LA",
			Make: "Toyota", Model: "Camry", Package: u.Ptr("SE"),
			Color: "White", Year: 2019, Category: "Sedan",
			Mileage: u.Ptr(int64(3999)), Price: u.Ptr(int64(2899000)),
		},
		"1I3XJRLLC": {
			ID:   "1I3XJRLLC",
			Make: "Toyota", Model: "Rav4", Package: u.Ptr("XSE"),
			Color: "Red", Year: 2018, Category: "SUV",
			Mileage: u.Ptr(int64(24001)), Price: u.Ptr(int64(2275000)),
		},
		"DKU43920S": {
			ID:   "DKU43920S",
			Make: "Ford", Model: "Bronco", Package: u.Ptr("Badlands"),
			Color: "Burnt Orange", Year: 2022, Category: "SUV",
			Mileage: u.Ptr(int64(1)), Price: u.Ptr(int64(4499000)),
		},
	}
}
