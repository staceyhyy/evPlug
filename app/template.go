package app

import (
	"golive/logger"
	"net/http"
)

//render is called by all the handlers to execute template based on the parameter passed
func (app *appDB) render(w http.ResponseWriter, tmpl string, param interface{}) {
	err := app.tpl.ExecuteTemplate(w, tmpl, param)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		logger.Error.Println("failed to execute template: ", tmpl, "err: ", err)
		return
	}
	return
}
