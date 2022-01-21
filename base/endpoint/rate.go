package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	src "github.com/raufhm/rtxTest/base"
	util "github.com/raufhm/rtxTest/base/util"
	"github.com/sirupsen/logrus"
)

func GetLatest(c *gin.Context) {
	db, err := util.GetPostgreClient()
	if err != nil {
		logrus.Error(err)
		return
	}

	statement := `SELECT currency, rate FROM (SELECT * FROM envelopes ORDER BY time DESC LIMIT 35) AS envs ORDER BY currency ASC`
	envelopes, err := db.Query(statement)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	dataEnvelopes := []*src.Cube3{}
	for envelopes.Next() {
		e := &src.Cube3{}
		err = envelopes.Scan(&e.Currency, &e.Rate)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		dataEnvelopes = append(dataEnvelopes, e)
	}

	if len(dataEnvelopes) == 0 {
		err = fmt.Errorf("an error occured")
		logrus.Error(err)
		c.JSON(http.StatusNotFound, err)
		return
	}
	currencyRateMap := map[string]float64{}
	for _, r := range dataEnvelopes {
		rate, _ := strconv.ParseFloat(r.Rate, 64)
		currencyRateMap[r.Currency] = rate
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"base":  "EUR",
		"rates": currencyRateMap,
	})

}

func GetByDate(c *gin.Context) {
	requestDate, err := time.Parse("2006-01-02", c.Param("YYYY-MM-DD"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{})
		return
	}

	db, err := util.GetPostgreClient()
	if err != nil {
		logrus.Error(err)
		return
	}

	statement := `SELECT currency, rate FROM (SELECT * FROM envelopes WHERE time = $1 LIMIT 35) AS envs ORDER BY currency ASC`
	envelopes, err := db.Query(statement, requestDate)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	dataEnvelopes := []*src.Cube3{}
	for envelopes.Next() {
		e := &src.Cube3{}
		err = envelopes.Scan(&e.Currency, &e.Rate)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{})
			return
		}
		dataEnvelopes = append(dataEnvelopes, e)
	}

	if len(dataEnvelopes) == 0 {
		err = fmt.Errorf("an error occured")
		logrus.Error(err)
		c.JSON(http.StatusNotFound, err)
		return
	}

	currencyRateMap := map[string]float64{}
	for _, r := range dataEnvelopes {
		rate, _ := strconv.ParseFloat(r.Rate, 64)
		currencyRateMap[r.Currency] = rate
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"base":  "EUR",
		"rates": currencyRateMap,
	})

}

func GetAnalize(c *gin.Context) {
	db, err := util.GetPostgreClient()
	if err != nil {
		logrus.Error(err)
		return
	}

	statement := `SELECT currency, MIN(rate) AS min_rate, MAX(rate) AS max_rate, AVG(rate) AS avg_rate FROM envelopes GROUP BY currency ORDER BY currency`
	envelopes, err := db.Query(statement)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	dataEnvelopes := []*src.Analyze{}
	for envelopes.Next() {
		e := &src.Analyze{}
		err = envelopes.Scan(&e.Currency, &e.Min, &e.Max, &e.Avg)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{})
			return
		}
		dataEnvelopes = append(dataEnvelopes, e)
	}

	if len(dataEnvelopes) == 0 {
		err = fmt.Errorf("an error occured")
		logrus.Error(err)
		c.JSON(http.StatusNotFound, err)
		return
	}

	currencyRateMap := map[string]map[string]interface{}{}
	for _, r := range dataEnvelopes {
		currencyRateMap[r.Currency] = map[string]interface{}{
			"MIN": r.Min,
			"MAX": r.Max,
			"AVG": r.Avg,
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"base":  "EUR",
		"rates": currencyRateMap,
	})
}


