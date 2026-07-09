package zap

import "go.uber.org/zap"

type ZapLoggerCover struct {
	logger *zap.Logger
}

func NewZapLoggerCover(logger *zap.Logger) ZapLoggerCover {
	return ZapLoggerCover{
		logger: logger,
	}
}

func (z ZapLoggerCover) Debug(msg string) {
	z.logger.Debug(msg)
}

func (z ZapLoggerCover) Info(msg string) {
	z.logger.Info(msg)
}

func (z ZapLoggerCover) Warn(msg string) {
	z.logger.Warn(msg)
}

func (z ZapLoggerCover) Error(msg string) {
	z.logger.Error(msg)
}
