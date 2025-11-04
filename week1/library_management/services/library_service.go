package services

import (
	"errors"
	"fmt"
	"library_management/models"
	"slices"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book

	ReserveBook(bookID int, memberID int) error
}

type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

func (l *Library) AddMember(m models.Member) {
	l.members[m.ID] = m
}

func (l *Library) GetMember(id int) (models.Member, bool) {
	m, ok := l.members[id]
	return m, ok
}

func (l *Library) GetBook(id int) (models.Book, bool) {
	b, ok := l.books[id]
	return b, ok
}

func (l *Library) AddBook(book models.Book) {
	// If not specified, default to Available
	if book.Status == "" {
		book.Status = "Available"
	}
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	if book.Status != "Available" {
		return errors.New("book is not available")
	}

	book.Status = "Borrowed"
	l.books[bookID] = book

	// Append to member's borrowed list.
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	// Verify that the member actually borrowed this book.
	idx := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("member did not borrow this book")
	}

	// Remove from member borrowed slice (keep order stable).
	member.BorrowedBooks = slices.Delete(member.BorrowedBooks, idx, idx+1)
	l.members[memberID] = member

	// Mark book available again.
	book.Status = "Available"
	l.books[bookID] = book
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	var result []models.Book
	for _, b := range l.books {
		if b.Status == "Available" {
			result = append(result, b)
		}
	}
	// Sort by ID for consistent output (Go 1.21 has slices.SortFunc)
	slices.SortFunc(result, func(a, b models.Book) int {
		return a.ID - b.ID
	})
	return result
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, ok := l.members[memberID]
	if !ok {
		return nil
	}
	// Return a copy for safety.
	out := make([]models.Book, len(member.BorrowedBooks))
	copy(out, member.BorrowedBooks)
	slices.SortFunc(out, func(a, b models.Book) int {
		return a.ID - b.ID
	})
	return out
}
