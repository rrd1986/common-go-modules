package assets

import (
	"context"
	"net/http"

	"github.com/rrd1986/common-go-modules/serviceClient"
	"github.com/stretchr/testify/mock"
)

/**********************************************************************************************************
*
*		Mock Asset client implements Client Type
*
**********************************************************************************************************/

type MockClient struct {
	mock.Mock
}

func (c *MockClient) GetAssets(ctx context.Context, params *map[string]string, assetType string) (assets []Asset, err error) {
	args := c.Called(params, assetType)
	assets, _ = args.Get(0).([]Asset)
	err, _ = args.Get(1).(error)
	return
}

func (c *MockClient) GetAssetsWithPageCount(ctx context.Context, params *map[string]string, assetType string) (assets *Assets, err error) {
	args := c.Called(params, assetType)
	assets, _ = args.Get(0).(*Assets)
	err, _ = args.Get(1).(error)
	return
}

func (c *MockClient) GetAssetById(ctx context.Context, assetType string, id string) (asset Asset, err error) {
	args := c.Called(assetType, id)
	asset, _ = args.Get(0).(Asset)
	err, _ = args.Get(1).(error)
	return
}

func (c *MockClient) PatchAssetById(ctx context.Context, assetType string, id string, payload interface{}) (err error) {
	args := c.Called(assetType, id, payload)
	err, _ = args.Get(0).(error)
	return
}

/**********************************************************************************************************
*
*		Test data
*
**********************************************************************************************************/

func getDummyAssetsBytes() []byte {
	return []byte(`{"pages":"1","count":"1","assets":[{"InitialDeployment":{"CallHomeConfigured":false,"InitialDeploymentCompleted":false,"ManagementNetworkConfigured":false,"NetworkServicesConfigured":false,"PlatformControllerConfigured":false,"ProductionNetworkConfigured":false,"RbacConfigured":false},"assetType":"CAM","lock":"","model":"CAM","objectName":"CAMAsset","quarantine":4,"standardDefaultsVersion":"2.0","uniqueIdentifier":"1"}]}`)
}

func getDummyAssetBytes() []byte {
	return []byte(`{"InitialDeployment":{"CallHomeConfigured":false,"InitialDeploymentCompleted":false,"ManagementNetworkConfigured":false,"NetworkServicesConfigured":false,"PlatformControllerConfigured":false,"ProductionNetworkConfigured":false,"RbacConfigured":false},"assetType":"CAM","lock":"","model":"CAM","objectName":"CAMAsset","quarantine":4,"standardDefaultsVersion":"2.0","uniqueIdentifier":"1"}`)
}

func getDummyAssetsType() []Asset {
	return []Asset{{"assetType": "CAM", "lock": "", "model": "CAM", "objectName": "CAMAsset", "quarantine": 4.0, "standardDefaultsVersion": "2.0", "uniqueIdentifier": "1", "InitialDeployment": map[string]interface{}{"InitialDeploymentCompleted": false, "ManagementNetworkConfigured": false, "NetworkServicesConfigured": false, "PlatformControllerConfigured": false, "ProductionNetworkConfigured": false, "RbacConfigured": false, "CallHomeConfigured": false}}}
}

func getDummyAssetType() Asset {
	return Asset{"assetType": "CAM", "lock": "", "model": "CAM", "objectName": "CAMAsset", "quarantine": 4.0, "standardDefaultsVersion": "2.0", "uniqueIdentifier": "1", "InitialDeployment": map[string]interface{}{"InitialDeploymentCompleted": false, "ManagementNetworkConfigured": false, "NetworkServicesConfigured": false, "PlatformControllerConfigured": false, "ProductionNetworkConfigured": false, "RbacConfigured": false, "CallHomeConfigured": false}}
}

func getNotFoundBytes() []byte {
	return []byte(`"{message: Resource not found}"`)
}

func getAssetsSuccessResponse() serviceClient.HeaderedResponse {
	return serviceClient.MockClientResponse{Code: http.StatusOK, Bytes: getDummyAssetsBytes(), Head: nil}
}

func getAssetSuccessResponse() serviceClient.HeaderedResponse {
	return serviceClient.MockClientResponse{Code: http.StatusOK, Bytes: getDummyAssetBytes(), Head: nil}
}

func getBadJsonResponse() serviceClient.HeaderedResponse {
	return serviceClient.MockClientResponse{Code: http.StatusOK, Bytes: []byte(`this is not json`), Head: nil}
}

func getNotFoundResponse() serviceClient.HeaderedResponse {
	return serviceClient.MockClientResponse{Code: http.StatusNotFound, Bytes: getNotFoundBytes(), Head: nil}
}

func getErrorResponse() serviceClient.HeaderedResponse {
	return serviceClient.MockClientResponse{Code: http.StatusInternalServerError, Bytes: nil, Head: nil}
}

func getDummyAssetPatch() interface{} {
	var qState struct {
		QuarantineState       int
		QuarantineDescription string
	}
	qState.QuarantineState = 5
	qState.QuarantineDescription = "Invalid configuration"
	return qState
}
func getDummyParams() *map[string]string {
	params := map[string]string{"test1": "val1", "test2": "val2"}
	return &params
}

func getDummyEmptyParams() *map[string]string {
	return nil
}

func getJsonHeader() *map[string]string {
	return serviceClient.GetAcceptJsonHeader()
}
