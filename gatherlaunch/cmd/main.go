package main

import (
	"faross/gatherlaunch"
	"faross/gatherlaunch/util"
	"fmt"
	"os"
)

func main() {
	purl, err := util.GetPurl(os.Args[1])
	if err != nil {
		return
	}
	gatherlaunch.InitGatherLaunch("./instruments.json")
	decision, err := gatherlaunch.Scan(purl)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(decision)
	}
}
