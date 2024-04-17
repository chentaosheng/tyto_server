package tyto

import "tyto/core/logging"

type Context interface {
	Logger() logging.Logger
}
