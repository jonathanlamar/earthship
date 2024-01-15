package main

import "testing"

func TestGetDeviceIdThrowsIfTwoReturned(t *testing.T) {
	// TODO: Mock so get request returns two devices
	_, err := getThermostatDeviceId("foo", "bar")

	if err == nil {
		t.Fatal("Did not throw")
	}
	msg := err.Error()
	if msg != "Expected one device.  Found 2" {
		t.Fatalf("Did not throw the correct error:  Threw \"%s\"", msg)
	}
}

func TestGetDeviceIdThrowsIfZeroReturned(t *testing.T) {
	// TODO: Mock so get request returns empty list
	_, err := getThermostatDeviceId("foo", "bar")

	if err == nil {
		t.Fatal("Did not throw")
	}
	msg := err.Error()
	if msg != "Expected one device.  Found 0" {
		t.Fatalf("Did not throw the correct error:  Threw \"%s\"", msg)
	}
}

func TestGetDeviceIdFiltersForThermostatsOnly(t *testing.T) {
	// TODO: Mock so getDeviceId must filter
	expectedDeviceId := "1234-5678"
	deviceId, err := getThermostatDeviceId("foo", "bar")

	if err != nil {
		t.Fatalf("Threw \"%s\"", err.Error())
	}
	if deviceId != expectedDeviceId {
		t.Fatalf("Expected device ID %s, returned %s", expectedDeviceId, deviceId)
	}
}

func TestGetThermostatReadingParsesCorrectFields(t *testing.T) {
	// TODO: Mock response
	reading, err := getThermostatReading("foo", "bar", "baz")

	if err != nil {
		t.Fatalf("Threw \"%s\"", err.Error())
	}
	// Assert reading is not null and has correct fields
}
