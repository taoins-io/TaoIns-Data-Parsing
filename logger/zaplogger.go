package logger

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	log "gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"sync"
	"tao/config"
	"time"
)

type Logger struct {
	L *zap.Logger
	*zap.SugaredLogger
	sync.Mutex
	Opts      *config.Log `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

var (
	logger                         *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer
	debugConsoleWS                 = zapcore.Lock(os.Stdout)
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

func (l *Logger) Printf(message string, data ...interface{}) {
	//if len(data) >= 2  {
	//	if data[0] == "panic" || data[1] == "error"{
	//		logger.Errorf(message, data)
	//	}
	//}
	logger.Infof(message, data)
}

func InitLogger() {
	data := &config.Config.Log
	if data.Level == "debug" {
		data.Config.Development = true
	} else {
		data.Config.Development = false
	}
	logger = &Logger{
		Opts: data,
	}
	if logger.inited {
		logger.Info("[initLogger] logger Inited")
		return
	}
	logger.loadCfg()
	logger.init()
	logger.Info("[initLogger] zap plugin initializing completed")
	logger.inited = true
}

// GetLogger returns logger
func GetLogger() (ret *Logger) {
	return logger
}

func (l *Logger) init() {
	l.setSyncers()
	var err error
	mylogger, err := l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	l.Opts = &config.Config.Log
	l.L = mylogger
	l.SugaredLogger = mylogger.Sugar()
	defer l.SugaredLogger.Sync()
}

func (l *Logger) loadCfg() {
	if l.Opts.Config.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}
	if l.Opts.Config.OutputPaths == nil || len(l.Opts.Config.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}
	if l.Opts.Config.ErrorOutputPaths == nil || len(l.Opts.Config.ErrorOutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stderr"}
	}
	if len(l.Opts.LogFileDir) == 0 {
		l.Opts.LogFileDir += sp + "logs"
	}
	os.MkdirAll(l.Opts.LogFileDir, os.ModePerm)
	if l.Opts.AppName == "" {
		l.Opts.AppName = "app"
	}
	if l.Opts.ErrorFileName == "" {
		l.Opts.ErrorFileName = "error.log"
	}
	if l.Opts.WarnFileName == "" {
		l.Opts.WarnFileName = "warn.log"
	}
	if l.Opts.InfoFileName == "" {
		l.Opts.InfoFileName = "info.log"
	}
	if l.Opts.DebugFileName == "" {
		l.Opts.DebugFileName = "debug.log"
	}
	if l.Opts.MaxSize == 0 {
		l.Opts.MaxSize = 100
	}
	if l.Opts.MaxBackups == 0 {
		l.Opts.MaxBackups = 30
	}
	if l.Opts.MaxAge == 0 {
		l.Opts.MaxAge = 30
	}
}

func (l *Logger) setSyncers() {
	f := func(fN string) zapcore.WriteSyncer {
		logf, _ := rotatelogs.New(l.Opts.LogFileDir+sp+l.Opts.AppName+"-"+fN+".%Y_%m%d",
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(time.Minute),
		)
		return zapcore.AddSync(logf)
	}
	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	return
}

func (l *Logger) cores() zap.Option {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})
	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, errWS, errPriority),
		zapcore.NewCore(consoleEncoder, warnWS, warnPriority),
		zapcore.NewCore(consoleEncoder, infoWS, infoPriority),
		zapcore.NewCore(consoleEncoder, debugWS, debugPriority),
	}
	if l.Opts.Config.Development {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		}...)
	} else {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}

type writer struct {
	log.Writer
}

// NewWriter writer
func NewWriter(w log.Writer) *writer {
	return &writer{Writer: w}
}

func (w *writer) Printf(message string, data ...interface{}) {
	logger.Infof(message, data...)

}
