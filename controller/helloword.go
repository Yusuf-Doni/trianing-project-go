package controller

import "net/http"

func HelloWorldController() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!!"))
	}
}