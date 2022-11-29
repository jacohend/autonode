package autonode

const printedLength = 8

type Config struct {
	Seeds []string `long:"seeds" description:"public seed ip"`
	Host  string   `long:"listen" description:"host/port combo to listen in on"`
}
