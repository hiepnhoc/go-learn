package darp

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/auth-microservice/config"
	authServerGRPC "github.com/AleksK1NG/auth-microservice/internal/user/delivery/grpc/service"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"sync"
)

var (
	daprClient dapr.Client
	doOnce     sync.Once
)

type darpHttp struct {
	logger   logger.Logger
	cfg      *config.Config
	userGrpc *authServerGRPC.UsersService
}

func NewDarpHttp(logger logger.Logger, cfg *config.Config, userService *authServerGRPC.UsersService) *darpHttp {
	return &darpHttp{logger: logger, cfg: cfg, userGrpc: userService}
}

func (s *darpHttp) AddInvocationHandler(log logger.Logger, service common.Service, invocationMethod string) {

	err := service.AddServiceInvocationHandler(invocationMethod, s.echoHandler)
	if err != nil {
		log.Fatalf("addInvocationHandler: error adding invocation handler: %s %v", invocationMethod, err)
	}
	log.Infof("addInvocationHandler: add the %s method success.", invocationMethod)

}

func (s *darpHttp) echoHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	//log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
	// do something with the invocation here

	//registerRequest := userService.RegisterRequest{}

	//err1 := convert.ProtoJsonToMessage(in.Data, &registerRequest)
	//
	//if err1 != nil {
	//	return
	//}

	fmt.Println("test log")

	registerRequest := userService.RegisterRequest{
		Role:      "admin",
		Password:  "123456",
		FirstName: "hiepln",
		LastName:  "hiepln1",
		Email:     "hiepln800@acbs.com.vn",
	}

	res, err := s.userGrpc.Register(ctx, &registerRequest)

	fmt.Println(res)
	//resByte, err := convert.ProtoBytes(res)

	out = &common.Content{
		Data:        in.Data,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}
