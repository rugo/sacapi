package apilog

import (
    "os"
    "github.com/op/go-logging"
)

var Log = logging.MustGetLogger("scapi")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
    "%{color}%{time:15:04:05.000} %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func Init() {
    backend := logging.NewLogBackend(os.Stderr, "", 0)

    backendFormatter := logging.NewBackendFormatter(backend, format)

    backendLevel := logging.AddModuleLevel(backend)
    backendLevel.SetLevel(logging.DEBUG, "")

    logging.SetBackend(backendLevel, backendFormatter)
}