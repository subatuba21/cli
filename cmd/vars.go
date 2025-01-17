package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (

	// initialize console colours
	Bold  = "\033[1m"
	Reset = "\033[0m"
	Green = "\033[32m"
	// Blue = "\033[34m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Red    = "\033[31m"
	// Gray = "\033[37m"
	// White = "\033[97m"
)

var (
	// Utility build version
	Version string

	cfgFile string
	log     = logrus.New()
	DEBUG   bool
	JSON    bool

	LOG_FILE = ""

	// store Hasura console session command,
	// to kill it later while shutting down dev environment
	hasuraConsoleSpawnProcess *os.Process
)
