package main

import (
    "fmt"
)

func Hello(text string) string {
	return text
}

func main() {
    log.Println("Starting webhook server ...")

    mux := http.NewServeMux()
    mux.HandleFunc("/mutateTimezone", handleMutateTimezone)

    s := &http.Server{
        Addr:           "8080",
        Handler:        mux,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1048576
    }

}
