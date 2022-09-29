package main

import (
	"fmt"
	"net/http"
)

func main() {
	myhttp := http.NewServeMux()
	fs := http.FileServer(http.Dir("./."))
	myhttp.Handle("/", http.StripPrefix("", fs))

	myhttp.HandleFunc("/chat", chatSocketCreate)
	// myhttp.HandleFunc("/content", socketReaderCreate)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", myhttp)
}
