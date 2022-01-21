package main

import (
	"github.com/gin-gonic/gin"
	endpointV1 "github.com/raufhm/rtxTest/base/endpoint"
	"github.com/raufhm/rtxTest/base/util"
	"github.com/sirupsen/logrus"
)

func main() {

	// initialize connection to db
	db, err := util.GetPostgreClient()
	if err != nil {
		logrus.Error(err)
		return
	}
	defer db.Close()

	// bulk import
	err = util.BulkImport(db)
	if err != nil {
		logrus.Error(err)
		return
	}

	//initialize gin
	r := gin.Default()

	rates := r.Group("/rates")
	rates.GET("/latest", endpointV1.GetLatest)
	rates.GET("/:YYYY-MM-DD", endpointV1.GetByDate)
	rates.GET("/analyze", endpointV1.GetAnalize)

	r.Run()
}
