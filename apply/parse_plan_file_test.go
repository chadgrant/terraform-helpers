package main

import "testing"

func TestParsePlanFileName(t *testing.T) {
	testParse(t, &planfile{true, "development", "sms-send-api", "api"}, "destroy_development_sms-send-api_api.plan")
	testParse(t, &planfile{false, "development", "sms-send-api", "api"}, "development_sms-send-api_api.plan")
	testParse(t, &planfile{false, "development", "", "network"}, "development_network.plan")
	testParse(t, &planfile{true, "development", "", "network"}, "destroy_development_network.plan")

	testParse(t, &planfile{true, "development", "service_with_underscores", "network"}, "destroy_development_service_with_underscores_network.plan")
}

func testParse(t *testing.T, expected *planfile, filename string) {

	parsed, err := parsePlanFile(filename)
	if err != nil {
		t.Error(err)
		return
	}

	if parsed.destroy != expected.destroy {
		t.Error("Should be destroy")
	}

	if len(parsed.environment) <= 0 && parsed.environment != expected.environment {
		t.Error("did not parse environment")
	}

	if len(parsed.service) <= 0 && parsed.service != expected.service {
		t.Error("did not parse service")
	}

	if len(parsed.stack) <= 0 && parsed.stack == expected.stack {
		t.Error("did not parse stack")
	}
}
