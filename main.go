package main 

import
(
	"fmt"
	"flag"
)

func main() {
	var foo string
	var htmlString string
	
	// Get command line parameter for html file
	flag.StringVar(&foo, "foo", "", "TODO: delete this")
	flag.Parse()

	if (foo == "") {
		flag.PrintDefaults()
		return
	}

	// Read html file to string
	fmt.Println("Welcome to main.go")
}