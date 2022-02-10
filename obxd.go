package main

import (
	"github.com/cpacia/obxd/repo"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
)

func main() {
	var emptyCfg repo.Config
	parser := flags.NewNamedParser("obxd", flags.Default)
	parser.AddGroup("Node Options", "Configuration options for the node", &emptyCfg)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	cfg, err := repo.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, err := BuildServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	for sig := range c {
		if sig == os.Kill || sig == os.Interrupt {
			log.Info("obxd gracefully shutting down")
			server.Close()
			os.Exit(1)
		}
	}
}
