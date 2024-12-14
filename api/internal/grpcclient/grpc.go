package grpcclient

import (
	"fmt"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/config"
	"github.com/Kudryavkaz/sztuea-api/internal/grpcclient/protos"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var Client protos.CrawlerClient

func init() {
	keepaliveParams := keepalive.ClientParameters{
		Time:                1 * time.Hour,
		Timeout:             30 * time.Second,
		PermitWithoutStream: false,
	}

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepaliveParams),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", config.Config.GetString("grpc_client.address"), config.Config.GetString("grpc_client.port")), opts...)
	if err != nil {
		log.Logger().Error("[grpcclient] grpc.Dial", zap.Error(err))
	}
	Client = protos.NewCrawlerClient(conn)
}
