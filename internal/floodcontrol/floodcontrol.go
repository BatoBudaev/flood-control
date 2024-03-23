package floodcontrol

import "context"

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}
