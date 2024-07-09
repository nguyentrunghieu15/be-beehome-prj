package addressclient

import (
	"context"

	addressapi "github.com/nguyentrunghieu15/be-beehome-prj/api/address-api"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/logwrapper"
	"google.golang.org/grpc"
)

type IAddressClient interface {
	CheckExistAddress(
		ctx context.Context,
		in *addressapi.CheckExistAddressRequest,
		opts ...grpc.CallOption,
	) (bool, error)
}

type AddressClientWrapper struct {
	client addressapi.AddressServiceClient
	conn   *grpc.ClientConn
	logger logwrapper.ILoggerWrapper
}

func NewAddressClientWrapper(conn *grpc.ClientConn, logger logwrapper.ILoggerWrapper) *AddressClientWrapper {
	return &AddressClientWrapper{
		client: addressapi.NewAddressServiceClient(conn),
		conn:   conn,
		logger: logger,
	}
}

func (client *AddressClientWrapper) CheckExistAddress(
	ctx context.Context,
	in *addressapi.CheckExistAddressRequest,
	opts ...grpc.CallOption,
) (bool, error) {
	result, err := client.client.CheckExistAddress(ctx, in, opts...)
	if err != nil {
		return false, err
	}
	if result != nil && result.Address != nil {
		return true, nil
	}
	return false, err
}
