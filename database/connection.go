package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileAccess struct {
	Name   string      `bson:"name"`
	Access []time.Time `bson:"access"`
	Policy string      `bson:"policy"`
}

type FileAge struct {
	Name string    `bson:"name"`
	Age  time.Time `bson:"age"`
}

var client *mongo.Client

func ConnectToDb() {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

}

func CheckExists(filename string) (bool, []time.Time) {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, _ = mongo.Connect(context.TODO(), clientOptions)
	collection := client.Database("hdfs").Collection("accesses")
	filter := bson.M{"name": filename}
	var result FileAccess
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		log.Fatal(err)
	}

	return true, result.Access
}

func InsertAccessAndAge(filename string, ts time.Time, policy string) {
	collectionAcc := client.Database("hdfs").Collection("accesses")

	dataAcc := FileAccess{
		Name:   filename,
		Access: []time.Time{ts},
		Policy: policy,
	}

	_, err := collectionAcc.InsertOne(context.TODO(), dataAcc)
	if err != nil {
		log.Fatal(err)
	}

	collectionAge := client.Database("hdfs").Collection("ages")

	dataAge := FileAge{
		Name: filename,
		Age:  ts,
	}

	_, err = collectionAge.InsertOne(context.TODO(), dataAge)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateAccess(filename string, data []time.Time) {
	collectionAcc := client.Database("hdfs").Collection("accesses")

	filter := bson.M{"name": filename}
	update := bson.M{"$set": bson.M{"name": filename, "access": data}}

	_, err := collectionAcc.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

}

func UpdatePolicy(filename string, data []time.Time, policy string) {
	collectionAcc := client.Database("hdfs").Collection("accesses")

	filter := bson.M{"name": filename}
	update := bson.M{"$set": bson.M{"name": filename, "access": data, "policy": policy}}

	_, err := collectionAcc.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

}

func FetchFromDatabase() (map[string][]time.Time, map[string]time.Time, map[string]string) {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, _ = mongo.Connect(context.TODO(), clientOptions)
	collectionAcc := client.Database("hdfs").Collection("accesses")
	cursor, err := collectionAcc.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var files []FileAccess
	if err = cursor.All(context.TODO(), &files); err != nil {
		log.Fatal(err)
	}

	fileAccessMap := make(map[string][]time.Time)
	storagePolicyMap := make(map[string]string)
	for _, file := range files {
		fileAccessMap[file.Name] = file.Access
		storagePolicyMap[file.Name] = file.Policy
	}

	collectionAge := client.Database("hdfs").Collection("ages")
	cursor, err = collectionAge.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var filesAge []FileAge
	if err = cursor.All(context.TODO(), &filesAge); err != nil {
		log.Fatal(err)
	}

	fileAgeMap := make(map[string]time.Time)
	for _, file := range filesAge {
		fileAgeMap[file.Name] = file.Age
	}

	return fileAccessMap, fileAgeMap, storagePolicyMap
}
