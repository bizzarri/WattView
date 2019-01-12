//******************************************************
//* griddata - read detailed grid data on a balancing authority
//*
//* Written by Maurice Bizzarri, January, 2019
//*
//* Version 0.0 - initial functionaility
//*
//* takes two to 4 arguments on command line:
//* -debug  displays debug information
//* -l location uses that location to get info
//* -a account to use
//* -p password for account
//* defaults to CAISO_ZP26
//* uses $HOME/.WattTime/account if account not specified
//*
//* if $HOME/.WattTime/ba there uses the file unless
//* overridden on command line with -l
//* see watttime.org for an interactive map to figure out
//* your grid designation or use gridregion with lat/long
//*
//*****************************************************

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Wtoken struct {
	token string `json:"token"`
}

type MakeAcct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Org      string `json:"org"`
}

type Datadef struct {
	Ba         string  `json:"ba"`
	Dtype      string  `json:"datatype"`
	Point_time string  `json:"point_time"`
	Frequency  float64 `json:"frequency"`
	Market     string  `json:"market"`
	Val        float64 `json:"value"`
	Fuel       string  `json:"fuel"`
}

func Check(err error, msg string) {
	if err != nil {
		fmt.Printf("Error: %s\n%v\n", msg, err)
		os.Exit(-1)
	}
}
func main() {

	//
	// location - balancing authority for parameter
	// debug - debug flag
	//
	defaultloc := "CAISO_ZP26"
	version := 0.0
	var location string
	var starttime string
	var endtime string
        var account string
	var password string
	boolPtr := flag.Bool("debug", false, "Debug flag")
	flag.StringVar(&location, "l", "nothing", "Balancing Authority abbrev.")
	flag.StringVar(&starttime, "s", "2019-01-02T00:00:00", "Start Time (RFC3339 format")
	flag.StringVar(&endtime, "e", "2019-01-02T00:00:05", "End Time (RFC3339 format)")
        flag.StringVar(&account,"a","nothing","Account Name")
	flag.StringVar(&password,"p","nothing","Account password")

	flag.Parse()
	debug := *boolPtr
	if debug {
		fmt.Printf("Debug flag true - in debug mode.\n")
		fmt.Printf("Version: %1.2f\n", version)
	}

	//
	// get account and password from $HOME/.WattTime/account
	// should be set by makeacct
	//
	homedir := os.Getenv("HOME")
	acctfile := homedir + "/.WattTime/account"
	bafile := homedir + "/.WattTime/ba"
	if account == "nothing" {
		accts, err := ioutil.ReadFile(acctfile)
		Check(err, "Accounts file not found or other read error")
		var macct MakeAcct
		err = json.Unmarshal(accts, &macct)
		Check(err, "Error unmarshalling accounts files")
		account = macct.Username
		password = macct.Password
	}
	
	if debug {
		fmt.Printf("Account Name: %s\n", account)
		fmt.Printf("Password: %s\n", password)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api2.watttime.org/v2/login", nil)
	req.SetBasicAuth(account, password)
	resp, err := client.Do(req)
	Check(err, "Error WattTime login request")
	defer resp.Body.Close()
        if resp.StatusCode != 200 {
		fmt.Printf("Error: Status Code: %d\n",resp.StatusCode)
		fmt.Printf("Status Error: %s\n",resp.Status)
		os.Exit(-1)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	Check(err, "Error reading body")
	if debug {
		fmt.Printf("body: %s\n", bodyText)
	}
	var wtoken map[string]interface{}
	err = json.Unmarshal(bodyText, &wtoken)
	Check(err, "Error unmarshalling body text from login")
	if debug {
		fmt.Printf("token: %s\n", wtoken["token"])
	}
	//
	// see if ba file created
	// command line takes precedent.
	// if loc is not "nothing" then check file
	// if file isn't there then default to defaultloc
	//
	var locate string
	if location != "nothing" {
		locate = location
	} else {
		blocate, err := ioutil.ReadFile(bafile)

		if err == nil {
			locate = string(blocate)
		} else {
			locate = defaultloc

		}
	}
	if debug {
		fmt.Printf("locate: %s\n", locate)
	}

	fmt.Printf("Grid Data for Balancing Authority  %s\n", locate)
	gridstr := "https://api2.watttime.org/v2/data/?ba="
	gridstr = gridstr + locate
	gridstr = gridstr + "&latitude=&longitude=&starttime="
	gridstr = gridstr + starttime
	gridstr = gridstr + "&endtime=" + endtime
	if debug {
		fmt.Printf("getstring: %s\n", gridstr)
	}
	req, err = http.NewRequest("GET", gridstr, nil)
	Check(err, "Error getting request")
	defer resp.Body.Close()
	bearer := "Bearer " + wtoken["token"].(string)
	req.Header.Add("Authorization", bearer)
	resp, err = client.Do(req)
	Check(err, "Error retrieving data")
        if resp.StatusCode != 200 {
		fmt.Printf("Error: Status Code: %d\n",resp.StatusCode)
		fmt.Printf("Status Error: %s\n",resp.Status)
		os.Exit(-1)
	}
	response, err := ioutil.ReadAll(resp.Body)
	Check(err, "Error reading data from GET")
	if debug {
		fmt.Printf("Response: %s\n", response)
	}

	var unwrap []Datadef
	err = json.Unmarshal(response, &unwrap)
	Check(err, "Error unmarshalling response")
	//
	//        for idx := 0; idx < 5; idx++ {
	for idx := range unwrap {
		fmt.Printf("Balancing Authority: %s\n", unwrap[idx].Ba)
		fmt.Printf("Data Type: %s\n", unwrap[idx].Dtype)
		fmt.Printf("Time Stamp: %s\n", unwrap[idx].Point_time)
		fmt.Printf("Frequency: %f seconds\n", unwrap[idx].Frequency)
		fmt.Printf("Value: %f\n", unwrap[idx].Val)
		fmt.Printf("Electricity Grid Region Market: %s\n", unwrap[idx].Market)
		fmt.Printf("Type of Fuel: %s\n", unwrap[idx].Fuel)
	}

}
