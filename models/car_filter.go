package models

type CarFilters struct {
	Make  string
	Model string
	Year  int // 0 means no year filter
}
