package client

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
)

const userAgent = "Bugyo-Client-Go/1.0.0"
const baseUri = "https://id.obc.jp"

type BugyoClient interface {
	Login() error
	IsLoggedIn() bool
}

type bugyoClient struct {
	client   *http.Client
	config   *BugyoConfig
	token    string
	userCode string
	debug    bool
}

type BugyoConfig struct {
	TenantCode string
	OBCiD      string
	Password   string
}

func NewClient(config *BugyoConfig, debug bool) (BugyoClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}
	if config == nil {
		return nil, errors.New("invalid argument: config required")
	}
	if config.TenantCode == "" {
		return nil, errors.New("invalid argument: TenantCode required")
	}
	if config.OBCiD == "" {
		return nil, errors.New("invalid argument: OBCiD required")
	}
	if config.Password == "" {
		return nil, errors.New("invalid argument: Password required")
	}
	return &bugyoClient{client: client, config: config, debug: debug}, nil
}

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

func (b *bugyoClient) Punchmark() error {
	// move PunchMark page
	// GET: https://hromssp.obc.jp/{tenantCode}/{userCode}/timeclock/punchmark/

	// ClockIn
	// POST https://hromssp.obc.jp/{tenantCode}/{userCode}/TimeClock/InsertReadDateTime/
	//headers={'Referer': 'https://hromssp.obc.jp/{tenantCode}/{userCode}/timeclock/punchmark/', '__RequestVerificationToken': '<token>', 'X-Requested-With': 'XMLHttpRequest'}
	//data={'ClockType': 'ClockIn', 'LaborSystemID': '0', 'LaborSystemCode': '', 'LaborSystemName': '', 'PositionLatitude': <CompanyLatitude>, 'PositionLongitude': <CompanyLongitude>, 'PositionAccuracy': '0'}

	return nil
}

func (b *bugyoClient) IsLoggedIn() bool {
	if err := b.getUserCode(); err != nil {
		return false
	}
	return b.userCode != ""
}

func (b *bugyoClient) get(uri string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)

	if b.debug {
		reqDump, _ := httputil.DumpRequest(req, false)
		log.Printf("request=%q\n", reqDump)
	}
	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if b.debug {
		resDump, _ := httputil.DumpResponse(res, true)
		log.Printf("response=%q\n", resDump)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New(res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (b *bugyoClient) post(domain, endpoint string, body url.Values) (*goquery.Document, error) {
	uri := fmt.Sprintf("%s/%s", domain, endpoint)
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("__RequestVerificationToken", b.token)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	if b.debug {
		reqDump, _ := httputil.DumpRequest(req, true)
		log.Printf("request=%q\n", reqDump)
	}
	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if b.debug {
		resDump, _ := httputil.DumpResponse(res, true)
		log.Printf("response=%q\n", resDump)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New(res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
