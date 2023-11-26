package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/seanburman/game-ws-server/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseOptions stores relevant database names.
type DatabaseOptions struct {
	Databases
	Collections
}

type Databases struct {
	Game string
}

type Collections struct {
	User string
}

type mongoClient struct {
	*mongo.Client
	DatabaseOptions
	Collections
}

func NewMongoClient(opts DatabaseOptions) *mongoClient {
	md := &mongoClient{DatabaseOptions: opts}
	if err := md.Connect(); err != nil {
		log.Panicln(err.Error())
	}
	return md
}

func (m *mongoClient) Connect() error {
	uri := config.Env().MONGO_URI
	t := reflect.TypeOf(bson.M{})
	reg := bson.NewRegistry()
	reg.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, t)
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri).SetTimeout(time.Second*30).SetRegistry(reg),
	)
	if err != nil {
		return err
	}
	m.Client = client
	fmt.Println("\n\033[32m Connected to MongoDB...\033[0m")
	return nil
}

// Disconnect() should be defered after calling Connect()
func (m *mongoClient) Disconnect() error {
	if err := m.Client.Disconnect(context.Background()); err != nil {
		return err
	}
	log.Println("Disconnected from MongoDB...")
	return nil
}
