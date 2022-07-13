package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//测试环境变量设置
func TestSetup(t *testing.T) {
	dir, _ := os.Getwd()
	fileName := dir + "/../.env"
	fmt.Println(fileName)
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}
func TestConnectMongo(t *testing.T) {
	// 判断服务是不是可用
	err := DB().Client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("mongodb 连接失败")
	}
}

var insertResult *mongo.InsertOneResult

func TestInsert(t *testing.T) {
	insertResult = DB().InsertOne("test", bson.M{"name": "test", "remark": "test insert on"})
	if insertResult.InsertedID == nil {
		log.Fatal("inesert error")
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func TestFindOne(t *testing.T) {
	result := bson.M{}
	singleResult := DB().FindOne("test", "_id", insertResult.InsertedID)
	singleResult.Decode(&result)
	if result["name"] != "test" {
		log.Fatal("find on error")
	}
	log.Println(result)
}

func TestDelete(t *testing.T) {
	result := DB().Delete("test", "_id", insertResult.InsertedID)
	if result <= 0 {
		log.Fatal("delete on error")
	}
	log.Println(result)
}
