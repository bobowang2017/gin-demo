package log

import (
	"fmt"
	"gin-demo/core/settings"
	"gin-demo/infra/common"
	rotateLogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	Logger = log.New()
)

type CustomFormatter struct{}

// Format 自定义日志输出格式
func (s *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format(common.TimeLayout)
	var file string
	var line int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}
	msg := fmt.Sprintf("%s [%s] [%s:%d] %s\n", timestamp, strings.ToUpper(entry.Level.String()), file, line, entry.Message)
	//msg := fmt.Sprintf("%s [%s] [%s:%d] %s\n", timestamp, strings.ToUpper(entry.Level.String()), file, line, entry.Message)
	return []byte(msg), nil
}

// Setup initialize the log instance
func Setup() {
	_ = getLogFilePath()
	fileName := getLogFileName()

	writer, _ := rotateLogs.New(settings.Config.Server.LogSavePath+
		"%Y-%m-%d_%H-%M"+fileName,
		//rotateLogs.WithLinkName(logFile),
		rotateLogs.WithMaxAge(settings.Config.Server.LogMaxAge*24*time.Hour),
		rotateLogs.WithRotationTime(settings.Config.Server.LogRotationTime*time.Hour),
	)
	mw := io.MultiWriter(os.Stdout, writer)
	Logger.SetOutput(mw)
	//Logger.SetFormatter(new(CustomFormatter))
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:03:04",
		ForceColors:     true,
		FullTimestamp:   true,
		DisableQuote:    true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理函数名
			fs := strings.Split(frame.Function, ".")
			fun := ""
			if len(fs) > 0 {
				fun = fs[len(fs)-1]
			}
			fileName := path.Base(frame.File)
			return fmt.Sprintf("[\033[1;34m%s\033[0m]", fun), fmt.Sprintf("[%s:%d]", fileName, frame.Line)
		},
	})
	Logger.SetReportCaller(true)
	switch strings.ToLower(settings.Config.Server.LogLevel) {
	case "panic":
		Logger.SetLevel(log.PanicLevel)
	case "fatal":
		Logger.SetLevel(log.FatalLevel)
	case "warn":
		Logger.SetLevel(log.WarnLevel)
	case "info":
		Logger.SetLevel(log.InfoLevel)
	case "debug":
		Logger.SetLevel(log.DebugLevel)
	case "trace":
		Logger.SetLevel(log.TraceLevel)
	default:
		Logger.SetLevel(log.InfoLevel)
	}
	//logger.Hooks.Add(NewContextHook())
}

func getLogFilePath() string {
	var (
		logPath                 = settings.Config.Server.LogSavePath
		currentPath, logAbsPath string
		err                     error
	)
	if currentPath, err = filepath.Abs("."); err != nil {
		panic(err)
	}
	logAbsPath = filepath.Join(currentPath, logPath)
	if _, err := os.Stat(logAbsPath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(logAbsPath, 0755)
			if err != nil {
				panic(err)
			}
		}
	}
	return logPath
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s.%s",
		time.Now().Format(common.DateLayout),
		"log",
	)
}
