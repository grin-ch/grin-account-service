package cmd

import (
	"context"

	"github.com/grin-ch/grin-account-service/pkg/auth"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func recoveryFunc(p interface{}) (err error) {
	log.Errorf("panic triggered: %v", p)
	return status.Errorf(codes.Unknown, "unknow error")
}

func authFunc(ctx context.Context) (context.Context, error) {
	claims, err := auth.ParseToken(ctx)
	if err != nil {
		return nil, err
	}
	err = claims.Claims.Valid()
	if err != nil {
		return nil, err
	}
	return ctx, nil
}
