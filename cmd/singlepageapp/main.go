package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/janmbaco/go-infrastructure/logs"
	"github.com/janmbaco/go-infrastructure/server/facades"
)

func main() {
	port := flag.String("port", ":8080", "port to listen on, like :8080")
	staticPath := flag.String("static", "./static", "path to static files")
	index := flag.String("index", "index.html", "index file name")
	flag.Parse()

	logger := logs.NewLogger()
	logger.Info(fmt.Sprintf("Iniciando servidor SPA en puerto %s, static: %s, index: %s", *port, *staticPath, *index))
	facades.SinglePageAppStart(*port, *staticPath, *index)
	fmt.Fprintln(os.Stderr, "server exited")
}
