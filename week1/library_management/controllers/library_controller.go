package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

type LibraryController struct {
	lib         services.LibraryManager
	libConcrete *services.Library
	in          *bufio.Reader
}

func NewLibraryController(lib *services.Library) *LibraryController {
	return &LibraryController{
		lib:         lib,
		libConcrete: lib,
		in:          bufio.NewReader(os.Stdin),
	}
}

func (c *LibraryController) readLine(prompt string) string {
	fmt.Print(prompt)
	text, _ := c.in.ReadString('\n')
	return strings.TrimSpace(text)
}

func (c *LibraryController) readInt(prompt string) (int, error) {
	raw := c.readLine(prompt)
	return strconv.Atoi(strings.TrimSpace(raw))
}

func (c *LibraryController) seedDemoData() {
	c.libConcrete.AddMember(models.Member{ID: 1, Name: "Alice"})
	c.libConcrete.AddMember(models.Member{ID: 2, Name: "Bob"})

	c.lib.AddBook(models.Book{ID: 100, Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Status: "Available"})
	c.lib.AddBook(models.Book{ID: 101, Title: "Clean Code", Author: "Robert C. Martin", Status: "Available"})
	c.lib.AddBook(models.Book{ID: 102, Title: "Introduction to Algorithms", Author: "Cormen et al.", Status: "Available"})
}

func (c *LibraryController) printMenu() {
	fmt.Println("\n=== Library Management ===")
	fmt.Println("1) Add Book")
	fmt.Println("2) Remove Book")
	fmt.Println("3) Borrow Book")
	fmt.Println("4) Return Book")
	fmt.Println("5) List Available Books")
	fmt.Println("6) List Member's Borrowed Books")
	fmt.Println("7) Add Member (helper)")
	fmt.Println("8) Reserve Book (async borrow)")
	fmt.Println("0) Exit")
}

func (c *LibraryController) addBook() {
	id, err := c.readInt("Enter Book ID: ")
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}
	title := c.readLine("Enter Title: ")
	author := c.readLine("Enter Author: ")
	c.lib.AddBook(models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	})
	fmt.Println("Book added.")
}

func (c *LibraryController) removeBook() {
	id, err := c.readInt("Enter Book ID to remove: ")
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}
	c.lib.RemoveBook(id)
	fmt.Println("Book removed (if it existed).")
}

func (c *LibraryController) borrowBook() {
	bookID, err := c.readInt("Enter Book ID to borrow: ")
	if err != nil {
		fmt.Println("Invalid Book ID.")
		return
	}
	memberID, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Println("Invalid Member ID.")
		return
	}
	if err := c.lib.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Borrow successful.")
}

func (c *LibraryController) returnBook() {
	bookID, err := c.readInt("Enter Book ID to return: ")
	if err != nil {
		fmt.Println("Invalid Book ID.")
		return
	}
	memberID, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Println("Invalid Member ID.")
		return
	}
	if err := c.lib.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Return successful.")
}

func (c *LibraryController) listAvailable() {
	books := c.lib.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	fmt.Println("\nAvailable Books:")
	for _, b := range books {
		fmt.Printf("- [%d] %s — %s (%s)\n", b.ID, b.Title, b.Author, b.Status)
	}
}

func (c *LibraryController) listBorrowed() {
	memberID, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Println("Invalid Member ID.")
		return
	}
	books := c.lib.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Println("No borrowed books for this member (or member not found).")
		return
	}
	fmt.Printf("\nBorrowed Books (Member %d):\n", memberID)
	for _, b := range books {
		fmt.Printf("- [%d] %s — %s (%s)\n", b.ID, b.Title, b.Author, b.Status)
	}
}

func (c *LibraryController) addMember() {
	id, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}
	name := c.readLine("Enter Member Name: ")
	c.libConcrete.AddMember(models.Member{ID: id, Name: name})
	fmt.Println("Member added.")
}

func (c *LibraryController) reserveBook() {
	bookID, err := c.readInt("Enter Book ID to reserve: ")
	if err != nil {
		fmt.Println("Invalid Book ID.")
		return
	}
	memberID, err := c.readInt("Enter Member ID: ")
	if err != nil {
		fmt.Println("Invalid Member ID.")
		return
	}
	if err := c.lib.ReserveBook(bookID, memberID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Reserved. If borrowing doesn't complete in 5s, it will auto-unreserve.")
}

func (c *LibraryController) Run() {
	c.seedDemoData()

	for {
		c.printMenu()
		choice, err := c.readInt("Choose an option: ")
		if err != nil {
			fmt.Println("Please enter a number.")
			continue
		}
		switch choice {
		case 1:
			c.addBook()
		case 2:
			c.removeBook()
		case 3:
			c.borrowBook()
		case 4:
			c.returnBook()
		case 5:
			c.listAvailable()
		case 6:
			c.listBorrowed()
		case 7:
			c.addMember()
		case 8:
			c.reserveBook()
		case 0:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown option.")
		}
	}
}
