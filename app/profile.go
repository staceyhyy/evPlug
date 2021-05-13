package app

import (
	"golive/session"
	"net/http"
)

type updatePage struct {
	Error  error
	UserID string
}

const profileTemplate = "updateProfile.html"

//updateProfile renders update profile page
func (app *appDB) updateProfile(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no cookie found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	up := updatePage{}
	userID, _ := session.Get(r)

	up.UserID = userID
	app.render(w, profileTemplate, up)
	return
}
