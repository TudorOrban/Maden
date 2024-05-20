package main

import (
	"fmt"
	"net/http"
	"sync"	
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/pods", createPodHandler).Methods("POST")
	r.HandleFunc("/pods", listPodsHandler).Methods("GET")
	r.HandleFunc("/pods/{id}", deletePodHandler).Methods("DELETE")
	r.HandleFunc("/nodes", registerNodeHandler).Methods("POST")
	r.HandleFunc("/nodes", listNodesHandler).Methods("GET")
	http.Handle("/", r)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Maden API Server")
}


var nodeDB = []Node{}
var podDB = []Pod{}
var mu sync.Mutex

