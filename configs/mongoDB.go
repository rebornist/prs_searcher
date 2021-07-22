package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	// Timeout 설정을 위한 Context생성
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

	// Auth에러 처리를 위한 client option 구성
	clientOptions := options.Client().ApplyURI("mongodb://192.168.160.15:27017")

	// MongoDB 연결
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB Connection Made")
	return client, ctx, cancel
}

func GetDefaultCollection(client *mongo.Client) (collection *mongo.Collection) {
	return client.Database("prs_search").Collection("search")
}
