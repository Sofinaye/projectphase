package services

import (
	"errors"
	"fmt"
	"library_management/models"
	"slices"
	"sync"
	"time"
)

// LibraryManager interface specifies the library operations.
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book

	// NEW: Reservation entry point
	ReserveBook(bookID int, memberID int) error
}

// internal reservation tracking
type reservation struct {
	memberID int
	cancel   chan struct{} // closes when borrow succeeds (cancels timer)
}

type reservationRequest struct {
	bookID   int
	memberID int
}

// Library implements LibraryManager.
type Library struct {
	mu           sync.Mutex
	books        map[int]models.Book
	members      map[int]models.Member
	reservations map[int]*reservation    // bookID -> reservation
	reserveCh    chan reservationRequest // queue of reservation processing
	wg           sync.WaitGroup
}

// NewLibrary creates a new Library service instance.
func NewLibrary() *Library {
	l := &Library{
		books:        make(map[int]models.Book),
		members:      make(map[int]models.Member),
		reservations: make(map[int]*reservation),
		reserveCh:    make(chan reservationRequest, 32), // buffered queue
	}
	go l.processReservations() // background worker
	return l
}

// ---- Optional helpers (not in the interface) ----

// AddMember registers a member (simple helper for seeding/CLI).
func (l *Library) AddMember(m models.Member) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.members[m.ID] = m
}

// GetMember returns a member by ID.  ✅ keep the param & return types
func (l *Library) GetMember(id int) (models.Member, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	m, ok := l.members[id]
	return m, ok
}

// GetBook returns a book by ID.  ✅ keep the param & return types
func (l *Library) GetBook(id int) (models.Book, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.books[id]
	return b, ok
}

// ---- Interface implementation ----

func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if book.Status == "" {
		book.Status = "Available"
	}
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// If reserved, cancel timer and drop reservation
	if res, ok := l.reservations[bookID]; ok {
		closeSafe(res.cancel)
		delete(l.reservations, bookID)
	}
	delete(l.books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	switch book.Status {
	case "Available":
		// OK to borrow directly
	case "Reserved":
		// Only the reserver may borrow
		res, ok := l.reservations[bookID]
		if !ok || res.memberID != memberID {
			return errors.New("book is reserved by another member")
		}
	default:
		return errors.New("book is not available")
	}

	// transition to Borrowed
	book.Status = "Borrowed"
	l.books[bookID] = book

	// attach to member
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	// if it was reserved, cancel timer + clear reservation
	if res, ok := l.reservations[bookID]; ok {
		closeSafe(res.cancel)
		delete(l.reservations, bookID)
	}
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	// verify member actually borrowed it
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

	// remove from member's borrowed list
	member.BorrowedBooks = slices.Delete(member.BorrowedBooks, idx, idx+1)
	l.members[memberID] = member

	// mark available; also drop any stale reservation (paranoia)
	book.Status = "Available"
	l.books[bookID] = book
	if res, ok := l.reservations[bookID]; ok {
		closeSafe(res.cancel)
		delete(l.reservations, bookID)
	}
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	var result []models.Book
	for _, b := range l.books {
		if b.Status == "Available" {
			result = append(result, b)
		}
	}
	slices.SortFunc(result, func(a, b models.Book) int { return a.ID - b.ID })
	return result
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, ok := l.members[memberID]
	if !ok {
		return nil
	}
	out := make([]models.Book, len(member.BorrowedBooks))
	copy(out, member.BorrowedBooks)
	slices.SortFunc(out, func(a, b models.Book) int { return a.ID - b.ID })
	return out
}

// ---- NEW: Reservation flow ----

// ReserveBook:
//   - If book is Available: mark as Reserved, start a 5s timer goroutine,
//     and enqueue an async borrow request via reserveCh.
//   - If already Reserved: error.
//   - If Borrowed/other: error.
func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	if _, ok := l.members[memberID]; !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	if book.Status == "Reserved" {
		return errors.New("book is already reserved")
	}
	if book.Status != "Available" {
		return errors.New("book is not available to reserve")
	}

	// mark reserved
	book.Status = "Reserved"
	l.books[bookID] = book

	// create reservation and timer cancel channel
	res := &reservation{
		memberID: memberID,
		cancel:   make(chan struct{}),
	}
	l.reservations[bookID] = res

	// start auto-unreserve timer goroutine
	go l.autoUnreserve(bookID, res.cancel)

	// enqueue async borrowing (processed by worker)
	l.reserveCh <- reservationRequest{bookID: bookID, memberID: memberID}

	return nil
}

// background worker: processes queued reservation borrow requests
func (l *Library) processReservations() {
	for req := range l.reserveCh {
		// process each request in its own goroutine to allow parallelism
		go l.tryAsyncBorrow(req)
	}
}

// attempt to complete borrow for an existing reservation
func (l *Library) tryAsyncBorrow(req reservationRequest) {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[req.bookID]
	if !ok {
		return
	}

	// must still be reserved by same member
	res, ok := l.reservations[req.bookID]
	if !ok || res.memberID != req.memberID || book.Status != "Reserved" {
		return
	}

	// perform the borrow
	member := l.members[req.memberID]
	book.Status = "Borrowed"
	l.books[req.bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[req.memberID] = member

	// cancel timer & clear reservation
	closeSafe(res.cancel)
	delete(l.reservations, req.bookID)
}

// autoUnreserve runs a 5s timer; if not canceled by successful borrow, it reverts to Available.
func (l *Library) autoUnreserve(bookID int, cancel <-chan struct{}) {
	select {
	case <-time.After(5 * time.Second):
		l.mu.Lock()
		defer l.mu.Unlock()
		// still reserved?
		if _, ok := l.reservations[bookID]; ok {
			// make the book available again if it's still Reserved
			if b, ok2 := l.books[bookID]; ok2 && b.Status == "Reserved" {
				b.Status = "Available"
				l.books[bookID] = b
			}
			delete(l.reservations, bookID)
		}
	case <-cancel:
		// borrow succeeded; just exit
		return
	}
}

// small helper to safely close a channel once
func closeSafe(ch chan struct{}) {
	defer func() { _ = recover() }()
	close(ch)
}
