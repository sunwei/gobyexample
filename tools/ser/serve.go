package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := "8000"
	publicDir := "public"
	fmt.Printf("Serving 深入理解Hugo at http://127.0.0.1:%s\n", port)
	http.ListenAndServe(":"+port, http.FileServer(http.Dir(publicDir)))
}
