/**
 *  Manu Cutillas Personal Website Server
 *
 *  created_at May 20, 2018
 *  author: Manu Cutillas<manucutillas@outlook.com>
 *  license: MIT
 */
package main

import (
	"flag"
	"./modules/log"
	"./router"
)

const (
	DefaultConfFilePath = "config/config.toml"
)

var (
	confFilePath string
	cmdHelp      bool
)

func init() {
	flag.StringVar(&confFilePath, "c", DefaultConfFilePath, "Configuration file path")
	flag.BoolVar(&cmdHelp, "h", false, "help")
	flag.Parse()

}

func main() {
	if cmdHelp {
		flag.PrintDefaults()
		return
	}
	log.Debugf("run with conf:%s", confFilePath)

	// Subdomain deployment
	router.RunSubdomains(confFilePath)
}




