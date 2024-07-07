package cerbosx

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/cerbos/cerbos-sdk-go/cerbos"
)

const (
	READ   = "read"
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

type CerbosClientConfig struct {
	CerbosAddress string
}

type CerbosClientWrapper struct {
	config *CerbosClientConfig
	cerbos *cerbos.GRPCClient
	once   sync.Once
}

var DefaultClient *CerbosClientWrapper

func NewCerbosClientWrapperWithConfig(config *CerbosClientConfig) *CerbosClientWrapper {
	return &CerbosClientWrapper{
		config: config,
	}
}

func (inst *CerbosClientWrapper) connect() error {
	if inst.config == nil {
		return errors.New("missing config")
	}
	c, err := cerbos.New(inst.config.CerbosAddress, cerbos.WithPlaintext())
	if err != nil {
		return err
	}
	inst.cerbos = c
	return nil
}

func (inst *CerbosClientWrapper) doConnect() {
	for {
		err := inst.connect()
		log.Println(err)
		if err != nil {
			time.Sleep(300)
			continue
		}
		break
	}
}

func (inst *CerbosClientWrapper) Setup() {
	inst.once.Do(inst.doConnect)
}

func (inst *CerbosClientWrapper) CanActive(
	ctx context.Context,
	principal *cerbos.Principal,
	resource *cerbos.Resource,
	action string,
) (bool, error) {
	return inst.cerbos.IsAllowed(ctx, principal, resource, action)
}
