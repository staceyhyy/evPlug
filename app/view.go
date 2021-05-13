package app

import (
	"errors"
	"golive/data"
	"golive/session"
	"net/http"
)

type viewPage struct {
	Error   error
	Current string
}

const viewTemplate = "view.html"

//getView renders 'Use My Current Location' setting page
func (app *appDB) getView(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	vp := viewPage{}
	userID, _ := session.Get(r)

	//retrieve 'use my current location' setting
	db := data.InitAppDB(app.client, app.mDB, app.env)
	user, err := db.RetrieveProfile(userID)
	if err != nil {
		vp.Error = errors.New("Unable to process. Please re-login.")
		app.render(w, viewTemplate, vp)
		return
	}
	if user.UseLocation {
		vp.Current = "Checked"
	}

	app.render(w, viewTemplate, vp)
	return
}

//postView saves 'Use My Current Location' setting update
func (app *appDB) postView(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vp := viewPage{}
	current := r.FormValue("current")
	r.ParseForm()
	var useLocation bool

	if current == "current" {
		useLocation = true
	}

	userID, _ := session.Get(r)

	db := data.InitAppDB(app.client, app.mDB, app.env)
	err := db.UpdateUserView(userID, useLocation)
	if err != nil {
		vp.Error = errors.New("Unable to process. Please re-login.")
		app.render(w, viewTemplate, vp)
		return
	}
	app.postMessage(w, "View profile updated")
	return
}
