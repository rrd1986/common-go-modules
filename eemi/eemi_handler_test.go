package eemi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/middleware"
	"github.com/rrd1986/common-go-modules/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// test that a handler that just returns an eemi error gets intercepted and the body is returned as an eemi error
func Test_EEMIHandler_Returns_EEMI_Response(t *testing.T) {
	currentFolder, _ := os.Getwd()
	SetLogger(log.NewLogger("", ""))
	eemiData, _ := LoadEemiFromFile(currentFolder+"/testData/eemi_test_data.json", utils.FileSystem{})

	handler := NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		return NewWithTemplateData(nil, "NGCI0002", map[string]string{"Name": "Stark", "Season": "Winter"})
	}, eemiData)

	request, err := http.NewRequest("GET", "/winterfell", nil)
	if err != nil {
		t.Fatal("Creating 'GET " + "" + " request failed!")
	}

	// define response
	respRecorder := httptest.NewRecorder()

	// create app router
	appRouter := mux.NewRouter().StrictSlash(true)
	appRouter.Handle("/winterfell", handler)
	wrapAppRouter(appRouter).ServeHTTP(respRecorder, request)

	eemi, _ := getErrorFromBody(respRecorder.Result())
	assert.Equal(t, 512, respRecorder.Code, "Status code from eemi error should be returned")
	assert.Equal(t, "example of a template error message: Stark", eemi.Message, "Message template should be applied")
	assert.Equal(t, "also needs to cater for this Stark season: Winter", eemi.ResponseAction, "Config Action template should be applied")

}

// Test that errors that are not eemi are caught and returned as an internal server error in the eemi format
func Test_EEMIHandler_Returns_EEMI_Response_For_Unhandled_Error(t *testing.T) {
	currentFolder, _ := os.Getwd()
	eemiData, _ := LoadEemiFromFile(currentFolder+"/testData/eemi_test_data.json", utils.FileSystem{})

	handler := NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("generic error should not surface")
	}, eemiData)

	request, err := http.NewRequest("GET", "/winterfell", nil)
	if err != nil {
		t.Fatal("Creating 'GET " + "" + " request failed!")
	}

	// define response
	respRecorder := httptest.NewRecorder()

	// create app router
	appRouter := mux.NewRouter().StrictSlash(true)
	appRouter.Handle("/winterfell", handler)
	wrapAppRouter(appRouter).ServeHTTP(respRecorder, request)

	eemi, _ := getErrorFromBody(respRecorder.Result())
	assert.Equal(t, 500, respRecorder.Code, "Status code from eemi error should be returned")
	assert.Equal(t, "An internal error has occurred", eemi.Message, "Message template should be applied")
	assert.Equal(t, "Check the logs and retry the operation. If the issue persists contact Dell support", eemi.ResponseAction, "Config Action template should be applied")

}

func wrapAppRouter(appRouter *mux.Router) http.Handler {
	// create new i18n bundle
	i18nBundle := i18n.NewBundle(language.English)

	// navigate to top level locale folder
	currentFolder, dirErr := os.Getwd()
	if dirErr != nil {
		logger.Fatal(dirErr)
	}
	// error messages
	i18nBundle.MustLoadMessageFile(currentFolder + "/testData/eemi_test_locale.en.json")

	// create new language matcher
	languageMatcher := language.NewMatcher([]language.Tag{
		language.English, language.German,
	})

	return middleware.LocaleSelectionMiddleware(i18nBundle, languageMatcher, log.NewLogger("", ""))(appRouter)
}

func getErrorFromBody(r *http.Response) (eemi Config, err error) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &eemi)
	if err != nil {
		return
	}
	return
}
