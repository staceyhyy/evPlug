package data

import (
	"context"
	"golive/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Vehicle struct {
	Model string `json:"model" bson:"model"`
}

type Charger struct {
	Type string `json:"type" bson:"type"`
}

//RetrieveVehicleList returns vehicle list
func (db *AppDB) RetrieveVehicleList() ([]Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	vehicles := []Vehicle{}

	cur, err := db.MDB.Collection(db.Env.DBVehicleColl).Find(ctx, bson.M{})
	if err != nil {
		logger.Error.Println("retrieveVehicle-query: ", err)
		return vehicles, err
	}

	for cur.Next(ctx) {
		var v Vehicle
		err := cur.Decode(&v)
		if err != nil {
			logger.Error.Println("retrieveVehicle-decoding: ", err)
			return vehicles, err
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

//RetrieveChargerList returns charger list
func (db *AppDB) RetrieveChargerList() ([]Charger, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	types := []Charger{}
	cur, err := db.MDB.Collection(db.Env.DBChargerColl).Find(ctx, bson.M{})
	if err != nil {
		logger.Error.Println("retriveChargeList-query: ", err)
		return types, err
	}

	for cur.Next(ctx) {
		var p Charger
		err := cur.Decode(&p)
		if err != nil {
			logger.Error.Println("retrieveVehicle-decoding: ", err)
			return types, err
		}
		types = append(types, p)
	}
	return types, nil
}
