package interop

import (
	fmt "fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	FlagInteropNodePort      = "interop-node-port"
	ReconnectInterval        = 10 * time.Second
	ErrInteropClientNotReady = "interop client not ready"
)

type InteropClientFactory struct {
	conn *grpc.ClientConn
}

func NewInteropClientFactory(log log.Logger, port uint16) *InteropClientFactory {
	interopClient := &InteropClientFactory{
		conn: nil,
	}

	go func() {
		target := fmt.Sprintf("localhost:%d", port)
		option := grpc.WithTransportCredentials(insecure.NewCredentials())
		for {
			conn, err := grpc.Dial(target, option)
			if err != nil {
				log.Error("failed to connect to interop node", "err", err)
				time.Sleep(ReconnectInterval)
				continue
			}

			interopClient.conn = conn
			return
		}
	}()

	return interopClient
}

func (w *InteropClientFactory) GetInteropClient() InteropClient {
	if w.conn == nil {
		return nil
	}

	return NewInteropClient(w.conn)
}
