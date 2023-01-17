package standards

import (
	"context"

	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/serviceClient"
)

// System Standards
const standardsServiceProtocolEnvVar = "STANDARDS_SVC_PROTOCOL"
const standardsServiceHostNameEnvVar = "STANDARDS_SVC_HOSTNAME"
const standardsServicePortEnvVar = "STANDARDS_SVC_PORT"
const standardsServiceDefaultProtocol = "http"
const standardsServiceDefaultHostname = "system-standards-service.ngci.svc.cluster.local"
const standardsServiceDefaultPort = ""
const standardsServiceBaseUrl = "/api/standards/defaults"
const latest = "latest"

/**********************************************************************************************************
*
*		System Standards client
*
**********************************************************************************************************/

type ClientType interface {
	GetLatest(ctx context.Context, assetType string) (serviceClient.BasicResponse, error)
}

type Client struct {
	serviceClient.Type
	logger log.LoggerType
}

func NewClient(logger log.LoggerType) ClientType {
	ci := serviceClient.ClientInput{
		DefaultProtocol: standardsServiceDefaultProtocol,
		ProtocolEnvVar:  standardsServiceProtocolEnvVar,
		DefaultHostName: standardsServiceDefaultHostname,
		HostNameEnvVar:  standardsServiceHostNameEnvVar,
		DefaultPort:     standardsServiceDefaultPort,
		PortEnvVar:      standardsServicePortEnvVar,
		BaseUrl:         standardsServiceBaseUrl,
	}
	client := serviceClient.New(&ci)

	logger = logger.WithCustomFields(map[string]interface{}{"component": "standards-client"})

	return Client{Type: client, logger: logger}
}

func (c Client) GetLatest(ctx context.Context, assetType string) (serviceClient.BasicResponse, error) {
	c.logger.Debugf("Sending get latest %s standards request to downstream service", assetType)
	uri := serviceClient.AppendToUri(c.GetApiPath(), assetType, latest)
	response, e := c.Get(ctx, nil, serviceClient.GetAcceptJsonHeader(), uri)
	if e != nil {
		c.logger.Errorf("Error retrieving latest standards %s from downstream service: %v", assetType, e)
	} else {
		c.logger.Debugf("Get latest %s standard request was successful", assetType)
	}
	return response, e
}
