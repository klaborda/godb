package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/klaborda/godb/models"
)

type mockBookModel struct{}

func (m *mockBookModel) All() ([]models.Book, error) {
	var bks []models.Book

	bks = append(bks, models.Book{1, "Stranger in a Strange Land", "Robert A Heinlein"})
	bks = append(bks, models.Book{2, "Friday", "Robert A Heinlein"})

	return bks, nil
}

func TestBooksIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)

	env := Env{books: &mockBookModel{}}

	http.HandlerFunc(env.booksIndex).ServeHTTP(rec, req)

	expected := "1, Stranger in a Strange Land, Robert A Heinlein\n2, Friday, Robert A Heinlein\n"
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}
