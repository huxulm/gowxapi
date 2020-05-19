package helper

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackdon/gowxapi/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongodbURL = config.C.Mongo.MongodbURL
	client     *mongo.Client
)

func init() {
	log.Println("Start initialize mongodb client...", mongodbURL)
	// Set client options
	clientOptions := options.Client().ApplyURI(mongodbURL)

	// Connect to MongoDB
	_client, err := mongo.Connect(context.TODO(), clientOptions)

	client = _client
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
}

// GetClient provide a method for get *mongo.client
func GetClient() *mongo.Client {
	return client
}

// ConnectDB : This is helper function to connect mongoDB
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func ConnectDB(collection string) *mongo.Collection {
	return client.Database(config.C.Mongo.DbName).Collection(collection)
}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
