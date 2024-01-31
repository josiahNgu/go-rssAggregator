package main

import "net/http"

// second params is the pointer to a http request
// the func name is standard to go library package expect
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
