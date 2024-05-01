package panicutil

import (
	"runtime/debug"
	"tyto/core/tyto"
)

func Recover(ctx tyto.Context) {
	if err := recover(); err != nil {
		switch v := err.(type) {
		case error:
			ctx.Logger().NoCallerError("panic:", v.Error())
		case string:
			ctx.Logger().NoCallerError("panic:", v)
		default:
			ctx.Logger().NoCallerError("panic: unknown error")
		}

		ctx.Logger().NoCallerError(string(debug.Stack()))
	}
}

func RecoverWith(ctx tyto.Context, f func()) {
	if err := recover(); err != nil {
		switch v := err.(type) {
		case error:
			ctx.Logger().NoCallerError("panic:", v.Error())
		case string:
			ctx.Logger().NoCallerError("panic:", v)
		default:
			ctx.Logger().NoCallerError("panic: unknown error")
		}

		ctx.Logger().NoCallerError(string(debug.Stack()))
		f()
	}
}

func PrintStack(ctx tyto.Context) {
	ctx.Logger().Info(string(debug.Stack()))
}
