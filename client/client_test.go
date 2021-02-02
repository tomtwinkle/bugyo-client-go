package client_test

import (
	"bugyo-client-go/client"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBugyoClient_Login(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	var config client.BugyoConfig
	if err := envconfig.Process("", &config); err != nil {
		t.Error(err)
		t.FailNow()
	}

	c, err := client.NewClient(&config, true)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Run("login success", func(t *testing.T) {
		err := c.Login()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		loggedIn := c.IsLoggedIn()
		assert.True(t, loggedIn)
	})
}
