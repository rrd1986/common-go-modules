package standards

import (
	"context"

	"github.com/rrd1986/common-go-modules/serviceClient"
	"github.com/stretchr/testify/mock"
)

/**********************************************************************************************************
*
*		Mock system standards svc client methods
*
**********************************************************************************************************/

type MockClient struct {
	mock.Mock
}

func (m MockClient) GetLatest(ctx context.Context, assetType string) (serviceClient.BasicResponse, error) {
	args := m.Called(assetType)

	resp, _ := args.Get(0).(serviceClient.BasicResponse)
	err, _ := args.Get(1).(error)

	return resp, err
}
