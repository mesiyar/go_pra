package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"github.com/justinas/alice"
)

// app encapsulates Env ,Router and middleware
type App struct {
	Router *mux.Router
	MiddleWares *MiddleWare
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkResp struct {
	ShortLink string `json:"shortlink"`
}

// Initialize is initialization of app
func (a *App) Initialize() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.MiddleWares = &MiddleWare{}
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	m := alice.New(a.MiddleWares.LoggingHandler, a.MiddleWares.RecoverHandler)
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortLink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.createShortLinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithError(w, StatusError{http.StatusBadRequest, fmt.Errorf("parse parameters failed %v", r.Body)})
		return
	}

	if err := validator.Validate(req); err != nil {
		responseWithError(w, StatusError{http.StatusBadRequest, fmt.Errorf("validate parameters failed %v", req)})
		return
	}

	defer r.Body.Close()
	fmt.Printf("%v\n", req)

}

func (a *App) createShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("ShortLink")
	fmt.Printf("%v\n", s)
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shrotlink"])
}

// run the app
func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, a.Router))
}
