package ports

import "context"

type TxFunc func(ctx context.Context) error

type TxManager interface {
	Run(ctx context.Context, fn TxFunc) error
}
