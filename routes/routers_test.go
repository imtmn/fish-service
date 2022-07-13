package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"mtmn.top/fish-service/entity"
)

func setup() {
	dir, _ := os.Getwd()
	fileName := dir + "/../.env"
	fmt.Println(fileName)
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestSetupRouter(t *testing.T) {
	r := SetupRouter()
	// other config

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

type (
	GoodsCase struct {
		Goods        entity.Goods
		ExceptCode   int
		ExceptRemark string
	}
)

var testIds []string

func TestSaveGoods(t *testing.T) {
	//用例数据
	var caseList []GoodsCase

	goods1 := entity.Goods{Name: "清蒸排骨", Remark: "test", Price: 2}
	caseList = append(caseList, GoodsCase{Goods: goods1, ExceptCode: 200, ExceptRemark: "正常数据"})

	goods2 := entity.Goods{Name: "红烧牛肉面", Remark: "test", Price: -1}
	caseList = append(caseList, GoodsCase{Goods: goods2, ExceptCode: 400, ExceptRemark: "价格不能为负数"})

	goods3 := entity.Goods{Name: "炒时蔬", Remark: "test", Price: 10000}
	caseList = append(caseList, GoodsCase{Goods: goods3, ExceptCode: 400, ExceptRemark: "价格不能超过1万，禁止杀猪"})

	goods4 := entity.Goods{Name: "水煮活鱼", Remark: "test", Price: 10}
	goods4.Specification = []entity.Specification{{Name: "微辣"}, {Name: "麻辣"}, {Name: "变态辣"}}
	caseList = append(caseList, GoodsCase{Goods: goods4, ExceptCode: 200, ExceptRemark: "价格不能为负数"})

	setup()
	r := SetupRouter()

	for _, item := range caseList {
		w := PostJson("/api/v1/goods/save", item.Goods, r)
		log.Println("body:", w.Body)
		assert.Equal(t, item.ExceptCode, w.Code)
		bodyBytes, _ := ioutil.ReadAll(w.Body)
		var body map[string]string
		json.Unmarshal(bodyBytes, &body)
		id := body["id"]
		if id != "" {
			testIds = append(testIds, id)
		}

	}
	log.Println(testIds)
}

func TestFindGoodsOne(t *testing.T) {
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/goods/"+testIds[0], nil)
	r.ServeHTTP(w, req)
	type Rep struct {
		Code int          `json:"code" `
		Msg  string       `json:"msg" `
		Data entity.Goods `json:"data" `
	}
	var re Rep
	bodyBytes, _ := ioutil.ReadAll(w.Body)
	log.Println(bodyBytes)
	json.Unmarshal(bodyBytes, &re)
	log.Println(re)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, testIds[0], re.Data.Id)
}

func TestFindGoodsPage(t *testing.T) {
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/goods", nil)

	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("size", "10")
	req.URL.RawQuery = q.Encode()
	r.ServeHTTP(w, req)
	type Rep struct {
		Code  int            `json:"code" `
		Msg   string         `json:"msg" `
		Data  []entity.Goods `json:"data" `
		Total int            `json:"total" `
		Page  int            `json:"page" `
	}
	var re Rep
	bodyBytes, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal(bodyBytes, &re)
	log.Println(re)
	assert.Equal(t, 200, w.Code)
	//查询条数应大于0
	assert.Less(t, 0, re.Total)
}

func TestDeleteGoods(t *testing.T) {
	r := SetupRouter()

	w := httptest.NewRecorder()

	type Rep struct {
		Code  int    `json:"code" `
		Msg   string `json:"msg" `
		Count int    `json:"count" `
	}

	for _, id := range testIds {
		req, _ := http.NewRequest("DELETE", "/api/v1/goods/"+id, nil)
		r.ServeHTTP(w, req)
		var re Rep
		bodyBytes, _ := ioutil.ReadAll(w.Body)
		json.Unmarshal(bodyBytes, &re)
		log.Println(re)
		assert.Equal(t, 200, w.Code)
		//查询条数应大于0
		assert.Less(t, 0, re.Count)
	}

}

//PostJson 根据特定请求uri和参数param，以Json形式传递参数，发起post请求返回响应
func PostJson(uri string, param interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	jsonByte, err := json.Marshal(param)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}
