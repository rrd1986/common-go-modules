package utils

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

const responseMismatchPrintFormat = "Actual response is not same as Expected data"

func TestRetrievePayloadMissingHeaderExpectJsonPayloadReturnsStruct(t *testing.T) {

	payloadData := `{"Feb":{"Fullform":"February"},"Jan":{"Fullform":"January"},"Mar":{"Fullform":"March"}}`
	payload := []byte(payloadData)
	request, _ := http.NewRequest("POST", "/sample/request", bytes.NewBuffer(payload))

	// sending empty map as DB record
	actual, err := RetrievePayload(request)

	// check if data is same as expected
	if err != nil {
		t.Errorf("Not expecting an exception in this test case")
	}

	expected := map[string]interface{}{
		"Feb": map[string]interface{}{"Fullform": "February"},
		"Jan": map[string]interface{}{"Fullform": "January"},
		"Mar": map[string]interface{}{"Fullform": "March"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(responseMismatchPrintFormat)
	}
}

func TestRetrievePayloadJsonPayloadReturnsStruct(t *testing.T) {

	payloadData := `{"Feb":{"Fullform":"February"},"Jan":{"Fullform":"January"},"Mar":{"Fullform":"March"}}`
	payload := []byte(payloadData)
	request, _ := http.NewRequest("POST", "/sample/request", bytes.NewBuffer(payload))
	request.Header.Set(ContentType, JSONContentType)

	// sending empty map as DB record
	actual, err := RetrievePayload(request)

	// check if data is same as expected
	if err != nil {
		t.Errorf("Not expecting an exception in this test case")
	}

	expected := map[string]interface{}{
		"Feb": map[string]interface{}{"Fullform": "February"},
		"Jan": map[string]interface{}{"Fullform": "January"},
		"Mar": map[string]interface{}{"Fullform": "March"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(responseMismatchPrintFormat)
	}
}

func TestRetrievePayloadXmlPayloadThrowsError(t *testing.T) {

	payloadData := `{"Feb":{"Fullform":"February"},"Jan":{"Fullform":"January"},"Mar":{"Fullform":"March"}}`
	payload := []byte(payloadData)
	request, _ := http.NewRequest("POST", "/sample/request", bytes.NewBuffer(payload))
	request.Header.Set(ContentType, XMLContentType)

	// sending empty map as DB record
	actual, err := RetrievePayload(request)

	// check if data is same as expected
	if err == nil {
		t.Errorf("Expecting an exception in this test case")
	}

	if actual != "" {
		t.Errorf(responseMismatchPrintFormat)
	}
}
