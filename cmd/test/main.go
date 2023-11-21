package main

import (
	"github.com/kehiy/RoboPac/client"
	"github.com/yudai/pp"
)

func main() {
	ip := "172.104.46.145:9090"
	c, _ := client.NewClient(ip)
	res, _ := c.IsValidator("tpc1pd9xmumgzsqd0mnmy3r5dvsku7d6xuxesmanwky")
	pp.Println(res)
}
