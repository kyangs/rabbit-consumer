package rabbitmq

import (
	"os"
	"syscall"
)

var DeadSignal = []os.Signal{
	syscall.SIGTERM,
	syscall.SIGINT,
	syscall.SIGKILL,
	syscall.SIGHUP,
	syscall.SIGQUIT,
}
