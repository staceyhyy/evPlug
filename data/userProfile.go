package data

import (
	"context"
	"golive/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Email       string   `json:"email" bson:"email"`
	UserID      string   `json:"userID" bson:"userID"`
	Password    []byte   `json:"password" bson:"password"`
	Vehicle     string   `json:"vehicle" bson:"vehicle"`
	ChargerType []string `json:"chargerType" bson:"chargerType"`
	UseLocation bool     `json:"useLocation" bson:"useLocation"`
}

//RetrieveProfile returns user profile for the requested userID
func (db *AppDB) RetrieveProfile(userID string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := User{}
	err := db.MDB.Collection(db.Env.DBUserColl).FindOne(ctx, bson.M{
		"$text": bson.M{
			"$search": userID}}).Decode(&result)

	if err != nil {
		logger.Info.Println("retrieveUser-query: ", err)
		return result, err
	}
	return result, nil
}

//SaveProfile saves user profile to the database
func (db *AppDB) SaveProfile(usr User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.MDB.Collection(db.Env.DBUserColl).InsertOne(ctx, usr)
	if err != nil {
		logger.Error.Println("saveUser-insert: ", err)
		return err
	}
	logger.Info.Println("userID " + usr.UserID + " added successfully")
	return nil
}

//UpdateUserPwd updates user's password
func (db *AppDB) UpdateUserPwd(userID string, pwd []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.MDB.Collection(db.Env.DBUserColl).UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"password": pwd}},
	)

	if err != nil {
		logger.Error.Println("update password in user coll: ", err)
		return err
	}
	return nil
}

//UpdateUserPwd updates user's email
func (db *AppDB) UpdateUserEmail(userID, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.MDB.Collection(db.Env.DBUserColl).UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"email": email}},
	)

	if err != nil {
		logger.Error.Println("update email in user coll: ", err)
		return err
	}
	return nil
}

//UpdateUserVehicle updates user's vehicle info
func (db *AppDB) UpdateUserVehicle(userID string, vehicle string, chargerType []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.MDB.Collection(db.Env.DBUserColl).UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"vehicle": vehicle, "chargerType": chargerType}},
	)

	if err != nil {
		logger.Error.Println("update vehicle in user coll: ", err)
		return err
	}
	return nil
}

//UpdateUserView updates 'use my current location' setting
func (db *AppDB) UpdateUserView(userID string, useLocation bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.MDB.Collection(db.Env.DBUserColl).UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"useLocation": useLocation}},
	)

	if err != nil {
		logger.Error.Println("update view in user coll: ", err)
		return err
	}
	return nil
}
