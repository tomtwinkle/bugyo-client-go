package bugyoclient

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
	if ref := b.refererForURL(b.lastReq); ref != "" {
		req.Header.Set("Referer", ref)
	}

	defer b.setLastReq(req.URL)

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
	if ref := b.refererForURL(b.lastReq); ref != "" {
		req.Header.Set("Referer", ref)
	}
	defer b.setLastReq(req.URL)

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
