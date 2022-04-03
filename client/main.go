package main

import (
	"file_trans/client/cmd"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func main() {
	cmd.Execute()
	Logger, _ = zap.NewProduction()
}
