package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := ":8080"
	http.HandleFunc("/", homeHandler)
	fmt.Println("server started on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("error starting server: %s", err)
	}
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome")
}
