package parser

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

type data struct {
	INN         string `json:"INN"`
	KPP         string `json:"KPP"`
	Leader      string `json:"Leader"`
	CompanyName string `json:"company_name"`
}

func (e *data) Export(exports chan interface{}) error {
	for res := range exports {
		js, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(js, e)
	}
	return nil
}
func GoParser(inn string) *data {
	var ex data
	url := fmt.Sprintf(`https://www.rusprofile.ru/search?query=%s+&search_inactive=2`, inn)
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs:         []string{url},
		RobotsTxtDisabled: true,
		ParseFunc:         parseFunc,
		Exporters:         []export.Exporter{&ex},
		LogDisabled:       true,
	}).Start()

	return &ex
}

func parseFunc(g *geziyor.Geziyor, r *client.Response) {

	r.HTMLDoc.Find(`div#anketa`).Each(func(i int, s *goquery.Selection) {
		leader := s.Find("a.gtm_main_fl").Text()
		if leader == "" {
			leader = s.Find("div.company-row.hidden-parent").Find("span.company-info__text").Text()
		}
		g.Exports <- map[string]interface{}{

			"INN":          s.Find("span#clip_inn").Text(),
			"KPP":          s.Find("span#clip_kpp").Text(),
			"leader":       leader,
			"company_name": s.Find("div.company-name").Text(),
		}

	})

}
