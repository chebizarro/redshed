package main

import (
	"github.com/chebizarro/redshed/cmd/redshed/config"
	"github.com/chebizarro/redshed/internal/logger"

	"github.com/chebizarro/redshed/internal/orm"
	"github.com/chebizarro/redshed/pkg/server"
)

// main
func main() {
	sc := config.Server()
	o, err := orm.Factory(sc)
	defer o.DB.Close()
	if err != nil {
		logger.Panic(err)
	}
	server.Run(sc, o)
}
