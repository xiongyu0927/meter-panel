package configs

import "log"

func Init() {
	log.SetPrefix("[meter-panel]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
