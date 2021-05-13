package app

import (
	"encoding/json"
	"regexp"
	"strings"

	"golive/data"
	"golive/logger"
	"golive/session"
	"html/template"
	"net/http"
	"strconv"
)

type indexPage struct {
	Login  bool
	UserID string
	URL    string
}

type Loc struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type filter struct {
	Provider []string `json:"provider"`
}

const (
	url  = "https://maps.googleapis.com/maps/api/js?key="
	parm = "&callback=initMap&libraries=&v=weekly"
)

const indexTemplate = "index.html"

type tpl *template.Template

//getInitMap prepares main page
func (app *appDB) getInitMap(w http.ResponseWriter, r *http.Request) {
	ip := indexPage{}

	//if user has logged-in, get the userID
	userID, found := session.Get(r)
	if !found {
		session.SetCookie(w, r)
	} else {
		ip.Login = true
		ip.UserID = userID
	}
	//set the google map URL
	ip.URL = url + app.env.APIKey + parm
	err := app.tpl.ExecuteTemplate(w, indexTemplate, ip)
	if err != nil {
		logger.Error.Println("rendering template index html in error")
	}
	return
}

//cancel is used by all cancel button
func cancel(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

//getSingapore returns singapore location in latlng :
//- if no user login
//- if user login but set not to use current location
func (app *appDB) getSingapore(w http.ResponseWriter, r *http.Request) {
	//if no cookie found, not allowed to proceed
	logger.Info.Println(r)
	if !session.CheckCookie(r) {
		logger.Error.Println("No cookie found ")
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}

	loc := Loc{}
	userID, found := session.Get(r)
	var useLocation bool
	if found {
		session.Refresh(r)
		//if user has logged-in, get the use current location setting
		db := data.InitAppDB(app.client, app.mDB, app.env)
		user, err := db.RetrieveProfile(userID)
		if err != nil {
			logger.Error.Println("unable to retrieve user profile: ", err)
			http.Error(w, "unable to process", http.StatusUnprocessableEntity)
			return
		}

		useLocation = user.UseLocation
	}

	if !useLocation {
		//set current location as singapore
		lat, _ := strconv.ParseFloat(app.env.Lat, 64)
		lng, _ := strconv.ParseFloat(app.env.Lng, 64)
		loc = Loc{Lat: lat, Lng: lng}
	} else {
		loc = Loc{Lat: 0, Lng: 0}
	}

	err := json.NewEncoder(w).Encode(loc)
	if err != nil {
		logger.Error.Println("Error converting to JSON: ", err)
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}
	return

}

//getLocationPoints returns all the locations points
func (app *appDB) getLocationPoints(w http.ResponseWriter, r *http.Request) {
	//if no cookie found, not allowed to proceed
	if !session.CheckCookie(r) {
		logger.Error.Println("No cookie found ")
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}

	//if user login, retrieve their vehicle charger type
	chargerType := make([]string, 0)
	db := data.InitAppDB(app.client, app.mDB, app.env)
	userID, found := session.Get(r)
	if found {
		session.Refresh(r)
		user, err := db.RetrieveProfile(userID)
		if err != nil {
			logger.Error.Println("unable to retrieve user profile: ", err)
			http.Error(w, "unable to process", http.StatusUnprocessableEntity)
			return
		}
		chargerType = append(user.ChargerType)
	}

	point, err := db.RetrieveAllPoints(chargerType)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, "unable to process", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(point)
	if err != nil {
		logger.Error.Println("Error converting to JSON: ", err)
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}
	return
}

//getLocationInfo returns information for the selected location
func (app *appDB) getLocationInfo(w http.ResponseWriter, r *http.Request) {
	//if no cookie found, not allowed to proceed
	if !session.CheckCookie(r) {
		logger.Error.Println("No cookie found ")
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}

	//get the location latlng
	q := r.URL.Query()
	lat, err := strconv.ParseFloat(q["lat"][0], 64)
	lng, err := strconv.ParseFloat(q["lng"][0], 64)
	if err != nil {
		logger.Fatal.Println("Invalid URL Query parameter")
		http.Error(w, "unable to process", http.StatusBadRequest)
		return
	}

	//if user login, retrieve their vehicle charger type
	location := data.NewPoint(lat, lng)
	db := data.InitAppDB(app.client, app.mDB, app.env)
	chargerType := make([]string, 0)
	userID, found := session.Get(r)
	if found {
		session.Refresh(r)
		user, err := db.RetrieveProfile(userID)
		if err != nil {
			logger.Error.Println("unable to retrieve user profile: ", err)
			http.Error(w, "unable to process", http.StatusUnprocessableEntity)
			return
		}
		chargerType = append(user.ChargerType)
	}

	info, err := db.RetrieveLocationInfo(location, chargerType)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, "unable to process", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(info)
	if err != nil {
		logger.Error.Println("Error converting to JSON: ", err)
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
	}

	return
}

//getFilterList returns provider list
func (app *appDB) getFilterList(w http.ResponseWriter, r *http.Request) {
	//if no cookie found, not allowed to proceed
	if !session.CheckCookie(r) {
		logger.Error.Println("No cookie found ")
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}

	session.Refresh(r)

	filter := filter{}
	filter.Provider = append(app.dataCollections.provider)

	err := json.NewEncoder(w).Encode(filter)
	if err != nil {
		logger.Error.Println("Error converting to JSON: ", err)
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
	}
	// }
	return
}

//getAddressInfo returns all locations that matched with the search address
func (app *appDB) getAddressInfo(w http.ResponseWriter, r *http.Request) {
	//if no cookie found, not allowed to proceed
	if !session.CheckCookie(r) {
		logger.Error.Println("No cookie found ")
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}

	//get the  address
	q := r.URL.Query()
	addr := q["addr"][0]

	//validates address
	if !(regexp.MustCompile("^[a-zA-Z0-9 ]*$").MatchString(addr)) {
		logger.Fatal.Println("Invalid URL Query parameter")
		http.Error(w, "unable to process", http.StatusBadRequest)
		return
	}

	//if user login, retrieve their vehicle charger type
	db := data.InitAppDB(app.client, app.mDB, app.env)
	chargerType := make([]string, 0)
	userID, found := session.Get(r)
	if found {
		user, err := db.RetrieveProfile(userID)
		if err != nil {
			logger.Error.Println("unable to retrieve user profile: ", err)
			http.Error(w, "unable to process", http.StatusUnprocessableEntity)
			return
		}

		chargerType = append(user.ChargerType)
	}

	point, err := db.RetrieveInfoByAddr(strings.Title(addr), chargerType)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, "unable to process", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(point)
	if err != nil {
		logger.Error.Println("Error converting to JSON: ", err)
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
	}
	return
}

//getProviderLocation returns all locations for the provider
func (app *appDB) getProviderLocation(w http.ResponseWriter, r *http.Request) {
	//if no cookie found, not allowed to proceed
	if !session.CheckCookie(r) {
		logger.Error.Println("No cookie found ")
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
		return
	}

	//get the provider
	q := r.URL.Query()
	provider := q["provider"][0]

	if !(regexp.MustCompile("^[a-zA-Z ]*$").MatchString(provider)) {
		logger.Fatal.Println("Invalid URL Query parameter")
		http.Error(w, "unable to process", http.StatusBadRequest)
		return
	}

	//if user login, retrieve their vehicle charger type
	db := data.InitAppDB(app.client, app.mDB, app.env)
	chargerType := make([]string, 0)
	userID, found := session.Get(r)
	if found {
		user, err := db.RetrieveProfile(userID)
		if err != nil {
			logger.Error.Println("unable to retrieve user profile: ", err)
			http.Error(w, "unable to process", http.StatusUnprocessableEntity)
			return
		}

		chargerType = append(user.ChargerType)
	}

	point, err := db.RetrieveLocationByProvider(provider, chargerType)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, "unable to process", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(point)
	if err != nil {
		logger.Error.Println("Error converting to JSON: ", err)
		http.Error(w, "unable to process", http.StatusUnprocessableEntity)
	}
	return
}
