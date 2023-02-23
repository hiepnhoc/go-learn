package darp

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/auth-microservice/config"
	authServerGRPC "github.com/AleksK1NG/auth-microservice/internal/user/delivery/grpc/service"
	"github.com/AleksK1NG/auth-microservice/pkg/convert"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	commonv1pb "github.com/dapr/go-sdk/dapr/proto/common/v1"
	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/proto"
	"strings"
)

const (
	ContentType     = "application/json"
	ContentTypeUTF8 = "application/json; utf-8"
)

// server is our user app
type darpGprc struct {
	pb.UnimplementedAppCallbackServer
	logger   logger.Logger
	cfg      *config.Config
	userGrpc *authServerGRPC.UsersService
}

func NewDarpGprc(logger logger.Logger, cfg *config.Config, userService *authServerGRPC.UsersService) *darpGprc {
	return &darpGprc{logger: logger, cfg: cfg, userGrpc: userService}
}

// EchoMethod is a simple demo method to invoke
func (s *darpGprc) EchoMethod() string {
	return "pong"
}

// EchoMethod is a simple demo method to invoke
func (s *darpGprc) Register(ctx context.Context, in *commonv1pb.InvokeRequest) ([]byte, error) {
	registerRequest := userService.RegisterRequest{}

	//registerRequest := userService.RegisterRequest{
	//	Role:      "admin",
	//	Password:  "123456",
	//	FirstName: "hiepln",
	//	LastName:  "hiepln1",
	//	Email:     "hiepln8@acbs.com.vn",
	//}
	//
	//registerRequestBytes, err := convert.ProtoBytes(&registerRequest)
	//
	//input := commonv1pb.InvokeRequest{
	//	Data: &any.Any{
	//		Value: registerRequestBytes,
	//	},
	//}

	//if err != nil {
	//	return nil, err
	//}

	err := s.convertToMessage(in, &registerRequest)

	if err != nil {
		return nil, err
	}

	res, err := s.userGrpc.Register(ctx, &registerRequest)

	if err != nil {
		return nil, err
	}

	return convert.ProtoBytes(res)
}

// This method gets invoked when a remote service has called the app through Dapr
// The payload carries a Method to identify the method, a set of metadata properties and an optional payload
func (s *darpGprc) OnInvoke(ctx context.Context, in *commonv1pb.InvokeRequest) (*commonv1pb.InvokeResponse, error) {
	var response []byte
	var err error

	switch in.Method {

	case "Register":
		response, err = s.Register(ctx, in)
		if err != nil {
			return nil, err
		}
	}

	return &commonv1pb.InvokeResponse{
		ContentType: in.ContentType,
		Data:        &any.Any{Value: response, TypeUrl: in.Data.TypeUrl},
	}, nil
}

// Dapr will call this method to get the list of topics the app wants to subscribe to. In this example, we are telling Dapr
// To subscribe to a topic named TopicA
func (s *darpGprc) ListTopicSubscriptions(ctx context.Context, in *empty.Empty) (*pb.ListTopicSubscriptionsResponse, error) {
	return &pb.ListTopicSubscriptionsResponse{
		Subscriptions: []*pb.TopicSubscription{
			{Topic: "TopicA"},
		},
	}, nil
}

// Dapr will call this method to get the list of bindings the app will get invoked by. In this example, we are telling Dapr
// To invoke our app with a binding named storage
func (s *darpGprc) ListInputBindings(ctx context.Context, in *empty.Empty) (*pb.ListInputBindingsResponse, error) {
	return &pb.ListInputBindingsResponse{
		Bindings: []string{"storage"},
	}, nil
}

// This method gets invoked every time a new event is fired from a registered binding. The message carries the binding name, a payload and optional metadata
func (s *darpGprc) OnBindingEvent(ctx context.Context, in *pb.BindingEventRequest) (*pb.BindingEventResponse, error) {
	fmt.Println("Invoked from binding")
	return &pb.BindingEventResponse{}, nil
}

// This method is fired whenever a message has been published to a topic that has been subscribed. Dapr sends published messages in a CloudEvents 0.3 envelope.
func (s *darpGprc) OnTopicEvent(ctx context.Context, in *pb.TopicEventRequest) (*pb.TopicEventResponse, error) {
	fmt.Println("Topic message arrived")
	return &pb.TopicEventResponse{}, nil
}

func (s *darpGprc) convertToMessage(input *commonv1pb.InvokeRequest, out proto.Message) error {

	if strings.HasPrefix(input.ContentType, ContentType) {
		return convert.ProtoJsonToMessage(input.Data.Value, out)
	}
	return convert.ProtoMessage(input.Data.Value, out)

}
