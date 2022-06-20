package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"phoebus/tropa/books/api"

	chi "github.com/go-chi/chi/v5"
)

type ErrorDto struct {
	Message string `json:"message"`
}

type Book struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

var books = []Book{
	Book{Id: 1, Title: "O Apanhador no Campo de Centeio"},
	Book{Id: 2, Title: "The Go Programming Language"},
}

func main() {
	server := api.NewServer("5000")

	server.Router.Get("/hello", hello)
	server.Router.Get("/books", findAll)
	server.Router.Get("/books/{book_id}", findOne)
	server.Router.Post("/books", addOne)
	server.Router.Delete("/books/{book_id}", removeOne)

	server.Start()
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello tropa\n"))
}

func findAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func addOne(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		RespondJSON(w, ErrorDto{Message: "Invalid body"}, http.StatusBadRequest)
		return
	}

	books = append(books, book)
	RespondJSON(w, books, http.StatusCreated)
}

func removeOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "book_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondJSON(w, ErrorDto{Message: "Invalid id"}, http.StatusBadRequest)
		return
	}

	nBooks := make([]Book, 0)
	for _, book := range books {
		if book.Id != id {
			nBooks = append(nBooks, book)
		}
	}

	books = nBooks
	RespondJSON(w, books, http.StatusOK)
}

func findOne(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "book_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondJSON(w, ErrorDto{Message: "Invalid id"}, http.StatusBadRequest)
		return
	}

	for _, book := range books {
		if book.Id == id {
			RespondJSON(w, book, http.StatusOK)
			return
		}
	}

	RespondJSON(w, ErrorDto{Message: "Book not found"}, http.StatusNotFound)
}

func RespondJSON(w http.ResponseWriter, body any, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
