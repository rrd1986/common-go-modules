package assets

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	pkgerr "github.com/pkg/errors"
	ngciErrors "github.com/rrd1986/common-go-modules/errors"
	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/serviceClient"
)

const assetsServiceProtocolEnvVar = "ASSETS_SVC_PROTOCOL"
const assetsServiceHostNameEnvVar = "ASSETS_SVC_HOSTNAME"
const assetsServicePortEnvVar = "ASSETS_SVC_PORT"
const assetsServiceDefaultProtocol = "http"
const assetsServiceDefaultHostname = "system-assets-service.ngci.svc.cluster.local"
const assetsServiceDefaultPort = ""
const assetsServiceBaseUrl = "/api/assets"

// error codes
const AssetNotFound = 100

/**********************************************************************************************************
*
*		System Assets client
*
**********************************************************************************************************/

type ClientType interface {
	GetAssets(ctx context.Context, params *map[string]string, assetType string) ([]Asset, error)
	GetAssetsWithPageCount(ctx context.Context, params *map[string]string, assetType string) (*Assets, error)
	GetAssetById(ctx context.Context, assetType string, id string) (Asset, error)
	PatchAssetById(ctx context.Context, assetType string, assetId string, payload interface{}) error
}

type client struct {
	serviceClient.Type
	logger log.LoggerType
}

func NewClient(logger log.LoggerType) ClientType {
	ci := serviceClient.ClientInput{
		DefaultProtocol: assetsServiceDefaultProtocol,
		ProtocolEnvVar:  assetsServiceProtocolEnvVar,
		DefaultHostName: assetsServiceDefaultHostname,
		HostNameEnvVar:  assetsServiceHostNameEnvVar,
		DefaultPort:     assetsServiceDefaultPort,
		PortEnvVar:      assetsServicePortEnvVar,
		BaseUrl:         assetsServiceBaseUrl,
	}
	c := serviceClient.New(&ci)

	logger = logger.WithCustomFields(map[string]interface{}{"component": "assets-client"})

	return client{Type: c, logger: logger}
}

func (c client) GetAssets(ctx context.Context, params *map[string]string, assetType string) ([]Asset, error) {

	tempResponse, e := c.GetAssetsWithPageCount(ctx, params, assetType)
	if e == nil {
		return tempResponse.Assets, nil
	}
	return nil, e

}

func (c client) GetAssetsWithPageCount(ctx context.Context, params *map[string]string, assetType string) (*Assets, error) {
	c.logger.Debugf("Sending get assets type %s request to downstream service", assetType)
	response, e := c.Get(ctx, params, serviceClient.GetAcceptJsonHeader(), serviceClient.AppendToUri(c.GetApiPath(), assetType))
	if response != nil && response.StatusCode() == http.StatusNotFound {
		c.logger.Errorf("Asset type %s not found", assetType)
		return nil, ngciErrors.NewErrorStr(fmt.Sprintf("Asset types %s not found", assetType), AssetNotFound)
	}
	if e != nil || response == nil || response.StatusCode() != http.StatusOK {
		if e == nil {
			e = errors.New("unexpected response")
		}
		c.logger.Errorf("Error retrieving asset type %s: %s", assetType, e)
		return nil, pkgerr.Wrapf(e, "Error retrieving asset type %s", assetType)
	} else {
		c.logger.Debugf("Get asset type %s request was successful", assetType)
	}

	assetsResponse, e := convertBytesToAssets(response.Body())
	if e != nil {
		return nil, pkgerr.Wrap(e, "Error converting assets response")
	}
	return &assetsResponse, e
}

func (c client) GetAssetById(ctx context.Context, assetType string, id string) (Asset, error) {
	c.logger.Debugf("Sending get asset %s, id %s request to downstream service", assetType, id)
	response, e := c.Get(ctx, nil, serviceClient.GetAcceptJsonHeader(), serviceClient.AppendToUri(c.GetApiPath(), assetType, id))
	if response != nil && response.StatusCode() == http.StatusNotFound {
		c.logger.Errorf("Asset type %s, id %s not found", assetType, id)
		return nil, ngciErrors.NewErrorStr(fmt.Sprintf("Asset %s with id %s not found", assetType, id), AssetNotFound)
	}
	if e != nil || response == nil || response.StatusCode() != http.StatusOK {
		if e == nil {
			e = errors.New("unexpected response")
		}
		c.logger.Errorf("Error retrieving asset type %s, id %s: %s", assetType, id, e)
		return nil, pkgerr.Wrapf(e, "Error retrieving asset type %s id %s", assetType, id)
	} else {
		c.logger.Debugf("Get asset type %s id %s request was successful", assetType, id)
	}

	assetResponse, e := convertBytesToAsset(response.Body())
	if e != nil {
		return nil, pkgerr.Wrap(e, "Error converting asset response")
	}
	return assetResponse, e
}

func (c client) PatchAssetById(ctx context.Context, assetType string, assetId string, payload interface{}) error {
	c.logger.Debugf("updating asset property of assetType %s, assetid %s request to downstream service", assetType, assetId)
	uri := serviceClient.AppendToUri(c.GetApiPath(), assetType, assetId)
	response, e := c.Patch(ctx, payload, serviceClient.GetAcceptJsonHeader(), uri)
	if response != nil && response.StatusCode() == http.StatusNotFound {
		c.logger.Errorf("Fail to updating asset property of assetType %s, id %s not found", assetType, assetId)
		return ngciErrors.NewErrorStr(fmt.Sprintf("Fail to updating asset property of assetType %s with id %s not found", assetType, assetId), AssetNotFound)
	}
	if e != nil || response == nil || response.StatusCode() != http.StatusOK {
		if e == nil {
			e = errors.New("unexpected response")
		}
		c.logger.Errorf("Error while Updating the asset property asset type %s, id %s: %s", assetType, assetId, e)
		return pkgerr.Wrapf(e, "Error while Updating the asset property asset type %s id %s", assetType, assetId)
	}
	c.logger.Debugf("Successfully updated asset property asset type %s id %s request was successful", assetType, assetId)
	return nil
}
