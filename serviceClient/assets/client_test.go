package assets

import (
	"context"
	"errors"
	"testing"

	ngciErrors "github.com/rrd1986/common-go-modules/errors"
	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/serviceClient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Setup test called before each test in the suite
func (suite *AssetTestSuite) SetupTest() {
	suite.mockClient = new(MockClient)
	suite.mockDbClient = new(serviceClient.MockClient)
}

const dummyAssetType = "Type1"
const dummyAssetId = "123"
const dummyApiPath = "DummyApiPath"
const dummyGetAssetsUri = dummyApiPath + "/" + dummyAssetType
const dummyGetAssetUri = dummyGetAssetsUri + "/" + dummyAssetId

type AssetTestSuite struct {
	suite.Suite
	mockClient   *MockClient
	mockDbClient *serviceClient.MockClient
}

func TestAssetServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AssetTestSuite))
}

var logger = log.NewLogger("", "")

/**********************************************************************************************************
*
*		Client tests
*
**********************************************************************************************************/

func (suite *AssetTestSuite) TestClient_GetAssets_Success() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyParams(), getJsonHeader(), dummyGetAssetsUri).Return(getAssetsSuccessResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssets(context.Background(), getDummyParams(), dummyAssetType)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err, "Unexpected error encountered")
	assert.Equal(suite.T(), getDummyAssetsType(), assetList, "Unexpected response")
}

func (suite *AssetTestSuite) TestClient_GetAssets_BadJson() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyParams(), getJsonHeader(), dummyGetAssetsUri).Return(getBadJsonResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssets(context.Background(), getDummyParams(), dummyAssetType)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
}

func (suite *AssetTestSuite) TestClient_GetAssets_NotFound() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyParams(), getJsonHeader(), dummyGetAssetsUri).Return(getNotFoundResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssets(context.Background(), getDummyParams(), dummyAssetType)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
	ngciErr, ok := err.(ngciErrors.Error)
	assert.True(suite.T(), ok, "Unexpected error type")
	assert.Equal(suite.T(), AssetNotFound, int(ngciErr.ErrorCode), "Unexpected response")
}

func (suite *AssetTestSuite) TestClient_GetAssets_ErrorReturned() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyParams(), getJsonHeader(), dummyGetAssetsUri).Return(getErrorResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssets(context.Background(), getDummyParams(), dummyAssetType)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
}

func (suite *AssetTestSuite) TestClient_GetAssets_ErrorEncountered() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyParams(), getJsonHeader(), dummyGetAssetsUri).Return(nil, errors.New("no assets service"))
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssets(context.Background(), getDummyParams(), dummyAssetType)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
}

func (suite *AssetTestSuite) TestClient_GetAsset_Success() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyEmptyParams(), getJsonHeader(), dummyGetAssetUri).Return(getAssetSuccessResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	asset, err := service.GetAssetById(context.Background(), dummyAssetType, dummyAssetId)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err, "Unexpected error encountered")
	assert.Equal(suite.T(), getDummyAssetType(), asset, "Unexpected response")
}

func (suite *AssetTestSuite) TestClient_GetAsset_NotFound() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyEmptyParams(), getJsonHeader(), dummyGetAssetUri).Return(getNotFoundResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssetById(context.Background(), dummyAssetType, dummyAssetId)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
	ngciErr, ok := err.(ngciErrors.Error)
	assert.True(suite.T(), ok, "Unexpected error type")
	assert.Equal(suite.T(), AssetNotFound, int(ngciErr.ErrorCode), "Unexpected response")
}

func (suite *AssetTestSuite) TestClient_GetAsset_BadJson() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyEmptyParams(), getJsonHeader(), dummyGetAssetUri).Return(getBadJsonResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssetById(context.Background(), dummyAssetType, dummyAssetId)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
}

func (suite *AssetTestSuite) TestClient_GetAsset_ErrorReturned() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyEmptyParams(), getJsonHeader(), dummyGetAssetUri).Return(getErrorResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssetById(context.Background(), dummyAssetType, dummyAssetId)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
}

func (suite *AssetTestSuite) TestClient_GetAsset_ErrorEncountered() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Get", getDummyEmptyParams(), getJsonHeader(), dummyGetAssetUri).Return(nil, errors.New("no assets service"))
	service := client{suite.mockDbClient, logger}

	// Act
	assetList, err := service.GetAssetById(context.Background(), dummyAssetType, dummyAssetId)

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), assetList, "Unexpected response encountered")
	assert.NotNil(suite.T(), err, "Expected error not encountered")
}

func (suite *AssetTestSuite) TestClient_PatchAssetById_Success() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Patch", getDummyAssetPatch(), getJsonHeader(), dummyGetAssetUri).Return(getAssetSuccessResponse(), nil)
	service := client{suite.mockDbClient, logger}
	// Act
	err := service.PatchAssetById(context.Background(), dummyAssetType, dummyAssetId, getDummyAssetPatch())
	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err, "Unexpected error encountered")
}

func (suite *AssetTestSuite) TestClient_Patch_NotFound() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Patch", getDummyAssetPatch(), getJsonHeader(), dummyGetAssetUri).Return(getNotFoundResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	err := service.PatchAssetById(context.Background(), dummyAssetType, dummyAssetId, getDummyAssetPatch())

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.NotNil(suite.T(), err, "Expected error not encountered")
	ngciErr, ok := err.(ngciErrors.Error)
	assert.True(suite.T(), ok, "Unexpected error type")
	assert.Equal(suite.T(), AssetNotFound, int(ngciErr.ErrorCode), "Unexpected response")
}

func (suite *AssetTestSuite) TestClient_Patch_ErrorReturned() {
	// Arrange
	suite.mockDbClient.On("GetApiPath").Return(dummyApiPath)
	suite.mockDbClient.On("Patch", getDummyAssetPatch(), getJsonHeader(), dummyGetAssetUri).Return(getErrorResponse(), nil)
	service := client{suite.mockDbClient, logger}

	// Act
	err := service.PatchAssetById(context.Background(), dummyAssetType, dummyAssetId, getDummyAssetPatch())

	// Assert
	suite.mockDbClient.AssertExpectations(suite.T())
	assert.NotNil(suite.T(), err, "Expected error not encountered")
}
