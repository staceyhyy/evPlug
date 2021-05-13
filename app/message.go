package app

import (
	"net/http"
)

type messagePage struct {
	Message string
}

const messageTemplate = "message.html"

//postMessage renders message page
func (app *appDB) postMessage(w http.ResponseWriter, message string) {
	mp := messagePage{}
	mp.Message = message
	app.render(w, messageTemplate, mp)
	return
}
