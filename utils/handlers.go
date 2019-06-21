package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func doRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FetchHandler ...
func FetchHandler(c *gin.Context) {
	params := c.Request.URL.Query()
	value := params["city"][0]
	respBody := map[string]interface{}{}
	if strings.Compare(value, "") == 0 {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric,uk&APPID=2647d7f2dab4bd9d604b908204051a80", value)
	wResp, err := doRequest(url)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	fmt.Println("Request to wapi done")
	if wResp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(wResp.Body)

		json.Unmarshal(bodyBytes, &respBody)
		main := respBody["main"].(map[string]interface{})
		kTemp := main["temp"].(float64)

		cResp, cErr := doRequest(fmt.Sprintf("http://convertor/convert?value=%f", kTemp))

		if cErr != nil {
			c.Status(http.StatusInternalServerError)
			fmt.Println("Request to convertor failed: ", cErr)
			return
		}

		var convertedResp map[string]interface{}
		cBytes, _ := ioutil.ReadAll(cResp.Body)
		defer cResp.Body.Close()
		json.Unmarshal(cBytes, &convertedResp)

		c.JSON(200, map[string]interface{}{
			"status": http.StatusOK,
			"temp":   convertedResp["cvalue"],
		})
	} else {
		c.Status(http.StatusInternalServerError)
	}

}

// Health ...
func Health(c *gin.Context) {
	c.Status(http.StatusOK)
}
