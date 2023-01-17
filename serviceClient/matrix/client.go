package matrix

import (
	"context"
	"strconv"
	"strings"

	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/serviceClient"
)

const matrixServiceProtocolEnvVar = "MATRIX_SVC_PROTOCOL"
const matrixServiceHostNameEnvVar = "MATRIX_SVC_HOSTNAME"
const matrixServicePortEnvVar = "MATRIX_SVC_PORT"
const matrixServiceDefaultProtocol = "http"
const matrixServiceDefaultHostname = "software-matrix-service.ngci.svc.cluster.local"
const matrixServiceDefaultPort = ""
const matrixServiceBaseUrl = "/api/matrix"
const matrixServiceVersionsEndpoint = "versions"

type AvailableVersionResponse struct {
	Results []AvailableVersion `json:"results"`
}

type InstallMatrix struct {
	AvailableVersion
	Components interface{} `json:"install_matrix"`
}

type AvailableVersion struct {
	Name        string `json:"name"`
	Id          int    `json:"id"`
	Description string `json:"description"`
}

/**********************************************************************************************************
*
*		Software Matrix client
*
**********************************************************************************************************/

type ClientType interface {
	GetVersions(ctx context.Context, assetType string) (serviceClient.BasicResponse, error)
	GetVersionsById(ctx context.Context, assetType string, id int) (serviceClient.BasicResponse, error)
}

type Client struct {
	serviceClient.Type
	logger log.LoggerType
}

func NewClient(logger log.LoggerType) ClientType {
	ci := serviceClient.ClientInput{
		DefaultProtocol: matrixServiceDefaultProtocol,
		ProtocolEnvVar:  matrixServiceProtocolEnvVar,
		DefaultHostName: matrixServiceDefaultHostname,
		HostNameEnvVar:  matrixServiceHostNameEnvVar,
		DefaultPort:     matrixServiceDefaultPort,
		PortEnvVar:      matrixServicePortEnvVar,
		BaseUrl:         matrixServiceBaseUrl,
	}

	logger = logger.WithCustomFields(map[string]interface{}{"component": "matrix-client"})

	client := serviceClient.New(&ci)
	return Client{Type: client, logger: logger}
}

func (c Client) GetVersions(ctx context.Context, assetType string) (serviceClient.BasicResponse, error) {
	c.logger.Debugf("Sending get versions by assets type %s request to downstream service", assetType)
	response, e := c.Get(ctx, nil, serviceClient.GetAcceptJsonHeader(), c.getVersionsUri(assetType))
	if e != nil {
		c.logger.Errorf("Error retrieving versions by assets type %s: %s", assetType, e)
	} else {
		c.logger.Debugf("Get versions by assets type %s request was successful", assetType)
	}
	return response, e
}

func (c Client) GetVersionsById(ctx context.Context, assetType string, id int) (serviceClient.BasicResponse, error) {
	c.logger.Debugf("Sending get version by id %v by assets type %s request to downstream service", id, assetType)
	response, e := c.Get(ctx, nil, serviceClient.GetAcceptJsonHeader(), c.getVersionByIdUri(assetType, id))
	if e != nil {
		c.logger.Errorf("Error retrieving versions by id: %v by assets type %s: %v", id, assetType, e)
	} else {
		c.logger.Debugf("Get versions by id: %v by assets type %s request was successful", id, assetType)
	}
	return response, e
}

func (c Client) getVersionsUri(assetType string) string {
	return serviceClient.AppendToUri(c.GetApiPath(), strings.ToLower(assetType), matrixServiceVersionsEndpoint)
}

func (c Client) getVersionByIdUri(assetType string, id int) string {
	return serviceClient.AppendToUri(c.GetApiPath(), strings.ToLower(assetType), matrixServiceVersionsEndpoint, strconv.Itoa(id))
}

func (c Client) getUpgradeUri(assetType string, assetId string, versionId int) string {
	return serviceClient.AppendToUri(c.GetApiPath(), strings.ToLower(assetType), assetId, matrixServiceVersionsEndpoint)
}
