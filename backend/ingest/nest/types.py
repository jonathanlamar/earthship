from pydantic import BaseModel, ConfigDict, Field


class SerializableDataClass(BaseModel):
    model_config = ConfigDict(populate_by_name=True)


class LambdaRequest(SerializableDataClass):
    projectId: str
    clientId: str
    clientSecret: str
    refreshToken: str


class GetAccessResponse(SerializableDataClass):
    accessToken: str
    expiresIn: int
    refreshToken: str
    scope: str
    tokenType: str


class RefreshAccessResponse(SerializableDataClass):
    accessToken: str
    expiresIn: int
    scope: str
    tokenType: str


# Not complete
class Device(SerializableDataClass):
    name: str
    deviceType: str


class GetDevicesResponse(SerializableDataClass):
    devices: list[Device]


class ThermostatReading(SerializableDataClass):
    temp: float
    humidity: int
    heatSetpointCelsius: float
    coolSetpointCelsius: float
    thermostatMode: str
    fanTimerMode: str
    hvacStatus: str


class Humidity(SerializableDataClass):
    ambientHumidityPercent: int


class Fan(SerializableDataClass):
    timerMode: str


class ThermostatMode(SerializableDataClass):
    mode: str


class HvacStatus(SerializableDataClass):
    status: str


class ThermostatTemperatureSetpoint(SerializableDataClass):
    heatCelsius: float
    coolCelsius: float


class Temperature(SerializableDataClass):
    ambientTemperatureCelsius: float


class Traits(SerializableDataClass):
    humidity: Humidity = Field(..., alias="sdm.devices.traits.Humidity")
    fan: Fan = Field(..., alias="sdm.devices.traits.Fan")
    thermostatMode: ThermostatMode = Field(..., alias="sdm.devices.traits.ThermostatMode")
    hvacStatus: HvacStatus = Field(..., alias="sdm.devices.traits.ThermostatHvac")
    thermostatTemperatureSetpoint: ThermostatTemperatureSetpoint = Field(
        ..., alias="sdm.devices.traits.ThermostatTemperatureSetpoint"
    )
    temperature: Temperature = Field(..., alias="sdm.devices.traits.Temperature")


class GetDeviceResponse(SerializableDataClass):
    name: str
    deviceType: str = Field(..., alias="type")
    traits: Traits
