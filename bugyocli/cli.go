package bugyocli

import (
	bugyoclient "github.com/tomtwinkle/bugyo-client-go"
	"github.com/tomtwinkle/bugyo-client-go/config"
	"log"
)

type cli struct {
	bugyoclient.BugyoClient
}

type CLI interface {
	PunchMark(clockType bugyoclient.ClockType) error
}

func NewCLI(verbose bool) CLI {
	cfg := config.NewConfig()
	bCfg, err := cfg.Init()
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		c, err := bugyoclient.NewClient(bCfg, bugyoclient.WithDebug())
		if err != nil {
			log.Fatal(err)
		}
		return &cli{c}
	}
	c, err := bugyoclient.NewClient(bCfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cli{c}
}

func (c cli) PunchMark(clockType bugyoclient.ClockType) error {
	if err := c.BugyoClient.Login(); err != nil {
		return err
	}
	if err := c.BugyoClient.Punchmark(clockType); err != nil {
		return err
	}
	log.Printf("success Punchmark [%s]", clockType)
	return nil
}
