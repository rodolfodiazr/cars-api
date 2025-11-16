package data

import "cars/models"

func Cars() map[string]models.Car {
	return map[string]models.Car{
		"JHK290XJ": {
			ID:   "JHK290XJ",
			Make: "Ford", Model: "F10", Package: "Base",
			Color: "Silver", Year: 2010, Category: "Truck",
			Mileage: 120123, Price: 1999900,
		},
		"FWL37LA": {
			ID:   "FWL37LA",
			Make: "Toyota", Model: "Camry", Package: "SE",
			Color: "White", Year: 2019, Category: "Sedan",
			Mileage: 3999, Price: 2899000,
		},
		"1I3XJRLLC": {
			ID:   "1I3XJRLLC",
			Make: "Toyota", Model: "Rav4", Package: "XSE",
			Color: "Red", Year: 2018, Category: "SUV",
			Mileage: 24001, Price: 2275000,
		},
		"DKU43920S": {
			ID:   "DKU43920S",
			Make: "Ford", Model: "Bronco", Package: "Badlands",
			Color: "Burnt Orange", Year: 2022, Category: "SUV",
			Mileage: 1, Price: 4499000,
		},
	}
}
