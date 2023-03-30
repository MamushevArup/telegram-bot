package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Mongo struct{}

type Filter struct {
	TelegramTag string   `bson:"telegram_tag"`
	Region      []string `bson:"region"`
	StartPrice  int      `bson:"startPrice"`
	EndPrice    int      `bson:"endPrice"`
	IsOwner     string   `bson:"is_owner"`
	URL         string   `bson:"url"`
}

func (m *Mongo) Connection(tag string, region []string, startPrice, endPrice int, isOwner string, URL string) error {
	var client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	var db = client.Database("test")
	d := db.Collection("first")
	defer client.Disconnect(context.Background())
	if err != nil {
		return err
	}

	_, err = d.InsertOne(context.Background(), bson.M{"telegram_tag": tag, "region": region, "startPrice": startPrice, "endPrice": endPrice, "is_owner": isOwner, "url": URL})
	if err != nil {
		return err
	}
	return nil
}

func FindByFilter(tag string) (*Filter, error) {
	var client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	var db = client.Database("test")
	d := db.Collection("first")
	defer client.Disconnect(context.Background())
	if err != nil {
		return &Filter{}, err
	}
	var f *Filter
	filter := bson.M{"telegram_tag": tag}
	err = d.FindOne(context.Background(), filter).Decode(&f)
	//_, err = d.UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"url": lastUrl}})
	//fmt.Println(f, "())))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))")
	if err != nil {
		log.Println("Error while find", err)
	}
	return f, nil
}

func UpdateURL(tag, url string) error {
	var client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	var db = client.Database("test")
	d := db.Collection("first")
	filter := bson.M{"telegram_tag": tag}
	upd := bson.M{"$set": bson.M{"url": url}}
	_, err = d.UpdateMany(context.Background(), filter, upd)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())
	return nil
}
