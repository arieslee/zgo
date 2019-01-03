// TODO log time rotate
// TODO log size rotate
// TODO multi logger
// TODO windows color
// TODO check tty

package log

import (
    "fmt"
    "log"
    "os"
    "sync"
)

const (
    LevelDebug = (iota + 1) * 10
    LevelInfo
    LevelWarn
    LevelError
    LevelFatal
)

const (
    TagDebug = "[D]"
    TagInfo  = "[I]"
    TagWarn  = "[W]"
    TagError = "[E]"
    TagFatal = "[F]"
    TagLog   = "[L]"
)

// ref: https://en.wikipedia.org/wiki/ANSI_escape_code
const (
    ColorDebug = "\033[37m"
    ColorInfo  = "\033[36m"
    ColorWarn  = "\033[33m"
    ColorError = "\033[31m"
    ColorFatal = "\033[35m"
    ColorReset = "\033[0m"
)

const (
    RotateNo = iota
    RotateTimeDay
    RotateTimeHour
    RotateTimeMinute
    RotateTimeSecond
    RotateSizeKB
    RotateSizeMB
    RotateSizeGB
)

type LevelLogger struct {
    *log.Logger
    fp *os.File
    mu sync.Mutex
    Prefix string
    Filename  string
    Level int
    Rotate int
    MaxSize int
    MaxBackup int
}

var (
    defaultLog *LevelLogger
)

func GetLogger() *LevelLogger {
    if defaultLog == nil {
        fmt.Println("can not GetLogger before Install")
        os.Exit(1)
    }
    return defaultLog
}

func Install(dest string) *LevelLogger {
    if defaultLog != nil {
        defaultLog.Warn("can not install logger twice !!!")
        return defaultLog
    }

    var base *log.Logger
    var fp *os.File

    if dest == "stdout" {
        fp = os.Stdout
        base = log.New(fp, "",log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    } else {
        fp, err := os.OpenFile(dest, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
        if err != nil {
            fmt.Printf("can not open logfile: %v\n", err)
        }
        base = log.New(fp, "",log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    }

    l := LevelLogger{
        Logger: base,
        fp: fp,
        Prefix: "",
        Filename:   dest,
        Level:  LevelDebug,
    }
    defaultLog = &l
    return &l
}

func (l *LevelLogger) Log(level int, depth int, format string, v ...interface{} ) {
    if level >= l.Level {
        var tag,color,message string
        switch level {
            case LevelDebug:
                tag = TagDebug
                color = ColorDebug
            case LevelInfo:
                tag = TagInfo
                color = ColorInfo
            case LevelWarn:
                tag = TagWarn
                color = ColorWarn
            case LevelError:
                tag = TagError
                color = ColorError
            case LevelFatal:
                tag = TagFatal
                color = ColorFatal
            default:
                tag = TagLog
                color = ColorReset
        }
        if len(v) == 0 {
            message = format
        } else {
            message = fmt.Sprintf(format, v...)
        }
        if l.Filename == "stdout" {
            // XXX debug only, slow with 4 lock
            l.mu.Lock()
            l.Logger.SetPrefix(color)
            l.Logger.Output(depth, fmt.Sprint(tag, " ", l.Prefix, message, ColorReset))
            l.Logger.SetPrefix("")
            l.mu.Unlock()
        } else {
            l.Logger.Output(depth, fmt.Sprint(tag, " ", l.Prefix, message))
        }
    }
}

func (l *LevelLogger) SetPrefix(prefix string) {
    l.Prefix = prefix
}

func (l *LevelLogger) SetLevel(level int) {
    l.Level = level
}

func (l *LevelLogger) Debug(format string, v ...interface{}) {
    l.Log(LevelDebug, 3, format, v...)
}

func (l *LevelLogger) Info(format string, v ...interface{}) {
    l.Log(LevelInfo, 3, format, v...)
}

func (l *LevelLogger) Warn(format string, v ...interface{}) {
    l.Log(LevelWarn, 3, format, v...)
}

func (l *LevelLogger) Error(format string, v ...interface{}) {
    l.Log(LevelError, 3, format, v...)
}

func (l *LevelLogger) Fatal(format string, v ...interface{}) {
    l.Log(LevelFatal, 3, format, v...)
    os.Exit(1)
}

func Debug(format string, v ...interface{}) {
    defaultLog.Log(LevelDebug, 3, format, v...)
}

func Info(format string, v ...interface{}) {
    defaultLog.Log(LevelInfo, 3, format, v...)
}

func Warn(format string, v ...interface{}) {
    defaultLog.Log(LevelWarn, 3, format, v...)
}

func Error(format string, v ...interface{}) {
    defaultLog.Log(LevelError, 3, format, v...)
}

func Fatal(format string, v ...interface{}) {
    defaultLog.Log(LevelFatal, 3, format, v...)
    os.Exit(1)
}

func Debugd(depth int, format string, v ...interface{}) {
    defaultLog.Log(LevelDebug, 3+depth, format, v...)
}

func Infod(depth int, format string, v ...interface{}) {
    defaultLog.Log(LevelInfo, 3+depth, format, v...)
}

func Warnd(depth int, format string, v ...interface{}) {
    defaultLog.Log(LevelWarn, 3+depth, format, v...)
}

func Errord(depth int, format string, v ...interface{}) {
    defaultLog.Log(LevelError, 3+depth, format, v...)
}

func Fatald(depth int, format string, v ...interface{}) {
    defaultLog.Log(LevelFatal, 3+depth, format, v...)
    os.Exit(1)
}