package flag

import "flag"

var IsTestMode bool // for api test

func init() {
	flag.BoolVar(&IsTestMode, "t", false, "if start in test mode")

	flag.Parse()
}
