package app

import (
	"errors"
	"golive/data"
	"golive/logger"
	"golive/session"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type chgPwdPage struct {
	Error error
}

const chgPwdTemplate = "chgPwd.html"

//getChangePwd renders change password page
func (app *appDB) getChangePwd(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.render(w, chgPwdTemplate, nil)
	return
}

//postChangePwd saves user password update
func (app *appDB) postChangePwd(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cp := chgPwdPage{}
	oldPwd := r.FormValue("oldPassword")
	newPwd1 := r.FormValue("newPassword1")
	newPwd2 := r.FormValue("newPassword2")

	err := validateChgPwd(oldPwd, newPwd1, newPwd2)
	if err != nil {
		cp.Error = err
		logger.Error.Println(err)
		app.render(w, chgPwdTemplate, cp)
		return
	}

	userID, _ := session.Get(r)

	//encrypt new password
	bPwd, _ := bcrypt.GenerateFromPassword([]byte(newPwd1), bcrypt.MinCost)

	db := data.InitAppDB(app.client, app.mDB, app.env)
	err = db.UpdateUserPwd(userID, bPwd)
	if err != nil {
		app.render(w, chgPwdTemplate, errors.New("Unable to process. Please re-login."))
		return
	}

	logger.Info.Println("password updated: " + userID)
	app.postMessage(w, "Password updated")
	return
}

func validateChgPwd(oldPwd, newPwd1, newPwd2 string) error {
	if oldPwd == "" || newPwd1 == "" || newPwd2 == "" {
		return errors.New("please enter all the fields")
	}

	if !validPwd(oldPwd) || !validPwd(newPwd1) || !validPwd(newPwd2) || (newPwd1 != newPwd2) {
		return errors.New("please enter a valid password")
	}

	return nil
}
