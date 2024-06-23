package exceptUtils

import (
	"code-gen/utils/logger"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Recover() {
	v := recover()
	if v == nil {
		return
	}

	logger.Logger.Error("出现异常", zap.Error(errors.New(fmt.Sprintf("%v", v))))
}
func CatchError(fn func() error) {
	defer Recover()
	err := fn()
	if err != nil {
		logger.Logger.Error("出现异常", zap.Error(errors.WithStack(err)))
	}
}
func CatchErrorWithMessage(errorMsg string, fn func() error) {
	defer Recover()
	err := fn()
	if err != nil {
		logger.Logger.Error(errorMsg, zap.Error(errors.WithStack(err)))
	}
}
