package app

import (
	"net/http"

	"github.com/twilio/twilio-go/twiml"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /answer", app.Answer)
	mux.HandleFunc("POST /route", app.Route)
	return mux
}

func (a *app) Answer(w http.ResponseWriter, r *http.Request) {
	say := &twiml.VoiceSay{
		Message: "Welcome the Rapid Response Hotline, press 1 for Spanish, press 2 for English",
		Loop:    "3",
	}
	gather := twiml.VoiceGather{
		NumDigits: "1",
		Action:    "/route",
		Method:    "POST",
	}
	gather.InnerElements = []twiml.Element{say}
	twimlRes, err := twiml.Voice([]twiml.Element{gather})
	if err != nil {
		http.Error(w, "Failed to say response", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(twimlRes))
	}
}

func (a *app) Route(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	var message string
	option := r.FormValue("Digits")
	if option == "1" {
		message = "Routing to a Spanish speaking agent"
	} else if option == "2" {
		message = "Routing to an English speaking agent"
	} else {
		message = "Routing to next available agent"
	}

	say := &twiml.VoiceSay{
		Message: message,
	}

	twimlRes, err := twiml.Voice([]twiml.Element{say})
	if err != nil {
		http.Error(w, "Failed to say response", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(twimlRes))
	}
}
