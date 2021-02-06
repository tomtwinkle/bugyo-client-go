package bugyoclient

import (
	"errors"
	"fmt"
	"net/url"
)

type ClockType string

const (
	ClockTypeClockIn  = ClockType("ClockIn")
	ClockTypeClockOut = ClockType("ClockOut")
	ClockTypeGoOut    = ClockType("GoOut")
	ClockTypeReturned = ClockType("Returned")
)

func (b *bugyoClient) Punchmark(clockType ClockType) error {
	// move PunchMark page
	if err := b.movePunchMarkPage(); err != nil {
		return err
	}

	// ClockIn/ClockOut
	if err := b.insertReadDateTime(clockType); err != nil {
		return err
	}

	return nil
}

func (b *bugyoClient) movePunchMarkPage() error {
	if b.userCode == "" {
		return errors.New("login required")
	}
	// GET: https://hromssp.obc.jp/{tenantCode}/{userCode}/timeclock/punchmark/
	uri := fmt.Sprintf(urlPunchmarkPage, b.config.TenantCode, b.userCode)
	doc, err := b.get(uri)
	if err != nil {
		return err
	}
	if token, ok := doc.Find("input[name=__RequestVerificationToken]").Attr("value"); ok {
		b.token = token
		return nil
	}
	return errors.New("no token")
}

func (b *bugyoClient) insertReadDateTime(clockType ClockType) error {
	if b.userCode == "" {
		return errors.New("login required")
	}
	// POST https://hromssp.obc.jp/{tenantCode}/{userCode}/TimeClock/InsertReadDateTime/
	uri := fmt.Sprintf(urlInsertReadDateTime, b.config.TenantCode, b.userCode)

	body := url.Values{}
	body.Set("ClockType", string(clockType))
	body.Set("LaborSystemID", "0")
	body.Set("LaborSystemCode", "")
	body.Set("LaborSystemName", "")
	body.Set("PositionLatitude", "35.6812")   // FIXME: tokyo station
	body.Set("PositionLongitude", "139.7671") // FIXME: tokyo station
	body.Set("PositionAccuracy", "0")

	_, err := b.post(uri, body)
	if err != nil {
		return err
	}

	return nil
}
