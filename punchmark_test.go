package bugyoclient

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBugyoClient_Punchmark(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "" {
		t.Log("Running on CI")
		t.SkipNow()
	}
	if err := godotenv.Load(".env"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	var config BugyoConfig
	if err := envconfig.Process("", &config); err != nil {
		t.Error(err)
		t.FailNow()
	}

	c, err := NewClient(&config, WithDebug())
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Run("clockin success", func(t *testing.T) {
		err := c.Login()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		err = c.Punchmark(ClockTypeClockIn)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	})

	t.Run("goout success", func(t *testing.T) {
		err := c.Login()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		err = c.Punchmark(ClockTypeGoOut)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	})

	t.Run("returned success", func(t *testing.T) {
		err := c.Login()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		err = c.Punchmark(ClockTypeReturned)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	})

	t.Run("clockout success", func(t *testing.T) {
		err := c.Login()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		err = c.Punchmark(ClockTypeClockOut)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	})
}
