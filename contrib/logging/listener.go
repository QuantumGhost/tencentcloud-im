package logging

import (
	"github.com/QuantumGhost/tencentcloud-im/events"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Listener struct {
	logger *zap.Logger
}

func (l Listener) Event(event events.EventType) {
	l.logger.Info("received event", zap.String("event", string(event)))
}

func (l Listener) EventKV(event events.EventType, kvs []events.KVPair) {
	l.logger.Info("received event", kvsToFields(kvs, zap.String("event", string(event))))
}

func (l Listener) Timing(event events.EventType, nanoseconds int64) {
	l.logger.Info("request timing",
		zap.String("event", string(event)),
		zap.Duration("duration", time.Duration(nanoseconds)))
}

func (l Listener) TimingKv(event events.EventType, nanoseconds int64, kvs []events.KVPair) {
	fields := kvsToFields(
		kvs, zap.String("event", string(event)), zap.Duration("duration", time.Duration(nanoseconds)))
	l.logger.Info("request timing", fields...)
}

func kvsToFields(kvs []events.KVPair, otherFields ...zapcore.Field) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(kvs)+len(otherFields))
	fields = append(fields, otherFields...)
	for _, kv := range kvs {
		fields = append(fields, zap.Any(kv.Key(), kv.Value()))
	}
	return fields
}
