
## Data Structures

- **Book**
  - `ID int`
  - `Title string`
  - `Author string`
  - `Status string` - `"Available"` or `"Borrowed"`

- **Member**
  - `ID int`
  - `Name string`
  - `BorrowedBooks []Book`

- **Library (service)**
  - `books map[int]Book`
  - `members map[int]Member`

## Interface

```go
type LibraryManager interface {
    AddBook(book Book)
    RemoveBook(bookID int)
    BorrowBook(bookID int, memberID int) error
    ReturnBook(bookID int, memberID int) error
    ListAvailableBooks() []Book
    ListBorrowedBooks(memberID int) []Book
}
