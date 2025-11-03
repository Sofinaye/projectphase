package services

import (
	"errors"
	"fmt"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
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
