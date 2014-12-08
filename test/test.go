package main

import (
	"time"

	"github.com/squiidz/nexus"
	"github.com/squiidz/syphon"
	"github.com/squiidz/syphon/provider"
)

func main() {
	yj := provider.NewYPage("http://api.sandbox.yellowapi.com", "FindBusiness", "Restaurant", "Joliette", "Dev", "eh2vk49jvdgmm66dymcre2xy")
	yl := provider.NewYPage("http://api.sandbox.yellowapi.com", "FindBusiness", "Restaurant", "Laval", "Dev", "qs9x872kthgk4aur4u6x2xr9")

	sypJ := syphon.NewSyphon(yj, 2)
	//sypJ.Nexus.SetStarter(SleepyStart)
	sypJ.Do("jolietteData.json")

	sypL := syphon.NewSyphon(yl, 1)
	sypL.Nexus.SetStarter(SleepyStart)
	sypL.Do("lavalData.json")
}

func SleepyStart(n *nexus.Nexus) {
	for _, p := range n.Probes {
		n.WaitStack.Add(1)
		go p.Work()
		time.Sleep(time.Second * 1)
	}
	n.WaitStack.Wait()
}
