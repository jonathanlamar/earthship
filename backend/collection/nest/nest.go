package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
    "os"
	"strings"
)

func main() {
	projectId := os.Getenv("PROJECT_ID") 
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN")

	newAccessToken := refreshAccessToken(clientId, clientSecret, refreshToken)
	fmt.Printf("New access token: %s\n", newAccessToken)

	deviceIds := getThermostatDeviceIds(projectId, newAccessToken)
	for _, deviceId := range deviceIds {
		fmt.Printf("Device Name: %s\n", deviceId)
	}
	if len(deviceIds) != 1 {
		log.Fatalf("Expected one device.  Found %v", deviceIds)
	}

	deviceId := deviceIds[0]
	thermostatReading := getThermostatReading(projectId, deviceId, newAccessToken)

	fmt.Printf("Thermostat Reading: %v\n", thermostatReading)
}

type GetAccessResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func getAccessToken(clientId, clientSecret, refreshToken string) (string, string) {
	fullUrl := "https://www.googleapis.com/oauth2/v4/token" +
		"?client_id=" + clientId +
		"&client_secret=" + clientSecret +
		"&grant_type=authorization_code&redirect_uri=https://www.google.com"

	body := sendEmptyPost(fullUrl)

	var result GetAccessResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Fatal("Can not unmarshal JSON")
	}

	return result.AccessToken, result.RefreshToken
}

func sendEmptyPost(url string) []byte {
	postBody, _ := json.Marshal(map[string]string{})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}

	return body
}

type RefreshAccessResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func refreshAccessToken(clientId, clientSecret, refreshToken string) string {
	fullUrl := "https://www.googleapis.com/oauth2/v4/token" +
		"?client_id=" + clientId +
		"&client_secret=" + clientSecret +
		"&refresh_token=" + refreshToken +
		"&grant_type=refresh_token"

	body := sendEmptyPost(fullUrl)

	var result RefreshAccessResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result.AccessToken
}

// Not complete
type GetDevicesResponse struct {
	Devices []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"devices"`
}

func getThermostatDeviceIds(projectId, accessToken string) []string {
	fullUrl := "https://smartdevicemanagement.googleapis.com/v1/enterprises/" + projectId + "/devices"
	body := sendGetRequestWithAccessToken(fullUrl, accessToken)

	var result GetDevicesResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	var deviceNames []string
	for _, device := range result.Devices {
		if device.Type == "sdm.devices.types.THERMOSTAT" {
			deviceNames = append(deviceNames, getDeviceNameFromPath(device.Name))
		}
	}

	return deviceNames
}

func sendGetRequestWithAccessToken(url, accessToken string) []byte {
	// Empty body.  Header has to be set after the request is created.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}

	return body
}

func getDeviceNameFromPath(device string) string {
	pathSplit := strings.Split(device, "/")

	return pathSplit[len(pathSplit)-1]
}

type ThermostatReading struct {
	Mode                string
	Temp                float64
	Humidity            int
	HeatSetpointCelsius float64
	CoolSetpointCelsius float64
	ThermostatMode      string
	FanTimerMode        string
	HvacStatus          string
}

type GetDeviceResponse struct {
	Name   string `json:"name"`
	Tyep   string `json:"type"`
	Traits struct {
		Humidity struct {
			AbientHumidityPercent int `json:"ambientHumidityPercent"`
		} `json:"sdm.devices.traits.Humidity"`
		Fan struct {
			TimerMode string `json:"timerMode"`
		} `json:"sdm.devices.traits.Fan"`
		ThermostatMode struct {
			Mode string `json:"mode"`
		} `json:"sdm.devices.traits.ThermostatMode"`
		HvacStatus struct {
			Status string `json:"status"`
		} `json:"sdm.devices.traits.ThermostatHvac"`
		ThermostatTemperatureSetpoint struct {
			HeatCelsius float64 `json:"heatCelsius"`
			CoolCelsius float64 `json:"coolCelsius"`
		} `json:"sdm.devices.traits.ThermostatTemperatureSetpoint"`
		Temperature struct {
			TemperatureCelsius float64 `json:"ambientTemperatureCelsius"`
		} `json:"sdm.devices.traits.Temperature"`
	} `json:"traits"`
}

func getThermostatReading(projectId, deviceId, accessToken string) ThermostatReading {
	fullUrl := "https://smartdevicemanagement.googleapis.com/v1/enterprises/" + projectId + "/devices/" + deviceId
	body := sendGetRequestWithAccessToken(fullUrl, accessToken)

	var result GetDeviceResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return getThermostatReadingFromResponse(result)
}

func getThermostatReadingFromResponse(response GetDeviceResponse) ThermostatReading {
	traits := response.Traits
	var thermostatReading ThermostatReading

	thermostatReading.Mode = traits.ThermostatMode.Mode
	thermostatReading.Temp = traits.Temperature.TemperatureCelsius
	thermostatReading.Humidity = traits.Humidity.AbientHumidityPercent
	thermostatReading.HeatSetpointCelsius = traits.ThermostatTemperatureSetpoint.HeatCelsius
	thermostatReading.CoolSetpointCelsius = traits.ThermostatTemperatureSetpoint.CoolCelsius
	thermostatReading.ThermostatMode = traits.ThermostatMode.Mode
	thermostatReading.FanTimerMode = traits.Fan.TimerMode
	thermostatReading.HvacStatus = traits.HvacStatus.Status

	return thermostatReading
}
