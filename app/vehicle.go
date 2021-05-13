package app

import (
	"errors"
	"golive/data"
	"golive/logger"
	"golive/session"
	"net/http"
	"strconv"
)

type stringMap map[string]string

type vehiclePage struct {
	Error   error
	Vehicle string
	Model   stringMap
	Charger stringMap
}

const vehicleTemplate = "vehicle.html"

//getVehicle renders user vehicle update page
func (app *appDB) getVehicle(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vp := vehiclePage{}
	userID, _ := session.Get(r)

	//retrieve current user vehicle data
	vp = app.loadVehicleInputData(userID)
	app.render(w, vehicleTemplate, vp)
	return
}

//postVehicle saves user vehicle update
func (app *appDB) postVehicle(w http.ResponseWriter, r *http.Request) {
	//redirect to main page if there's no session found
	if !session.Check(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vp := vehiclePage{}
	userID, _ := session.Get(r)

	//retrieve current user vehicle data
	vp = app.loadVehicleInputData(userID)

	//get user input
	r.ParseForm()
	charger := make([]string, 0)
	for k, v := range r.Form {
		if k != "model" {
			charger = append(v)
		}
	}

	modelIdx, _ := strconv.Atoi(r.FormValue("model"))

	//validate input
	if vp.Vehicle == "" && modelIdx == 0 {
		vp.Error = errors.New("please select one of the option")
		app.render(w, vehicleTemplate, vp)
		return
	}

	if (modelIdx > len(app.dataCollections.vehicle)) || (len(charger) == 0) {
		vp.Error = errors.New("please select one of the option")
		app.render(w, vehicleTemplate, vp)
		return
	}
	model := ""
	if vp.Vehicle != "" && modelIdx == 0 {
		model = vp.Vehicle
	} else {
		model = app.dataCollections.vehicle[modelIdx-1]
	}

	//updates changes
	db := data.InitAppDB(app.client, app.mDB, app.env)
	err := db.UpdateUserVehicle(userID, model, charger)

	if err != nil {
		vp.Error = errors.New("Unable to process. Please re-login.")
		app.render(w, vehicleTemplate, vp)
		return
	}
	app.postMessage(w, "Vehicle profile added")
	return
}

//loadVehicleInputData retrieve current user vehicle data
func (app *appDB) loadVehicleInputData(userID string) vehiclePage {
	vp := vehiclePage{}

	db := data.InitAppDB(app.client, app.mDB, app.env)
	user, err := db.RetrieveProfile(userID)
	if err != nil {
		logger.Error.Println("userID not found")
		vp.Error = errors.New("Unable to process. Please re-login.")
		return vp
	}

	vp.Vehicle = user.Vehicle
	vp.Model = make(stringMap)
	for i, v := range app.dataCollections.vehicle {
		vp.Model[strconv.Itoa(i+1)] = v
	}

	vp.Charger = make(stringMap)
	for _, v := range app.dataCollections.charger {
		vp.Charger[v] = ""
	}
	for _, v := range user.ChargerType {
		vp.Charger[v] = "checked"
	}

	return vp
}
