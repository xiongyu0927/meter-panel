package configs

import "log"

func init() {
	log.SetPrefix("[meter-panel]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
