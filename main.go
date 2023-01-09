package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	errStorageNotFound = errors.New("storage: key not found")
	errValidateEmpty   = errors.New("validate: key empty")
)

func main() {
	// store
	storage := newStore()

	// server
	srv := &http.Server{
		Addr: ":9000",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			req := r.URL.Query().Get("id")

			value, err := logic(storage, req)
			if errors.Is(err, errStorageNotFound) {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			if errors.Is(err, errValidateEmpty) {
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			res, _ := json.Marshal(value)
			w.Write(res)

		default:
			w.WriteHeader(http.StatusForbidden)
		}
	})

	log.Print("listening...")
	srv.ListenAndServe()
}

// storage
type storage struct {
	store map[string]interface{}
}

func newStore() *storage {
	return &storage{
		store: map[string]interface{}{
			"id_0": "john smith",
			"id_1": "joe blogs",
		},
	}
}

func (s storage) getFromStore(key string) (interface{}, error) {
	value, ok := s.store[key]
	if !ok {
		return nil, errStorageNotFound
	}

	return value, nil
}

// validate
func reqEmpty(key string) error {
	if key == "" {
		return errValidateEmpty
	}

	return nil
}

// logic
func logic(storage *storage, key string) (interface{}, error) {
	if err := reqEmpty(key); err != nil {
		return nil, err
	}

	value, err := storage.getFromStore(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}
