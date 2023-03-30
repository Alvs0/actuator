package impl

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jasonlvhit/gocron"
)

func (g *GeneratorService) Stop(ctx context.Context, _ *empty.Empty) (res *empty.Empty, err error) {
	gocron.Clear()

	res = &empty.Empty{}

	return
}
