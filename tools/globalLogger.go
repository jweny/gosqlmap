package tools

import (
	"log"
	"os"
)

var (
	Info *log.Logger
	Warning *log.Logger
	Error * log.Logger
	Payload * log.Logger
	Debug * log.Logger
	TrafficOut * log.Logger
)

func init(){
	Debug = log.New(os.Stdout,"[DEBUG] ",log.Ldate | log.Ltime | log.Lshortfile)
	TrafficOut = log.New(os.Stdout,"[TRAFFIC OUT] ",log.Ldate | log.Ltime | log.Lshortfile)
	Info = log.New(os.Stdout,"[INFO] ",log.Ldate | log.Ltime | log.Lshortfile)
	Warning = log.New(os.Stdout,"[WARNING] ",log.Ldate | log.Ltime | log.Lshortfile)
	Error = log.New(os.Stderr,"[ERROR] ",log.Ldate | log.Ltime | log.Lshortfile)
	Payload = log.New(os.Stderr,"[PAYLOAD] ",log.Ldate | log.Ltime | log.Lshortfile)
}
