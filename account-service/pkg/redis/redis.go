package redis

import (
	"2margin.vn/account-service/config"
	"2margin.vn/account-service/pkg/logger"
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gookit/goutil/errorx"
	"google.golang.org/protobuf/proto"
)

type redis struct {
	config *config.Config
	log    logger.Logger
	client dapr.Client
}

func NewRedisClient(config *config.Config) (*redis, error) {
	var log = logger.NewAppLogger(config)
	var client, err = dapr.NewClient()

	if err != nil {
		log.Fatal(errorx.Wrap(err, "NewRedisClient dapr.NewClient()"))
	}

	return &redis{config: config, log: log, client: client}, nil
}

func (s *redis) Get(ctx context.Context, key string) (out *proto.Message, wrapError error) {

	state, err := s.client.GetState(ctx, s.config.DaprComponents.StateStore.ComponentName, key, nil)

	if err != nil {
		wrapError = errorx.Wrap(err, "RedisClient.Get.GetState")
		s.log.Error(wrapError)
		return nil, wrapError
	}

	if state.Value == nil {
		wrapError = errorx.Wrap(err, "RedisClient.Get.StateValue.Nil")
		s.log.Error(wrapError)
		return nil, wrapError
	}

	if err := unmarshal(state.Value, out); err != nil {
		wrapError = errorx.Wrap(err, "RedisClient.Get.unmarshal")
		s.log.Error(wrapError)
	}

	return

}

func unmarshal(input []byte, out *proto.Message) error {
	return proto.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(input, *out)
}
