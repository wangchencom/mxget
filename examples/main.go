package main

import (
	"fmt"

	"github.com/winterssy/mxget/pkg/provider/netease"
)

func main() {
	client := netease.New(nil)
	resp, err := client.GetSong("36990266")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
