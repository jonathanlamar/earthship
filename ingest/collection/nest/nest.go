package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"strings"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event *LambdaRequest) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	projectId := event.ProjectId
	clientId := event.ClientId
	clientSecret := event.ClientSecret
	refreshToken := event.RefreshToken

	newAccessToken, err := refreshAccessToken(clientId, clientSecret, refreshToken)
	if err != nil {
		return nil, errors.New("Failed to refresh access token.")
	}
	log.Printf("New access token: %s\n", newAccessToken)

	deviceId, err := getThermostatDeviceId(projectId, newAccessToken)
	if err != nil {
		return nil, errors.New("Failed to get device ID")
	}
	thermostatReading, err := getThermostatReading(projectId, deviceId, newAccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to get thermostat reading from device with ID %v", deviceId)
	}

	message := fmt.Sprintf("Thermostat Reading: %+v\n", thermostatReading)
	return &message, nil
}

func refreshAccessToken(clientId, clientSecret, refreshToken string) (string, error) {
	fullUrl := "https://www.googleapis.com/oauth2/v4/token" +
		"?client_id=" + clientId +
		"&client_secret=" + clientSecret +
		"&refresh_token=" + refreshToken +
		"&grant_type=refresh_token"

	body, err := sendEmptyPost(fullUrl)
	if err != nil {
		return "", errors.New("Failed to send empty POST")
	}

	// Parse []byte to go struct pointer
	var result RefreshAccessResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", errors.New("Can not unmarshal JSON")
	}

	return result.AccessToken, nil
}

func getThermostatDeviceId(projectId, accessToken string) (string, error) {
	fullUrl := "https://smartdevicemanagement.googleapis.com/v1/enterprises/" + projectId + "/devices"
	body, err := sendGetRequestWithAccessToken(fullUrl, accessToken)
	if err != nil {
		return "", errors.New("Failed to set get request with access token.")
	}

	// Parse []byte to go struct pointer
	var result GetDevicesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", errors.New("Can not unmarshal JSON")
	}

	var deviceIds []string
	for _, device := range result.Devices {
		if device.Type == "sdm.devices.types.THERMOSTAT" {
			deviceIds = append(deviceIds, getDeviceNameFromPath(device.Name))
		}
	}
	if len(deviceIds) != 1 {
		return "", fmt.Errorf("Expected one device.  Found %v", deviceIds)
	}

	log.Printf("Device Id: %s\n", deviceIds[0])
	return deviceIds[0], nil
}

func getDeviceNameFromPath(device string) string {
	pathSplit := strings.Split(device, "/")
	return pathSplit[len(pathSplit)-1]
}

func getThermostatReading(projectId, deviceId, accessToken string) (*ThermostatReading, error) {
	fullUrl := "https://smartdevicemanagement.googleapis.com/v1/enterprises/" + projectId + "/devices/" + deviceId
	body, err := sendGetRequestWithAccessToken(fullUrl, accessToken)
	if err != nil {
		return nil, errors.New("Failed to send get request with access token.")
	}

	var result GetDeviceResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Printf("Can not unmarshal JSON")
	}

	thermostatReading := getThermostatReadingFromResponse(result)
	return &thermostatReading, nil
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

// TODO: This is never called.  Remove it in a future commit.
func getAccessToken(clientId, clientSecret, refreshToken string) (string, string, error) {
	fullUrl := "https://www.googleapis.com/oauth2/v4/token" +
		"?client_id=" + clientId +
		"&client_secret=" + clientSecret +
		"&grant_type=authorization_code&redirect_uri=https://www.google.com"

	body, err := sendEmptyPost(fullUrl)
	if err != nil {
		return "", "", errors.New("Failed to send empty POST")
	}

	var result GetAccessResponse
	// Parse []byte to golang struct pointer
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", errors.New("Can not unmarshal JSON")
	}

	return result.AccessToken, result.RefreshToken, nil
}
