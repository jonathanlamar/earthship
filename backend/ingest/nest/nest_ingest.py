import logging

from ingest.nest import http
from ingest.nest.types import *

logger = logging.getLogger(__name__)


def HandleRequest(event: LambdaRequest, context: dict) -> str:
    projectId = event.projectId
    clientId = event.clientId
    clientSecret = event.clientSecret
    refreshToken = event.refreshToken

    newAccessToken = refreshAccessToken(clientId, clientSecret, refreshToken)
    logger.info(f"New access token: {newAccessToken}")

    deviceId = getThermostatDeviceId(projectId, newAccessToken)
    thermostatReading = getThermostatReading(projectId, deviceId, newAccessToken)

    return f"Thermostat Reading: {thermostatReading}"


def refreshAccessToken(clientId: str, clientSecret: str, refreshToken: str) -> str:
    fullUrl = (
        "https://www.googleapis.com/oauth2/v4/token"
        + "?client_id="
        + clientId
        + "&client_secret="
        + clientSecret
        + "&refresh_token="
        + refreshToken
        + "&grant_type=refresh_token"
    )

    body = http.sendEmptyPost(fullUrl)
    result = RefreshAccessResponse.model_validate(body)

    return result.accessToken


def getThermostatDeviceId(projectId: str, accessToken: str) -> str:
    fullUrl = f"https://smartdevicemanagement.googleapis.com/v1/enterprises/{projectId}/devices"
    body = http.sendGetRequestWithAccessToken(fullUrl, accessToken)
    result = GetDevicesResponse.model_validate(body)

    deviceIds = [
        getDeviceNameFromPath(device.name)
        for device in result.devices
        if device.deviceType == "sdm.devices.types.THERMOSTAT"
    ]
    if len(deviceIds) != 1:
        raise RuntimeError(f"Expected one device.  Found {len(deviceIds)}.")

    logger.info(f"Device Id: {deviceIds[0]}")
    return deviceIds[0]


def getDeviceNameFromPath(device: str) -> str:
    return device.split("/")[-1]


def getThermostatReading(projectId: str, deviceId: str, accessToken: str) -> ThermostatReading:
    fullUrl = f"https://smartdevicemanagement.googleapis.com/v1/enterprises/{projectId}/devices/{deviceId}"
    body = http.sendGetRequestWithAccessToken(fullUrl, accessToken)
    result = GetDeviceResponse.model_validate(body)

    return ThermostatReading(
        temp=result.traits.temperature.ambientTemperatureCelsius,
        humidity=result.traits.humidity.ambientHumidityPercent,
        heatSetpointCelsius=result.traits.thermostatTemperatureSetpoint.heatCelsius,
        coolSetpointCelsius=result.traits.thermostatTemperatureSetpoint.coolCelsius,
        thermostatMode=result.traits.thermostatMode.mode,
        fanTimerMode=result.traits.fan.timerMode,
        hvacStatus=result.traits.hvacStatus.status,
    )


# TODO: This is never called.  Remove it in a future commit.
def getAccessToken(clientId: str, clientSecret: str) -> tuple[str, str]:
    fullUrl = (
        "https://www.googleapis.com/oauth2/v4/token"
        + "?client_id="
        + clientId
        + "&client_secret="
        + clientSecret
        + "&grant_type=authorization_code&redirect_uri=https://www.google.com"
    )

    body = http.sendEmptyPost(fullUrl)
    result = GetAccessResponse.model_validate(body)

    return result.accessToken, result.refreshToken
