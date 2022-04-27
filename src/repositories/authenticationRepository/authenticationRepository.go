package authenticationRepository

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/xAutoBot/iqoption-sdk/src/configs"
)

type ResponseLogin struct {
	Code    string
	Ssid    string
	Message string
}

func GetSSID() (responseSsid chan string, responseError chan error) {
	responseSsid = make(chan string)
	responseError = make(chan error)

	go func() {

		requestBody, _ := json.Marshal(map[string]string{
			"identifier": configs.GetIqOptionEmail(),
			"password":   configs.GetIqOptionPassword(),
		})

		response, err := http.Post("https://"+configs.IqoptionAuthHOst+"/api/v2/login", "application/json", bytes.NewBuffer(requestBody))

		if err != nil {
			log.Printf("An Error Occured %v", err)
			responseSsid <- ""
			responseError <- err

			return
		}

		defer response.Body.Close()

		responseJson, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Printf("An Error Occured %v", err)
			responseSsid <- ""
			responseError <- err
			return
		}

		var responseLogin ResponseLogin
		err = json.Unmarshal(responseJson, &responseLogin)

		if err != nil {
			responseSsid <- ""
			responseError <- err
		}
		if responseLogin.Code != "success" {
			responseSsid <- ""
			responseError <- errors.New(responseLogin.Message)
			return
		}
		responseSsid <- string(responseLogin.Ssid)
		responseError <- nil
	}()
	return
}
