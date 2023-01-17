package matrix

import (
	"context"

	"github.com/rrd1986/common-go-modules/serviceClient"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (c MockClient) GetVersions(ctx context.Context, assetType string) (resp serviceClient.BasicResponse, err error) {
	args := c.Called(assetType)

	resp, _ = args.Get(0).(serviceClient.BasicResponse)
	err, _ = args.Get(1).(error)
	return
}

func (c MockClient) GetVersionsById(ctx context.Context, assetType string, id int) (resp serviceClient.BasicResponse, err error) {
	args := c.Called(assetType, id)

	resp, _ = args.Get(0).(serviceClient.BasicResponse)
	err, _ = args.Get(1).(error)
	return
}
