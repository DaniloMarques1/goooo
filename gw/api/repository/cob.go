package repository

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/danilomarques1/godemo/gw/api/model"
	"github.com/danilomarques1/godemo/gw/api/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CobMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewCobMongoRepository(client *mongo.Client, collName string) *CobMongoRepository {
	dbName := os.Getenv("DATABASE")
	collection := client.Database(dbName).Collection(collName)
	return &CobMongoRepository{client: client, collection: collection}
}

func (cmr *CobMongoRepository) Save(cob *model.Cob) error {
	_, err := cmr.collection.InsertOne(context.Background(), cob)
	if err != nil {
		return err
	}
	return nil
}

func (cmr *CobMongoRepository) FindById(txid string) (*model.Cob, error) {
	var cob *model.Cob
	filter := bson.D{{Key: "_id", Value: txid}}
	if err := cmr.collection.FindOne(context.Background(), filter).Decode(&cob); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, response.NewApiError("Cob not found", http.StatusNotFound)
		}

		return nil, err
	}
	return cob, nil
}
