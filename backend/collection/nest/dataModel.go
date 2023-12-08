package main

type LambdaRequest struct {
	ProjectId    string `json:"projectId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RefreshToken string `json:"refreshToken"`
}

type GetAccessResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type RefreshAccessResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// Not complete
type GetDevicesResponse struct {
	Devices []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"devices"`
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
