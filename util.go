package main

import (
	"encoding/json"
	"fmt"
)

func isLocalBooking(isLocal bool) bool {
	return isLocal || global.ApiEndpoint == "" || global.RemoteBookingEnabled != true
}

// print the contents of the obj
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
