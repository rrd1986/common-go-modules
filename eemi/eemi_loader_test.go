package eemi

import (
	"os"
	"testing"

	"github.com/rrd1986/common-go-modules/utils"
	"github.com/stretchr/testify/assert"
)

func Test_LoadEemiFromFile_Returns_Error_Data(t *testing.T) {
	currentFolder, _ := os.Getwd()

	result, err := LoadEemiFromFile(currentFolder+"/testData/eemi_test_data.json", utils.FileSystem{})

	assert.Nil(t, err, "Error parsing eemi errors")
	assert.Equal(t, 3, len(result))
	assert.Equal(t, 512, result["NGCI0002"].Status)
}
