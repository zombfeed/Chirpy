package main

import "net/http"

func main() {
	smux := http.NewServeMux()
	server := http.Server{}
	server.Handler = smux
	server.Addr = ":8080"

	server.ListenAndServe()
}
