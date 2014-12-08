package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Fetcher struct {
	BaseUrl string
	Service string
	What    string
	Where   string
	PageLen int
	NbrPage int
	Dist    int
	Format  string
	Lang    string
	UID     string
	Key     string
}

type Payload struct {
	Summary  Summary  `json:"summary"`
	Listings []Entity `json:"listings"`
}

type Summary struct {
	What            string `json:"what"`
	Where           string `json:"where"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	FirstListing    int    `json:"firstListing"`
	LastListing     int    `json:"lastListing"`
	TotalListings   int    `json:"totalListings"`
	PageCount       int    `json:"pageCount"`
	CurrentPage     int    `json:"currentPage"`
	ListingsPerPage int    `json:"listingsPerPage"`
	Prov            string `json:"prov"`
}

type Entity struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Address  Address `json:"address"`
	GeoCode  GeoCode `json:"geocode"`
	Distance string  `json:"distance"`
	ParentId string  `json:"parentId"`
	IsParent bool    `json:"isParent"`
	Content  Content `json:"content"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Prov   string `json:"prov"`
	PCode  string `json:"pCode"`
}

type GeoCode struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Content struct {
	Video   Watcher `json:"video"`
	Photo   Watcher `json:"photo"`
	Profile Watcher `json:"profile"`
	DspAd   Watcher `json:"dspAd"`
	Logo    Watcher `json:"logo"`
	Url     Watcher `json:"url"`
}

type Watcher struct {
	Avail bool `json:"avail"`
	InMkt bool `json:"inMkt"`
}

func NewYPage(base, serv, what, where, uid, key string) *Fetcher {
	fetch := &Fetcher{
		BaseUrl: base,
		Service: serv,
		What:    what,
		Where:   where,
		PageLen: 700,
		NbrPage: 1,
		Dist:    3,
		Format:  "JSON",
		Lang:    "en",
		UID:     uid,
		Key:     key,
	}

	return fetch
}

func (f *Fetcher) Fetch() []byte {
	//DEBUG: fmt.Println("Fetcher Start")
	p := &Payload{}
	c := http.Client{}
	req := f.builder()
	resp, err := c.Do(req)
	if err != nil {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(p)
	if err != nil {
		return nil
	}
	raw, _ := json.MarshalIndent(p, "", "  ")
	//DEBUG: fmt.Println("Data Len:", len(raw))
	return raw
}

func (f *Fetcher) Size() int {
	return f.FindBusiness()
}

func (f *Fetcher) builder() *http.Request {
	//DEBUG: fmt.Println("Build Request")
	url := fmt.Sprintf("%s/%s/?what=%s&where=%s&pgLen=%d&pg=%d&dist=%d&fmt=%s&lang=%s&UID=%s&apikey=%s",
		f.BaseUrl,
		f.Service,
		f.What,
		f.Where,
		f.PageLen,
		f.NbrPage,
		f.Dist,
		f.Format,
		f.Lang,
		f.UID,
		f.Key,
	)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	return r
}

func (f *Fetcher) FindBusiness() int {
	//DEBUG: fmt.Println("Finding Business")
	c := http.Client{}
	p := &Payload{}
	req := f.builder()
	resp, err := c.Do(req)
	if err != nil {
		return 0
		fmt.Println("ERROR AT FINDBUSINESS FUNC")
	}
	err = json.NewDecoder(resp.Body).Decode(p)
	if err != nil {
		return 0
	}
	return p.Summary.PageCount
}

func (p *Payload) Readable() (string, error) {
	stack, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(stack), nil
}

func (p *Payload) Writable() ([]byte, error) {
	stack, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, err
	}
	return stack, nil
}
