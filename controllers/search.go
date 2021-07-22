package controllers

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"prsSearcher/configs"
	"prsSearcher/customTypes"
	"prsSearcher/middlewares"
	"prsSearcher/models"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Search(c echo.Context) error {

	// logger 접속
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	// mongodb 접속
	client := c.Request().Context().Value("CLIENT").(*mongo.Client)

	// Response 생성
	var resp customTypes.Response

	// Keyword 문장 변수에 담기
	keyword := c.QueryParam("keyword")

	condition := bson.D{}
	if keyword != "" {
		// 조건 생성
		condition = bson.D{{"content", primitive.Regex{Pattern: keyword, Options: ""}}}
	}

	// search 모델 생성
	var searchs []models.Search

	// Declare an index model object to pass to CreateOne()
	// db.members.createIndex( { "SOME_FIELD": 1 }, { unique: true } )
	mod := mongo.IndexModel{
		Keys: bson.M{
			"code": 1, // index in ascending order
		}, Options: nil,
	}

	col := configs.GetDefaultCollection(client)
	_, err := col.Indexes().CreateOne(context.TODO(), mod)

	// Check if the CreateOne() method returned any errors
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = nil
		middlewares.CreateLogger(logger, http.StatusInternalServerError, err.Error())
		return c.JSON(http.StatusInternalServerError, resp)
	}

	res, err := col.Find(context.TODO(), condition)
	// Check if the CreateOne() method returned any errors
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = nil
		middlewares.CreateLogger(logger, http.StatusInternalServerError, err.Error())
		return c.JSON(http.StatusInternalServerError, resp)
	}

	// 결과를 변수에 담기
	if err = res.All(context.TODO(), &searchs); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		resp.Data = nil
		middlewares.CreateLogger(logger, http.StatusInternalServerError, err.Error())
		return c.JSON(http.StatusInternalServerError, resp)
	}

	data := make(map[string][]string)
	for _, v := range searchs {

		// 검색 결과 정제
		key := fmt.Sprintf("%s||%s", v.Code, v.Keyword)
		value := searchKey(data, key)
		data[key] = append(value, v.Content)

	}

	var resps []customTypes.SearchResponse
	resps = sortData(data, resps)

	resp.Code = http.StatusOK
	resp.Message = fmt.Sprintf("%d keywords search success", len(searchs))
	resp.Data = resps
	middlewares.CreateLogger(logger, http.StatusOK, resp.Message)
	return c.JSON(http.StatusOK, resp)
}

func searchKey(data map[string][]string, val string) []string {
	value := data[val]
	return value
}

func sortData(data map[string][]string, resps []customTypes.SearchResponse) []customTypes.SearchResponse {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// To perform the opertion you want
	for _, k := range keys {
		keys := strings.Split(k, "||")
		var response customTypes.SearchResponse
		response.Code = keys[0]
		response.Keyword = keys[1]
		response.Datas = data[k]
		resps = append(resps, response)
	}

	return resps
}
