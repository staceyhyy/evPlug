package app

import (
	"errors"
	"golive/data"
	"golive/logger"
	"golive/session"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type loginPage struct {
	Error error
}

const loginTemplate = "login.html"

//getLogin renders login page
func (app *appDB) getLogin(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no cookie found
	if !session.CheckCookie(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.render(w, loginTemplate, nil)
	return
}

//postLogin validates login
func (app *appDB) postLogin(w http.ResponseWriter, r *http.Request) {
	session.SetCookie(w, r)

	lp := loginPage{}
	userID := r.FormValue("userID")
	password := r.FormValue("password")

	//validates input
	if !validUserID(userID) || !validPwd(password) {
		lp.Error = errors.New("invalid userID/Password")
		app.render(w, loginTemplate, lp)
		return
	}

	//check if userID exist
	db := data.InitAppDB(app.client, app.mDB, app.env)
	user, err := db.RetrieveProfile(strings.ToLower(userID))
	if err != nil {
		lp.Error = errors.New("invalid userID/Password")
		app.render(w, loginTemplate, lp)
		return
	}

	//compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error.Println("invalid password")
		lp.Error = errors.New("invalid userID/password")
		app.render(w, loginTemplate, lp)
		return
	}

	//start user session
	session.Set(w, r, userID)
	logger.Info.Println("login success: " + userID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (app *appDB) logout(w http.ResponseWriter, r *http.Request) {
	session.Delete(r)
	logger.Info.Println("user logout")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
