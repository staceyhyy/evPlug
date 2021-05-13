package app

import (
	"errors"
	"golive/data"
	"golive/logger"
	"golive/session"
	"net/http"
)

type chgEmailPage struct {
	Error error
	Email string
}

const chgEmailTemplate = "email.html"

//getChangeEmail renders change email page
func (app *appDB) getChangeEmail(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cp := chgEmailPage{}
	userID, _ := session.Get(r)

	//retrieve current user vehicle data
	db := data.InitAppDB(app.client, app.mDB, app.env)
	user, err := db.RetrieveProfile(userID)
	if err != nil {
		cp.Error = errors.New("Unable to process. Please re-login.")
		app.render(w, chgEmailTemplate, cp)
		return
	}
	cp.Email = user.Email
	app.render(w, chgEmailTemplate, cp)
	return
}

//postChangeEmail saves user email update
func (app *appDB) postChangeEmail(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cp := chgEmailPage{}
	email := r.FormValue("email")

	if !validEmail(email) {
		cp.Error = errors.New("please enter a valid email address")
		logger.Error.Println("invalid email")
		app.render(w, chgEmailTemplate, cp)
		return
	}

	userID, _ := session.Get(r)

	db := data.InitAppDB(app.client, app.mDB, app.env)
	err := db.UpdateUserEmail(userID, email)
	if err != nil {
		app.render(w, chgEmailTemplate, errors.New("Unable to process. Please re-login."))
		return
	}
	logger.Info.Println("Email updated: " + userID)
	app.postMessage(w, "Email address Updated")
	return
}
