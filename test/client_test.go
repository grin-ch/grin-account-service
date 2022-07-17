package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grin-ch/grin-api/api/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	conn, err := grpc.Dial("192.168.1.102:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//创建endpoint，并指明grpc调用的接口名和方法名
	client := account.NewUserServiceClient(conn)
	rsp, err := client.SignIn(ctx, &account.SignInReq{
		Contact:  "15200008888",
		Password: "rootroot",
		Key:      "f0bbc28b02c3a17c725b304fa351c312",
		Value:    "bdsc",
	})

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", rsp)
}
