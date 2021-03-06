//******************************************************
//* passrecover - recover from a forgotten password
//* this depends on a valid email address being supplied when
//* account was created
//*
//* Written by Maurice Bizzarri, January, 2019
//*
//* Version 0.0 - initial functionaility
//*
//* takes two (one optional) arguments on command line:
//* -debug  displays debug information
//* -a account account name
//*
//*****************************************************

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

)

type Wtoken struct {
	token string `json:"token"`
}

type Datadef struct {
	Abbrev string `json:"abbrev"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
}

func Check(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n%v\n", msg, err))

	}
}
func main() {

	//
	// debug - debug flag
	// account - account to password reset
	//
	version := 0.0
	var account string
	boolPtr := flag.Bool("debug", false, "Debug flag")
	flag.StringVar(&account, "a", "default", "Account to recover")
	flag.Parse()
	debug := *boolPtr
	if debug {
		fmt.Printf("Debug flag true - in debug mode.\n")
		fmt.Printf("Version: %1.2f\n", version)
	}

        if account == "default" {
		panic(fmt.Sprintf("Please specify account to recover using -a\n"))

	}
	fmt.Printf("Recover password for account: %s\n", account)
	client := &http.Client{}
	request := "https://api2.watttime.org/v2/password/?username="
	request = request + account
	req, err := http.NewRequest("GET", request, nil)
	Check(err, "Error creating NewRequest")
	resp, err := client.Do(req)
        defer resp.Body.Close()
	Check(err, "Error GET Request")
        if resp.StatusCode != 200 {
		fmt.Printf("Error: Status Code: %d\n",resp.StatusCode)
		panic(fmt.Sprintf("Status Error: %s\n",resp.Status))

	}

	response, err := ioutil.ReadAll(resp.Body)
	Check(err, "Error reading response")
	if debug {
		fmt.Printf("Response: %s\n", response)
	}
	var answer map[string]interface{}
	err = json.Unmarshal(response, &answer)
	fmt.Printf("\n%s\n", answer["ok"])
	fmt.Printf("Use the link in the email to reset your password.\n\n")

}
