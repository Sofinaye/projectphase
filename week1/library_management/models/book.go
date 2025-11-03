package models

// Book represents a library book.
type Book struct {
	ID     int
	Title  string
	Author string
	// Status can be "Available" or "Borrowed"
	Status string
}
