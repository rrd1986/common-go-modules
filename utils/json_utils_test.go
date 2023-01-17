package utils

import (
	"strings"
	"testing"
)

var request = []map[string]interface{}{
	{
		"interface":            1,
		"switchPortMode":       "Access",
		"vlans":                []int{102},
		"spanningTreeProtocol": "",
		"cos":                  "",
		"mtu":                  15000,
		"speed":                "",
		"negotiate":            "",
		"templateName":         "ngci_idrac_template",
		"interfaceState":       "Noshutdown",
	},
	{
		"interface":            2,
		"switchPortMode":       "Access",
		"vlans":                []int{102},
		"spanningTreeProtocol": "Edge",
		"cos":                  "",
		"mtu":                  1500,
		"speed":                "",
		"negotiate":            "",
		"templateName":         "ngci_idrac_template",
		"interfaceState":       "Noshutdown",
	},
}

func TestJSONMarshallSuccess(t *testing.T) {

	response, err := JSONMarshall(request)

	if err != nil {
		t.Error("Not expecting an exception in this test case")
	}

	expected := `[{"cos":"","interface":1,"interfaceState":"Noshutdown","mtu":15000,"negotiate":"","spanningTreeProtocol":"","speed":"","switchPortMode":"Access","templateName":"ngci_idrac_template","vlans":[102]},{"cos":"","interface":2,"interfaceState":"Noshutdown","mtu":1500,"negotiate":"","spanningTreeProtocol":"Edge","speed":"","switchPortMode":"Access","templateName":"ngci_idrac_template","vlans":[102]}]`
	if strings.Compare(string(response), expected) != 0 {
		t.Errorf("Expecting response to match expected information, response to be set to '%s'. Got '%s'", expected, string(response))
	}
}
