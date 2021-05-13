package data

import (
	"context"
	"golive/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Location is a GeoJSON type.
type Location struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type charger struct {
	Type    string `json:"type" bson:"type"`
	Detail  string `json:"detail" bson:"detail"`
	Price   string `json:"price" bson:"price"`
	Station string `json:"station" bson:"station"`
	Match   bool   `json:"match" bson:"match"`
}

type Point struct {
	Provider    string    `json:"provider" bson:"provider"`
	Address     string    `json:"address" bson:"address"`
	Postal      string    `json:"postal" bson:"postal"`
	Operator    string    `json:"operator" bson:"operator"`
	Requirement string    `json:"requirement" bson:"requirement"`
	Charger     []charger `json:"charger" bson:"charger"`
	Parking     string    `json:"parking" bson:"parking"`
	Hour        string    `json:"hour" bson:"hour"`
	Facility    string    `json:"facility" bson:"facility"`
	Location    Location  `json:"location" bson:"location"`
	Website     string    `json:"website" bson:"website"`
}

type PointJS struct {
	Provider string    `json:"provider"`
	Address  string    `json:"address"`
	Postal   string    `json:"postal"`
	Location []float64 `json:"location"`
	Charger  []charger `json:"charger"`
}

type PointInfoJS struct {
	Provider    string    `json:"provider"`
	Address     string    `json:"address"`
	Operator    string    `json:"operator"`
	Requirement string    `json:"requirement"`
	Charger     []charger `json:"charger"`
	Parking     string    `json:"parking"`
	Hour        string    `json:"hour"`
	Facility    string    `json:"facility"`
	Location    []float64 `json:"location"`
	Website     string    `json:"website"`
}

type Coordinate struct {
	Lat float64
	Lng float64
}

type Provider struct {
	Name string `json:"name" bson:"name"`
}

// NewPoint returns a GeoJSON Point with longitude and latitude.
func NewPoint(lat, lng float64) Location {
	return Location{
		"Point",
		[]float64{lng, lat},
	}
}

//RetrieveProviderList retrusn provider list
func (db *AppDB) RetrieveProviderList() ([]Provider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	providers := []Provider{}

	cur, err := db.MDB.Collection(db.Env.DBProviderColl).Find(ctx, bson.M{})
	if err != nil {
		logger.Error.Println("retrieveProvider-query: ", err)
		return providers, err
	}

	for cur.Next(ctx) {
		var p Provider
		err := cur.Decode(&p)
		if err != nil {
			logger.Error.Println("retrieveVehicle-decoding: ", err)
			return providers, err
		}
		providers = append(providers, p)
	}

	return providers, nil
}

//RetrieveAllPoints returns all point locations
func (db *AppDB) RetrieveAllPoints(chargerType []string) ([]PointJS, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []PointJS
	cur, err := db.MDB.Collection(db.Env.DBPointColl).Find(ctx, bson.M{})
	if err != nil {
		logger.Error.Println("retrieveAllPoints-query: ", err)
		return []PointJS{}, err
	}

	for cur.Next(ctx) {
		var p Point
		err := cur.Decode(&p)
		if err != nil {
			logger.Error.Println("retrieveAllPoints-decoding: ", err)
			return []PointJS{}, err
		}

		pJS := parsePoint(p, chargerType)
		results = append(results, pJS)
	}
	return results, nil
}

//RetrieveLocationInfo returns location info for selected location
func (db *AppDB) RetrieveLocationInfo(location Location, chargerType []string) (PointInfoJS, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results PointInfoJS
	filter := bson.D{
		primitive.E{Key: "location", Value: bson.D{
			primitive.E{Key: "$near", Value: bson.D{
				primitive.E{Key: "$geometry", Value: location},
			}},
		}},
	}

	cur, err := db.MDB.Collection(db.Env.DBPointColl).Find(ctx, filter)
	if err != nil {
		return PointInfoJS{}, err
	}
	for cur.Next(ctx) {
		var p Point
		err := cur.Decode(&p)
		if err != nil {
			logger.Error.Println("retrieveLocationInfo-decoding: ", err)
			return PointInfoJS{}, err
		}
		pJS := parsePointInfo(p, chargerType)
		return pJS, nil
	}

	return results, nil
}

//RetrieveInfoByAddr returns list of point location for selected address
func (db *AppDB) RetrieveInfoByAddr(address string, chargerType []string) ([]PointInfoJS, error) {
	results := []PointInfoJS{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := db.MDB.Collection(db.Env.DBPointColl).Find(ctx, bson.M{
		"$text": bson.M{
			"$search": address}})

	if err != nil {
		logger.Error.Println("retrieveInfoByAddr-query: ", err)
		return []PointInfoJS{}, err
	}
	for cur.Next(ctx) {
		var p Point
		err := cur.Decode(&p)
		if err != nil {
			logger.Error.Println("retrieveInfoByAddr-decoding: ", err)
			return []PointInfoJS{}, err
		}
		pJS := parsePointInfo(p, chargerType)
		results = append(results, pJS)
	}

	return results, nil
}

//RetrieveLocationByProvider returns list of providers for selected provider
func (db *AppDB) RetrieveLocationByProvider(provider string, chargerType []string) ([]PointInfoJS, error) {
	results := []PointInfoJS{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := db.MDB.Collection(db.Env.DBPointColl).Find(ctx, bson.M{
		// bson.M{"provider": provider})
		"$or": []bson.M{{"provider": provider}, {"operator": provider}}})

	if err != nil {
		logger.Error.Println("retrieveInfoByProvider-query: ", err)
		return []PointInfoJS{}, err
	}
	for cur.Next(ctx) {
		var p Point
		err := cur.Decode(&p)
		if err != nil {
			logger.Error.Println("retrieveInfoByProvider-decoding: ", err)
			return []PointInfoJS{}, err
		}
		pJS := parsePointInfo(p, chargerType)
		results = append(results, pJS)
	}

	return results, nil
}

//parsePoint returns selected fields
//to be used for returning json data
func parsePoint(p Point, chargerType []string) PointJS {
	pJS := PointJS{}

	pJS.Provider = p.Provider
	pJS.Address = p.Address
	pJS.Postal = p.Postal
	pJS.Location = append(pJS.Location, p.Location.Coordinates[1])
	pJS.Location = append(pJS.Location, p.Location.Coordinates[0])
	pJS.Charger = append(p.Charger)

	for _, v := range chargerType {
		for k, n := range pJS.Charger {
			if v == n.Type {
				pJS.Charger[k].Match = true
			}
		}
	}

	return pJS
}

//parsePoint returns selected fields
//to be used for returning json data
func parsePointInfo(p Point, chargerType []string) PointInfoJS {
	pJS := PointInfoJS{}

	pJS.Provider = p.Provider
	pJS.Address = p.Address
	pJS.Operator = p.Operator
	pJS.Requirement = p.Requirement
	pJS.Charger = p.Charger
	pJS.Parking = p.Parking
	pJS.Hour = p.Hour
	pJS.Facility = p.Facility
	pJS.Website = p.Website
	pJS.Location = append(pJS.Location, p.Location.Coordinates[1])
	pJS.Location = append(pJS.Location, p.Location.Coordinates[0])

	for _, v := range chargerType {
		for k, n := range pJS.Charger {
			if v == n.Type {
				pJS.Charger[k].Match = true
			}
		}
	}
	return pJS
}
