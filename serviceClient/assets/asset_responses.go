package assets

import (
	"encoding/json"
)

/**********************************************************************************************************
*
*		Responses returned from the system assets service
*
**********************************************************************************************************/
type Assets struct {
	Pages  string  `json:"pages"`
	Count  string  `json:"count"`
	Assets []Asset `json:"assets"`
}

type Asset map[string]interface{}

/**********************************************************************************************************
*
*		Conversion functions
*
**********************************************************************************************************/

func convertBytesToAssets(data []byte) (assets Assets, err error) {
	err = json.Unmarshal(data, &assets)
	return
}

func convertBytesToAsset(data []byte) (asset Asset, err error) {
	err = json.Unmarshal(data, &asset)
	return
}
