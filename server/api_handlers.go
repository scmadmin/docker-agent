package server


import (
	"fmt"
	"net/http"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
}

func Serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/containers/stats", statsHandler)
	http.ListenAndServe(":8080", mux)
}
