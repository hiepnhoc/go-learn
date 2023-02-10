package main

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/pkg/convert"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	commonv1pb "github.com/dapr/go-sdk/dapr/proto/common/v1"
	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"time"
)

func main() {
	HandlerInvoke()
}

const (
	address = "localhost:5000"
)

func HandlerInvoke() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewAppCallbackClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "account-service", "server")

	registerRequest := &userService.RegisterRequest{
		Role:      "admin",
		Password:  "123456",
		FirstName: "hiepln",
		LastName:  "hiepln1",
		Email:     "hiepln" + string(rand.Intn(1000)) + "@acbs.com.vn",
	}

	registerRequestBytes, err := convert.ProtoToJsonBytes(registerRequest)

	log.Println(string(registerRequestBytes))

	input := commonv1pb.InvokeRequest{
		ContentType: "application/json",
		Method:      "Register",
		Data: &any.Any{
			Value: registerRequestBytes,
		},
	}

	r, err := c.OnInvoke(ctx, &input)

	if err != nil {
		panic(err)
	}

	registerResponse := userService.RegisterResponse{}

	err = convert.ProtoMessage(r.Data.Value, &registerResponse)
	if err != nil {
		panic(err)
	}

	log.Printf("Greeting: %s", registerResponse.String())
}
