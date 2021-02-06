package bugyoclient

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
)

const userAgent = "Bugyo-Client-Go/1.0.0"

const (
	urlLoginPage                 = "https://id.obc.jp/%s"
	urlCheckAuthenticationMethod = "https://id.obc.jp/%s/login/CheckAuthenticationMethod"
	urlAuthenticate              = "https://id.obc.jp/%s/login/login/?Length=5"
	urlUserCode                  = "https://id.obc.jp/%s/omredirect/redirect/"
	urlPunchmarkPage             = "https://hromssp.obc.jp/%s/%s/timeclock/punchmark/"
	urlInsertReadDateTime        = "https://hromssp.obc.jp/%s/%s/TimeClock/InsertReadDateTime/"
)

type BugyoClient interface {
	Login() error
	IsLoggedIn() bool
	Punchmark(clockType ClockType) error
}

type bugyoClient struct {
	client   *http.Client
	config   *BugyoConfig
	token    string
	userCode string
	lastReq  *url.URL
	debug    bool
}

type BugyoConfig struct {
	TenantCode string
	OBCiD      string
	Password   string
}

type Options func(*options)

type options struct {
	debug bool
}

func WithDebug() Options {
	return func(ops *options) {
		ops.debug = true
	}
}

func NewClient(config *BugyoConfig, opts ...Options) (BugyoClient, error) {
	var opt options
	for _, o := range opts {
		o(&opt)
	}
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
	return &bugyoClient{client: client, config: config, debug: opt.debug}, nil
}

func (b *bugyoClient) get(uri string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	defer b.setLastReq(req.URL)

	if b.debug {
		reqDump, err := httputil.DumpRequest(req, false)
		if err != nil {
			return nil, err
		}
		log.Printf("request=%q\n", reqDump)
	}
	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if b.debug {
		resDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
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

// nolint:unparam
func (b *bugyoClient) post(uri string, body url.Values) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("__RequestVerificationToken", b.token)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	if ref := b.refererForURL(b.lastReq); ref != "" {
		req.Header.Set("Referer", ref)
	}
	defer b.setLastReq(req.URL)

	if b.debug {
		reqDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		log.Printf("request=%q\n", reqDump)
	}
	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if b.debug {
		resDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
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

	b.lastReq = req.URL
	return doc, nil
}

func (b *bugyoClient) refererForURL(lastReq *url.URL) string {
	if lastReq == nil {
		return ""
	}
	return lastReq.String()
}

func (b *bugyoClient) setLastReq(lastReq *url.URL) {
	b.lastReq = lastReq
}
