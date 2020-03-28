package main

import (
	"github.com/chebizarro/redshed/config"
	"github.com/chebizarro/redshed/internal/api"
	"strconv"
)

func main() {
	r := api.InitRouter()
	// Run the server
	port := strconv.Itoa(config.Config.Server.Port)
	r.Run(":" + port)
}