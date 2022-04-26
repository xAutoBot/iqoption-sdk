package authenticationRepository

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/evecimar/iqoptionapi/src/configs"
)

type ResponseLogin struct {
	Code    string
	Ssid    string
	Message string
}

func GetSSID() (string, error) {

	requestBody, err := json.Marshal(map[string]string{
		"identifier": configs.GetIqOptionEmail(),
		"password":   configs.GetIqOptionPassword(),
	})

	response, err := http.Post("https://"+configs.IqoptionAuthHOst+"/api/v2/login", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return "", err
	}

	defer response.Body.Close()

	responseJson, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return "", err
	}

	var responseLogin ResponseLogin
	err = json.Unmarshal(responseJson, &responseLogin)
	log.Printf(string(responseJson))

	if responseLogin.Code != "success" {
		return "", errors.New(responseLogin.Message)
	}

	return string(responseLogin.Ssid), nil

}
