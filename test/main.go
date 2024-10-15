package main

import (
    "io/ioutil"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/api.php", func(w http.ResponseWriter, r *http.Request) {
        // Read the request body
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Error reading body", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        // Log the request details
        log.Printf("Received request: %s\n", body)

        // Respond with a simple message
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Request received"))
    })

    log.Println("Listening on :8081...")
    log.Fatal(http.ListenAndServe(":8081", nil))
}

