package repository

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/danilomarques1/godemo/provider/api/model"
	"github.com/danilomarques1/godemo/provider/api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			apiErr := util.NewApiError("Cob not found", http.StatusNotFound)
			return nil, apiErr
		}
		return nil, err
	}
	return cob, nil
}

func (cmr *CobMongoRepository) Update(cob *model.Cob) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: cob.Status}}}}
	_, err := cmr.collection.UpdateByID(context.Background(), cob.TxId, update, options.Update())
	if err != nil {
		return err
	}
	return nil
}
