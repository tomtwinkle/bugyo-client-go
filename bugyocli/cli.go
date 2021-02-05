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

func NewCLI() CLI {
	cfg := config.NewConfig()
	bCfg, err := cfg.Init()
	if err != nil {
		log.Fatal(err)
	}
	c, err := bugyoclient.NewClient(bCfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cli{c}
}

func (c cli) PunchMark(clockType bugyoclient.ClockType) error {
	if err := c.Login(); err != nil {
		return err
	}
	if err := c.PunchMark(clockType); err != nil {
		return err
	}
	return nil
}
