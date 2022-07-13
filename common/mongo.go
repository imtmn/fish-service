package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mgo struct {
	Database string
	Client   *mongo.Client
}

var instance *mgo

var mu sync.Mutex

// public
func DB() *mgo {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = &mgo{
			Database: os.Getenv("MONGO_DATABASE"),
		}
		instance.Client = Connect()
	}
	return instance
}

// 连接设置
func Connect() *mongo.Client {
	log.Println("获取mongo连接")
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")

	var uri string
	if user != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?w=majority", user, pass,
			os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"), os.Getenv("MONGO_DATABASE"))
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s/%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"), os.Getenv("MONGO_DATABASE"))
		log.Println(uri)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(20)) // 连接池
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func (m *mgo) Coll(collection string) *mongo.Collection {
	coll, err := m.Client.Database(m.Database).Collection(collection).Clone()
	if err != nil {
		log.Fatal("获取mongo集合失败")
	}
	return coll
}

// 查询单个
func (m *mgo) FindOne(coll string, key string, value interface{}) *mongo.SingleResult {
	collection := m.Coll(coll)
	//collection.
	filter := bson.D{primitive.E{Key: key, Value: value}}
	singleResult := collection.FindOne(context.TODO(), filter)
	return singleResult
}

// 查询单个
func (m *mgo) FindID(coll string, id string) *mongo.SingleResult {
	collection := m.Coll(coll)
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter := bson.D{primitive.E{Key: "_id", Value: oid}}
		singleResult := collection.FindOne(context.TODO(), filter)
		return singleResult
	}
	//collection.
	return nil
}

//插入单个
func (m *mgo) InsertOne(coll string, value interface{}) *mongo.InsertOneResult {
	collection := m.Coll(coll)
	insertResult, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		fmt.Println(err)
	}
	return insertResult
}

//根据ID 编辑
func (m *mgo) UpdateByID(coll string, id string, update bson.D) *mongo.UpdateResult {
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		collection := m.Coll(coll)
		updateResult, err := collection.UpdateByID(context.TODO(), oid, update)
		if err != nil {
			fmt.Println(err)
		}
		return updateResult
	}
	return &mongo.UpdateResult{}
}

//根据ID 编辑
func (m *mgo) FindOneAndUpdate(coll string, filter interface{}, update bson.D) *mongo.SingleResult {
	return m.Coll(coll).FindOneAndUpdate(context.TODO(), filter, update)
}

//根据ID 编辑
func (m *mgo) UpdateOne(coll string, filter interface{}, update bson.D) *mongo.UpdateResult {
	updateResult, err := m.Coll(coll).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	return updateResult
}

//查询集合里有多少数据
func (m *mgo) CollectionCount(coll string) int64 {
	collection := m.Coll(coll)
	size, _ := collection.EstimatedDocumentCount(context.TODO())
	return size
}

//根据条件查询集合里有多少数据
func (m *mgo) CountByFilter(coll string, filter bson.M) int64 {
	collection := m.Coll(coll)
	size, _ := collection.CountDocuments(context.TODO(), filter)
	return size
}

//按选项查询集合 Skip 跳过 Limit 读取数量 sort 1 ，-1 . 1 为最初时间读取 ， -1 为最新时间读取
func (m *mgo) CollectionDocumentsFilter(coll string, Skip, Limit int64, filter bson.M, sort bson.D) *mongo.Cursor {
	collection := m.Coll(coll)
	// SORT := bson.D{primitive.E{Key: "_id", Value: sort}} //filter := bson.D{{key,value}}
	findOptions := options.Find().SetSort(sort).SetLimit(Limit).SetSkip(Skip)
	//findOptions.SetLimit(i)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	return temp
}

//按选项查询集合 Skip 跳过 Limit 读取数量 sort 1 ，-1 . 1 为最初时间读取 ， -1 为最新时间读取
func (m *mgo) CollectionDocuments(coll string, Skip, Limit int64, sort bson.D) *mongo.Cursor {
	collection := m.Coll(coll)
	// SORT := bson.D{primitive.E{Key: "_id", Value: sort}} //filter := bson.D{{key,value}}
	filter := bson.D{{}}
	findOptions := options.Find().SetSort(sort).SetLimit(Limit).SetSkip(Skip)
	//findOptions.SetLimit(i)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	return temp
}

//删除
func (m *mgo) DeleteById(coll string, id string) int64 {
	collection := m.Coll(coll)
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter := bson.D{primitive.E{Key: "_id", Value: oid}}
		count, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			return 0
		}
		return count.DeletedCount
	}
	return 0
}

//删除文章
func (m *mgo) Delete(coll string, key string, value interface{}) int64 {
	collection := m.Coll(coll)
	filter := bson.D{primitive.E{Key: key, Value: value}}
	count, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount
}

//删除多个
func (m *mgo) DeleteMany(coll string, key string, value interface{}) int64 {
	collection := m.Coll(coll)
	filter := bson.D{primitive.E{Key: key, Value: value}}

	count, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	return count.DeletedCount
}
