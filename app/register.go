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

type registerPage struct {
	Error   error
	Email   string
	UserID  string
	Current string
}

const registerTemplate = "register.html"

//getRegister renders register page
func (app *appDB) getRegister(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no cookie found
	if !session.CheckCookie(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.render(w, registerTemplate, nil)
	return
}

//postRegister validates and save new user profile
func (app *appDB) postRegister(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no cookie found
	if !session.CheckCookie(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rp := registerPage{}
	email := r.FormValue("email")
	userID := r.FormValue("userID")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	current := r.FormValue("current")

	rp.Email = email
	rp.UserID = userID
	rp.Current = current
	err := validateRegisterInput(email, userID, password1, password2)
	if err != nil {
		rp.Error = err
		logger.Error.Println(err)
		app.render(w, registerTemplate, rp)
		return
	}

	//check if userID already exist
	db := data.InitAppDB(app.client, app.mDB, app.env)
	user, err := db.RetrieveProfile(strings.ToLower(userID))
	if err == nil {
		logger.Error.Println("userID already exist")
		rp.Error = errors.New("userID already taken. please provide a new userID")
		app.render(w, registerTemplate, rp)
		return
	}

	//encrypt password
	bPwd, _ := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.MinCost)

	//get 'use my current location' setting
	if current == "current" {
		user.UseLocation = true
	}

	user.UserID = strings.ToLower(userID)
	user.Email = strings.ToLower(email)
	user.Password = bPwd

	err = db.SaveProfile(user)
	if err != nil {
		logger.Error.Println("failed to save profile")
		app.render(w, registerTemplate, errors.New("Unable to proceed. Please retry."))
		return
	}

	//start user session
	session.Set(w, r, userID)
	app.postMessage(w, "Account Created")
	return
}
