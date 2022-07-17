package cmd

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
)

func etcdConncter(builder resolver.Builder) connecter {
	return func(s string) (*grpc.ClientConn, error) {
		addr := fmt.Sprintf("etcd:///%s", s)
		return grpc.Dial(addr,
			grpc.WithResolvers(builder),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, roundrobin.Name)),
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:                10 * time.Second,
					Timeout:             100 * time.Millisecond,
					PermitWithoutStream: true},
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}
