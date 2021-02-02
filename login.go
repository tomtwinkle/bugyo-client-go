package bugyo_client_go

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (b *bugyoClient) Login() error {
	// Get token
	if err := b.getToken(); err != nil {
		return err
	}

	// Check authentication
	if err := b.checkAuthentication(); err != nil {
		return err
	}

	// authenticate
	if err := b.authenticate(); err != nil {
		return err
	}

	// top page redirect, get user code
	if err := b.getUserCode(); err != nil {
		return err
	}

	return nil
}

func (b *bugyoClient) getToken() error {
	// https://id.obc.jp/{tenantCode}
	uri := fmt.Sprintf("%s/%s", baseUri, b.config.TenantCode)
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

func (b *bugyoClient) checkAuthentication() error {
	// https://id.obc.jp/{tenantCode}/login/CheckAuthenticationMethod
	uri := fmt.Sprintf("%s/%s", baseUri, b.config.TenantCode)
	endpoint := "login/CheckAuthenticationMethod"

	body := url.Values{}
	body.Set("OBCiD", b.config.OBCiD)
	body.Set("isBugyoCloud", "false")

	_, err := b.post(uri, endpoint, body)
	if err != nil {
		return err
	}
	return nil
}

func (b *bugyoClient) authenticate() error {
	// https://id.obc.jp/{tenantCode}/login/login/?Length=5
	uri := fmt.Sprintf("%s/%s", baseUri, b.config.TenantCode)
	endpoint := "login/login/?Length=5"

	body := url.Values{}
	body.Set("OBCiD", b.config.OBCiD)
	body.Set("Password", b.config.Password)
	body.Set("__RequestVerificationToken", b.token)
	body.Set("X-Requested-With", "XMLHttpRequest")

	_, err := b.post(uri, endpoint, body)
	if err != nil {
		return err
	}
	return nil
}

func (b *bugyoClient) getUserCode() error {
	// https://id.obc.jp/{tenantCode}/omredirect/redirect/
	uri := fmt.Sprintf("%s/%s/omredirect/redirect/", baseUri, b.config.TenantCode)
	doc, err := b.get(uri)
	if err != nil {
		return err
	}
	b.userCode = ""
	if homeUri, ok := doc.Find("#ApplicationRoot").Attr("href"); ok {
		// /{tenantCode}/{userCode}
		bCodes := strings.Split(homeUri, "/")
		if len(bCodes) == 3 {
			b.userCode = bCodes[2]
		}
	} else {
		return errors.New("no user code")
	}
	if token, ok := doc.Find("input[name=__RequestVerificationToken]").Attr("value"); ok {
		b.token = token
	} else {
		return errors.New("no token")
	}
	return nil
}

func (b *bugyoClient) IsLoggedIn() bool {
	if err := b.getUserCode(); err != nil {
		return false
	}
	return b.userCode != ""
}
