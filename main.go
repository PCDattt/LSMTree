package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

var s *Storage

func main() {
	var err error
	s, err = NewStorage(CONST_CURRENT_FILE_NAME, CONST_COMPACTED_FILE_NAME)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /Add", AddHandler)
	mux.HandleFunc("POST /Find", FindHandler)
	mux.HandleFunc("POST /Flush", FlushHandler)

	fmt.Println("LSM Tree Server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}

type NodeRequest struct {
	Key int `json:"key"`
	Value string `json:"value"`
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req NodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid json "+err.Error(), http.StatusBadRequest)
		return
	}
	s.Add(req.Key, req.Value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func FindHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req NodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid json "+err.Error(), http.StatusBadRequest)
		return
	}
	value, err := s.Find(req.Key)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Fail to find: "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value+"\n"))
}

func FlushHandler(w http.ResponseWriter, r *http.Request) {
	// Compact file before flush to keep file's keys in sorted order
	err := s.Compact()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Fail to compact: "+err.Error(), http.StatusInternalServerError)
	}
	err = s.Flush()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Fail to flush: "+err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

