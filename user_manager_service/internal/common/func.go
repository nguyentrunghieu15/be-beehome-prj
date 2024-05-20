package common

import (
	"context"
	"fmt"

	"google.golang.org/grpc/peer"
)

func StandardMsgError(ctx context.Context, req string, err error) string {
	return fmt.Sprintf("from IP:%v invoke request:%v by error: %v", GetClientIp(ctx), req, err)
}

func StandardMsgInfor(ctx context.Context, req string, info string) string {
	return fmt.Sprintf("from IP :%v invoke request:%v by error: %v", GetClientIp(ctx), req, info)
}

func StandardMsgWarn(ctx context.Context, req string, warn string) string {
	return fmt.Sprintf("from IP :%v invoke request:%v by error: %v", GetClientIp(ctx), req, warn)
}

func GetClientIp(ctx context.Context) string {
	p, _ := peer.FromContext(ctx)
	return p.Addr.String()
}
