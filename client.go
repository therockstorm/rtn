package rtn

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	agreementURL = "https://www.frbservices.org/EPaymentsDirectory/submitAgreement"
	fileURL      = "https://www.frbservices.org/EPaymentsDirectory/FedACHdir.txt"
)

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type client struct {
	http   httpClient
	reader func(r io.Reader) ([]byte, error)
}

func newClient(cookieJarCreator func(o *cookiejar.Options) (*cookiejar.Jar, error)) *client {
	jar, err := cookieJarCreator(nil)
	if err != nil {
		log.Println(err)
		panic("CookieJar required.")
	}

	return &client{&http.Client{Jar: jar}, ioutil.ReadAll}
}

func (c *client) getFile() (string, error) {
	err := c.agree()
	if err != nil {
		return "", err
	}

	res, err := c.http.Get(fileURL)
	if err != nil {
		log.Println("GET file failed.")
		return "", err
	}
	defer res.Body.Close()

	body, err := c.reader(res.Body)
	if err != nil {
		log.Println("Failed to read response body.")
		return "", err
	}

	strBody := string(body)
	if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected statusCode=%v, body='%v'", res.StatusCode, strBody)
	}

	return strBody, err
}

func (c *client) agree() error {
	form := url.Values{}
	form.Add("agreementValue", "Agree")

	res, err := c.http.PostForm(agreementURL, form)
	if err != nil {
		log.Println("POST agreement failed.")
		return err
	}
	defer res.Body.Close()

	return nil
}
