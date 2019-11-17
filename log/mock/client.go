package mock

import (
	"context"

	"github.com/caigwatkin/go/log"
)

var (
	MockClientSuccess log.Client = mockClientSuccess{}
)

type mockClientSuccess struct{}

func (c mockClientSuccess) Debug(_ context.Context, _ string, _ ...log.Field) {}

func (c mockClientSuccess) Info(_ context.Context, _ string, _ ...log.Field) {}

func (c mockClientSuccess) Notice(_ context.Context, _ string, _ ...log.Field) {}

func (c mockClientSuccess) Warn(_ context.Context, _ string, _ ...log.Field) {}

func (c mockClientSuccess) Error(_ context.Context, _ string, _ ...log.Field) {}

func (c mockClientSuccess) Fatal(_ context.Context, _ string, _ ...log.Field) {}
