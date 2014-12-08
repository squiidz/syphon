package syphon

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/squiidz/nexus"
)

type Provider interface {
	Fetch() []byte
	Size() int
}

type Syphon struct {
	Provider
	Nexus  *nexus.Nexus
	worker int
	At     int
}

func NewSyphon(p Provider, w int) *Syphon {
	return &Syphon{p, nexus.New(), w, 0}
}

func (s *Syphon) Do(fname string) {
	for i := 0; i < s.worker; i++ {
		s.Nexus.NewProbe().NewJob(s.Fetch)

	}
	s.Nexus.Start()
	s.ViewData()

	s.WriteFile(fname)
}

func (s *Syphon) ViewData() {
	fmt.Println("######## Stats #######")
	for _, p := range s.Nexus.Probes {
		b := bytes.NewBuffer([]byte(""))

		stats := p.Extract(b)
		fmt.Println(stats)
	}
	fmt.Println("")
}

func (s *Syphon) WriteFile(filename string) error {
	for _, p := range s.Nexus.Probes {
		b := bytes.NewBuffer([]byte(""))
		p.Extract(b)

		err := ioutil.WriteFile(filename, b.Bytes(), 0600)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
