package client

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
)

const userAgent = "Bugyo-Client-Go/1.0.0"
const baseUri = "https://id.obc.jp"
const tokenSelector = "input[name=__RequestVerificationToken]"

type BugyoClient interface {
}

type bugyoClient struct {
	client *http.Client
	config *BugyoConfig
	token  string
	debug  bool
}

type BugyoConfig struct {
	TenantCode string
	OBCiD      int
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
	if config.OBCiD == 0 {
		return nil, errors.New("invalid argument: OBCiD required")
	}
	if config.Password == "" {
		return nil, errors.New("invalid argument: Password required")
	}
	return bugyoClient{client: client, config: config, debug: debug}, nil
}

func (b *bugyoClient) Login() error {
	// Get token
	if err := b.getToken(); err != nil {
		return err
	}

	// Check authentication

	// Login

	return nil
}

func (b *bugyoClient) getToken() error {
	// https://id.obc.jp/{tenantCode}
	uri := fmt.Sprintf("%s/%s", baseUri, b.config.TenantCode)
	doc, err := b.get(uri)
	if err != nil {
		return err
	}
	if token, ok := doc.Find(tokenSelector).Attr("value"); ok {
		b.token = token
		return nil
	}
	return errors.New("no token")
}

func (b *bugyoClient) Punchmark() {

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

func (b *bugyoClient) post(domain, endpoint string, body interface{}) (*goquery.Document, error) {
	uri := fmt.Sprintf("https://%s/%s", domain, endpoint)
	req, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
