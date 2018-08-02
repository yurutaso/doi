package main

import (
	"encoding/json"
	"fmt"
	"github.com/caltechlibrary/bibtex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	DOI_DOMAIN string = "http://www.doi.org"
	DOI_API    string = "/api/handles/"
)

type DoiResponse struct {
	ResponseCode int                `json:"responseCode",omitempty`
	Handle       string             `json:"handle",omitempty`
	Values       []DoiResponseValue `json:"values",omitempty`
}

type DoiResponseValue struct {
	Index     int             `json:"index",omitempty`
	Type      string          `json:"type",omitempty`
	Ttl       int             `json:"ttl",omitempty`
	Timestamp string          `json:"timestamp",omitempty`
	Data      DoiResponseData `json:"data",omiempty`
}

type DoiResponseData struct {
	Format string      `json:"format",omitempty`
	Value  interface{} `json:"value",omiempty`
}

func SearchDOI(doi string) ([]byte, error) {
	u, err := url.Parse(DOI_DOMAIN)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, DOI_API, doi)
	fmt.Println(u.String())
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ExtractBib(r io.Reader) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	elements, err := bibtex.Parse(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(elements))
	for _, element := range elements {
		if strings.Contains(bibtex.DefaultInclude, element.Type) {
			//fmt.Println(element)
			fmt.Println(element.Keys)
			fmt.Println(element.Type)
			//fmt.Println(element.Tags)
		}
	}
}

func main() {
	text, err := SearchDOI("10.1088/0004-637X/811/2/118")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(text))

	var doi DoiResponse
	err = json.Unmarshal(text, &doi)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doi)

	r, err := os.Open(`./example.bib`)
	if err != nil {
		log.Fatal(err)
	}
	ExtractBib(r)
}
