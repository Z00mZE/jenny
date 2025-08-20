package jenny

import (
	"context"

	"github.com/Z00mZE/jenny/pb/service/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type App struct {
}

func NewApp(ctx context.Context) (*App, error) {
	return new(App), nil
}

func (a *App) Send(ctx context.Context, event *types.Event) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (a *App) Stream(g grpc.ClientStreamingServer[types.Event, emptypb.Empty]) error {
	//TODO implement me
	panic("implement me")
}
