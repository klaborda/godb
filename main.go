package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/klaborda/godb/models"

	_ "github.com/mattn/go-sqlite3"
)

type Env struct {
	books interface {
		All() ([]models.Book, error)
	}
}

func main() {
	dbname := "./books.db"
	var db *sql.DB

	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := os.Stat(dbname); errors.Is(err, os.ErrNotExist) {
		log.Printf("%s does not exist", dbname)

		setupDb(dbname, db)
	} 

	// Initalise Env with a models.BookModel instance (which in turn wraps
	// the connection pool).
	env := &Env{
		books: models.BookModel{DB: db},
	}

	log.Print("Listening on :3000")

	http.HandleFunc("/books", env.booksIndex)
	http.ListenAndServe(":3000", nil)
}

func setupDb(dbname string, db *sql.DB) {
	log.Printf("Populating database, %s", dbname)

	sqlStmt := `
	create table books (id integer not null primary key, title text, author text);
	delete from books;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into books(id, title, author) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("Title%03d", i), fmt.Sprintf("Author%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (env *Env) booksIndex(w http.ResponseWriter, r *http.Request) {
	// Execute the SQL query by calling the All() method.
	bks, err := env.books.All()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%d, %s, %s\n", bk.Id, bk.Title, bk.Author)
	}
}
