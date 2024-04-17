package logs

import (
	"testing"
	"tyto/core/logging"
	"tyto/core/logs/mini"
)

func f11(logger logging.Logger) {
	logger.Error("test error, test error, test error, test error")
}

func f10(logger logging.Logger) {
	f11(logger)
}

func f9(logger logging.Logger) {
	f10(logger)
}

func f8(logger logging.Logger) {
	f9(logger)
}

func f7(logger logging.Logger) {
	f8(logger)
}

func f6(logger logging.Logger) {
	f7(logger)
}

func f5(logger logging.Logger) {
	f6(logger)
}

func f4(logger logging.Logger) {
	f5(logger)
}

func f3(logger logging.Logger) {
	f4(logger)
}

func f2(logger logging.Logger) {
	f3(logger)
}

func f1(logger logging.Logger) {
	f2(logger)
}

func TestLoggerImpl_Error(t *testing.T) {
	mlogger := mini.NewLogger()

	consoleSink := NewDefaultConsoleSink(mlogger)
	if consoleSink == nil {
		t.Error("NewDefaultConsoleSink failed")
		return
	}

	fileSink := NewDefaultFileSink(mlogger, "r:\\", "test.log")
	if fileSink == nil {
		t.Error("NewDefaultFileSink failed")
		return
	}

	logger := NewLoggerImpl(consoleSink, fileSink)
	defer logger.Close()

	for i := 0; i < 5; i++ {
		f1(logger)
	}
}

// 性能测试
// --- PASS: Test1 (2.13s)
// === RUN   Test1/info
// --- PASS: Test1/info (0.37s)
// === RUN   Test1/error
// --- PASS: Test1/error (1.76s)
func Test1(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		logger := NewLoggerImpl(NewDefaultFileSink(mini.NewLogger(), "r:\\", "test1.log"))
		defer logger.Close()

		for i := 0; i < 500000; i++ {
			logger.Info("this is a debug log message, this is a debug log message,", 10000000, 20000000, 30000000, 40000000, "5 6 7 8 9 10")
		}
	})

	t.Run("error", func(t *testing.T) {
		logger := NewLoggerImpl(NewDefaultFileSink(mini.NewLogger(), "r:\\", "test2.log"))
		defer logger.Close()

		for i := 0; i < 500000; i++ {
			logger.Error("this is a debug log message, this is a debug log message,", 10000000, 20000000, 30000000, 40000000, "5 6 7 8 9 10")
		}
	})
}

func Benchmark2(b *testing.B) {
	b.Run("logs-normal", func(b *testing.B) {
		logger := NewLoggerImpl(NewDiscardSink(mini.NewLogger()))
		defer logger.Close()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("this is a debug log message, this is a debug log message,", 10000000, 20000000, 30000000, 40000000, "5 6 7 8 9 10")
			}
		})
	})
}
