package session

import (
	"golive/envvar"
	"golive/logger"
	"net/http"
	"time"
)

type sessdetail struct {
	UserID    string
	StartTime time.Time
}

const timeoutMins = 10

//Sess stores session information
var Sess = map[string]sessdetail{}

//Set registers the sessionID
func Set(w http.ResponseWriter, r *http.Request, UserID string) {
	//get sessionID
	cookie, err := r.Cookie(envvar.CookieName())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		logger.Error.Println(err)
		return
	}

	//set session
	Sess[cookie.Value] = sessdetail{UserID, time.Now()}
	logger.Info.Println("start session for " + UserID)
	return
}

//LoggedIn returns true if user has active session
func LoggedIn(w http.ResponseWriter, r *http.Request) bool {
	var currentSession sessdetail

	cookie, err := r.Cookie(envvar.CookieName())
	if err != nil {
		logger.Error.Println(err)
		return false
	}

	//validate session
	currentSession = Sess[cookie.Value] //get session detail
	if currentSession.UserID == "" {
		logger.Warning.Println("user not logged in : " + cookie.Value)
		return false
	}
	return true
}

//Expire returns true if session time already expire.  TimeoutMins set as const
func Expire(w http.ResponseWriter, r *http.Request) bool {
	var currentSession sessdetail
	//get sessionID
	cookie, _ := r.Cookie(envvar.CookieName())
	currentSession = Sess[cookie.Value] //get session detail

	//check expire/timeout
	if time.Since(currentSession.StartTime) > (timeoutMins * time.Minute) {
		logger.Warning.Println("session expire : " + currentSession.UserID)

		//clear session & cookie
		Delete(r)
		return true
	}
	return false
}

//Get returns userID & session found status
func Get(r *http.Request) (string, bool) {
	var currentSession sessdetail
	cookie, err := r.Cookie(envvar.CookieName())
	if err != nil {
		return "", false
	}
	currentSession = Sess[cookie.Value] //get current session detail
	if currentSession.UserID == "" {
		return "", false
	}

	return currentSession.UserID, true
}

//Refresh resets the session time
func Refresh(r *http.Request) {
	var currentSession sessdetail
	cookie, _ := r.Cookie(envvar.CookieName())

	//reset session time
	currentSession = Sess[cookie.Value] //get current session detail
	Sess[cookie.Value] = sessdetail{currentSession.UserID, time.Now()}
	logger.Info.Println("session refreshed: " + currentSession.UserID)
	return
}

//Delete removes sessionID from the map
func Delete(r *http.Request) {
	cookie, _ := r.Cookie(envvar.CookieName())

	logger.Info.Println("remove session for " + cookie.Value)
	delete(Sess, cookie.Value)
	return
}

//Check validates if the user has done login and the session is still valid.
//It will returns false for unauthorised user
func Check(w http.ResponseWriter, r *http.Request) bool {
	//check for unauthorised access
	if !LoggedIn(w, r) {
		return false
	}

	if Expire(w, r) {
		return false
	}
	Refresh(r)
	return true
}
