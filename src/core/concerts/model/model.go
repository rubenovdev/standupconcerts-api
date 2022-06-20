package model

type Filters struct {
	Year       uint64
	ComedianId uint64
	SortBy     string // popular || new
	ExcludedId uint64
}
