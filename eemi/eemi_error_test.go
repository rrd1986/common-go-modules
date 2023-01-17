package eemi

import (
	"reflect"
	"testing"
)

// Test that the provided details are reflected in the error. Not testing the error content here.
// The translated error will be the responsibility of the api
func TestCreateNGCIError(t *testing.T) {
	ex := New(nil, "NAST00001")

	expectedFormat := map[string]interface{}{
		"messageID": "NAST00001",
	}
	if !reflect.DeepEqual(expectedFormat, ex.FormatError()) {
		t.Error("Expecting the error formatted message details to match!")
	}
}
