package cmd

import (
	"flag"
	"fmt"
)

var version = "unknown"
var revision = "unknown"

func main() {
	showVersion := false
	clockType := ""
	flag.BoolVar(&showVersion, "v", false, "show application version")
	flag.BoolVar(&showVersion, "version", false, "show application version")
	flag.StringVar(&clockType, "clocktype", "", "in:出勤, out:退勤, go:外出, returned:再入")
	flag.Parse()

	if showVersion {
		fmt.Println(fmt.Sprintf("bugyo-client-go version %s.rev-%s", version, revision))
	} else {

	}
}

func createConfig() {

}
