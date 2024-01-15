import unittest
from unittest.mock import patch

from ingest.nest.nest_ingest import getThermostatDeviceId, getThermostatReading
from ingest.nest.types import *


class TestNestIngest(unittest.TestCase):
    deviceOne = {"name": "enterprises/123-456/devices/foo", "deviceType": "sdm.devices.types.THERMOSTAT"}
    deviceTwo = {"name": "enterprises/123-456/devices/bar", "deviceType": "sdm.devices.types.THERMOSTAT"}
    deviceThree = {"name": "enterprises/123-456/devices/baz", "deviceType": "sdm.devices.types.SOMETHINGELSE"}
    genericDeviceResponse = {
        "name": "enterprises/123-456/devices/foo",
        "type": "sdm.devices.types.THERMOSTAT",
        "assignee": "enterprises/123-456/structures/blahblah",
        "traits": {
            "sdm.devices.traits.Info": {"customName": ""},
            "sdm.devices.traits.Humidity": {"ambientHumidityPercent": 49},
            "sdm.devices.traits.Connectivity": {"status": "ONLINE"},
            "sdm.devices.traits.Fan": {"timerMode": "OFF"},
            "sdm.devices.traits.ThermostatMode": {
                "mode": "HEATCOOL",
                "availableModes": ["HEAT", "COOL", "HEATCOOL", "OFF"],
            },
            "sdm.devices.traits.ThermostatEco": {
                "availableModes": ["OFF", "MANUAL_ECO"],
                "mode": "OFF",
                "heatCelsius": 10,
                "coolCelsius": 24.444443,
            },
            "sdm.devices.traits.ThermostatHvac": {"status": "OFF"},
            "sdm.devices.traits.Settings": {"temperatureScale": "FAHRENHEIT"},
            "sdm.devices.traits.ThermostatTemperatureSetpoint": {"heatCelsius": 19.152817, "coolCelsius": 23.333328},
            "sdm.devices.traits.Temperature": {"ambientTemperatureCelsius": 19.119995},
        },
        "parentRelations": [{"parent": "enterprises/123-456/structures/blahblahblah", "displayName": "Living Room"}],
    }

    @patch("ingest.nest.http.sendGetRequestWithAccessToken", return_value={"devices": [deviceOne, deviceTwo]})
    def test_getDeviceIdThrowsIfTwoDevicesReturned(self, mock_sendGetRequestWithAccessToken):
        with self.assertRaisesRegex(RuntimeError, "Expected one device.  Found 2."):
            getThermostatDeviceId(projectId="123", accessToken="ABCDDEF")

    @patch("ingest.nest.http.sendGetRequestWithAccessToken", return_value={"devices": []})
    def test_getDeviceIdThrowsIfZeroDevicesReturned(self, mock_sendGetRequestWithAccessToken):
        with self.assertRaisesRegex(RuntimeError, "Expected one device.  Found 0."):
            getThermostatDeviceId(projectId="123", accessToken="ABCDDEF")

    @patch("ingest.nest.http.sendGetRequestWithAccessToken", return_value={"devices": [deviceOne, deviceThree]})
    def test_getDeviceIdFiltersForThermostatsOnly(self, mock_sendGetRequestWithAccessToken):
        expectedDeviceId = "foo"
        deviceId = getThermostatDeviceId(projectId="123", accessToken="ABCDDEF")

        self.assertEqual(deviceId, expectedDeviceId)

    @patch("ingest.nest.http.sendGetRequestWithAccessToken", return_value=genericDeviceResponse)
    def test_getThermostatReadingParsesCorrectFields(self, mock_sendGetRequestWithAccessToken):
        expectedReading = ThermostatReading(
            temp=19.119995,
            humidity=49,
            heatSetpointCelsius=19.152817,
            coolSetpointCelsius=23.333328,
            thermostatMode="HEATCOOL",
            fanTimerMode="OFF",
            hvacStatus="OFF",
        )
        reading = getThermostatReading(projectId="123", deviceId="foo", accessToken="ABCDEF")

        self.assertEqual(reading, expectedReading)
