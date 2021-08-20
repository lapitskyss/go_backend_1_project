package repository

import "context"

type Store interface {
	Connect(ctx context.Context) error
	GetConnection() interface{}
	CloseConnection()

	Link() LinkInterface
	RedirectLog() RedirectLogInterface
}
