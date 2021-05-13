package app

import (
	"golive/data"
	"golive/logger"
	"html/template"
)

type dataCollection struct {
	vehicle  []string
	charger  []string
	provider []string
}

//initTemplate initialise html/template
func initTemplate() *template.Template {
	return template.Must(template.ParseGlob("template/*.html"))
}

//initAppDB initialise appDB struct
func initAppDB(db *data.AppDB) *appDB {
	return &appDB{client: db.Client, mDB: db.MDB, env: db.Env}
}

//loadInitData load vehicles, charges and providers from database
//to be stored in memory array
func loadInitData(db *data.AppDB) dataCollection {
	dataCollections := dataCollection{}

	vehicles, err := db.RetrieveVehicleList()
	if err != nil {
		logger.Fatal.Println("loading vehicle list failed", err)
		return dataCollections
	}
	chargers, err := db.RetrieveChargerList()
	if err != nil {
		logger.Fatal.Println("loading charger list failed")
		return dataCollections
	}

	providers, err := db.RetrieveProviderList()
	if err != nil {
		logger.Fatal.Println("loading provider list failed")
		return dataCollections
	}

	for _, v := range vehicles {
		dataCollections.vehicle = append(dataCollections.vehicle, v.Model)
	}

	for _, v := range chargers {
		dataCollections.charger = append(dataCollections.charger, v.Type)
	}

	for _, v := range providers {
		dataCollections.provider = append(dataCollections.provider, v.Name)
	}

	return dataCollections
}
