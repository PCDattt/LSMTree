package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

var s *Storage

func main() {
	var err error
	s, err = NewStorage()
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /Add", AddHandler)
	mux.HandleFunc("POST /Find", FindHandler)
	mux.HandleFunc("POST /Flush", FlushHandler)
	mux.HandleFunc("POST /Compact", CompactHandler)

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
		http.Error(w, "Fail to find"+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value+"\n"))
}

func FlushHandler(w http.ResponseWriter, r *http.Request) {
	err := s.Flush()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Fail to flush"+err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func CompactHandler(w http.ResponseWriter, r *http.Request) {
err := s.Compact()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Fail to compact"+err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
