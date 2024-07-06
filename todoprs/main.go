package main

import (
	"fmt"
	"log"
	"os/exec"
)

// TODO
func main() {
	// Need the project path.
	// argument to check which method to extract the TODOs. (Check if rg is available, prio it if it is and the user didn't specify a method)
	// argument for the output method (json, xml, csv)
	// for now we start with rg -> json.
	path := "/Users/medkhalil/dev/Go/htmlParsingProject/htmlprs"
	cmd := exec.Command("rg", "-nw", "TODO", path)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The output is %s\n", out)

}
