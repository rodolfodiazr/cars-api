package models

// CarFilters defines optional criteria used to filter cars when listing them.
type CarFilters struct {
	Make  string
	Model string
	Year  int // 0 means no year filter
}
