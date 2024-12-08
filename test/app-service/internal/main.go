package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "App Service is running")
	})

	fmt.Println("App Service started on port 8081")
	http.ListenAndServe(":8081", nil)
}
