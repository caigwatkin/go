package mock

import (
	"context"

	"github.com/caigwatkin/go/log"
)

var (
	Client log.Client = client{}
)

type client struct{}

func (c client) Debug(_ context.Context, _ string, _ ...log.Field) {}

func (c client) Info(_ context.Context, _ string, _ ...log.Field) {}

func (c client) Notice(_ context.Context, _ string, _ ...log.Field) {}

func (c client) Warn(_ context.Context, _ string, _ ...log.Field) {}

func (c client) Error(_ context.Context, _ string, _ ...log.Field) {}

func (c client) Fatal(_ context.Context, _ string, _ ...log.Field) {}
