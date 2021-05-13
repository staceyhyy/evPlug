package data

import (
	"context"
	"golive/envvar"
	"golive/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type AppDB struct {
	Client *mongo.Client
	MDB    *mongo.Database
	Env    envvar.Env
}

//Connect connects to database and create indexes
func Connect(ctx context.Context, env envvar.Env) (*AppDB, error) {
	credential := options.Credential{
		Username: env.DBUser,
		Password: env.DBPwd,
	}
	db := &AppDB{}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.DBHost).SetAuth(credential))
	if err != nil {
		logger.Error.Println("db connect failed: ", err)
		return db, err
	}
	logger.Info.Println("db connected.")

	//ping the primary server
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error.Println("ping failed: ", err)
		return db, err
	}

	//create index for points collection
	db.MDB = client.Database(env.DBName)
	db.Env = env
	err = db.createLocationIndex(ctx)
	if err != nil {
		return db, err
	}

	//create index for users collection
	err = db.createUserIndex(ctx)
	if err != nil {
		return db, err
	}

	return db, nil
}

//createLocationIndex create indexes for points collection
func (db *AppDB) createLocationIndex(ctx context.Context) error {
	indexOpts := options.CreateIndexes().
		SetMaxTime(time.Second * 10)

	// creating index on point collection for field location 2dsphere type.
	pointIndexModel := mongo.IndexModel{
		Keys: bsonx.MDoc{"location": bsonx.String("2dsphere")},
	}

	pointIndexes := db.MDB.Collection(db.Env.DBPointColl).Indexes()
	_, err := pointIndexes.CreateOne(
		ctx,
		pointIndexModel,
		indexOpts,
	)
	if err != nil {
		logger.Error.Println("create location index ", err)
		return err
	}
	logger.Info.Println("index location created")

	// creating index on point collection for field address text type
	pointIndexModel = mongo.IndexModel{
		Keys: bsonx.MDoc{"address": bsonx.String("text")},
	}

	pointIndexes = db.MDB.Collection(db.Env.DBPointColl).Indexes()
	_, err = pointIndexes.CreateOne(ctx, pointIndexModel, indexOpts)
	if err != nil {
		logger.Error.Println("create address index ", err)
		return err
	}
	logger.Info.Println("index address created")
	return nil
}

//createUserIndex create indexes for users collection
func (db *AppDB) createUserIndex(ctx context.Context) error {
	indexOpts := options.CreateIndexes().
		SetMaxTime(time.Second * 10)

	// creating index on user collection for field userID text type
	userIndexModel := mongo.IndexModel{
		Keys: bsonx.MDoc{"userID": bsonx.String("text")},
	}

	userIndex := db.MDB.Collection(db.Env.DBUserColl).Indexes()
	_, err := userIndex.CreateOne(ctx, userIndexModel, indexOpts)
	if err != nil {
		logger.Error.Println("create user index ", err)
		return err
	}
	logger.Info.Println("index user created")
	return nil
}

//InitAppDB returns initialised struct
func InitAppDB(client *mongo.Client, mDB *mongo.Database, env envvar.Env) *AppDB {
	return &AppDB{client, mDB, env}
}
