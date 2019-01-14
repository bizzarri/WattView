//******************************************************
//* emissions - read the emissions information from the watttime.org API
//*
//* Written by Maurice Bizzarri, January, 2019
//*
//* Version 0.0 - initial functionaility
//*
//* takes two arguments (one is optional) on command line:
//* -debug  displays debug information
//* -l location uses that location to get info
//* also reads $HOME/.WattTime/ba file and uses that for
//* the balancing authority if nothing on command line
//* defaults to CAISO_ZP26 if nothing supplied or in ba file
//*
//* see watttime.org for an interactive map to figure out
//* your grid designation or use gridregion which creates
//* and rewrites $HOME/.WattTime/ba file
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
	"time"
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

type Response struct {
	Barea      string `json:"ba"`
	Valid      int    `json:"validFor"`
	Validuntil string `json:"validUntil"`
	Rating     string `json:"rating"`
	Green      string `json:"switch"`
	Percent    string `json:"percent"`
	Freq       string `json:"freq"`
	Market     string `json:"market"`
}

func Check(val error, explain string) {
	if val != nil {
		panic(fmt.Sprintf("Error: %s\n%v\n", explain, val))

	}
}

func main() {
	//
	// check for debug option
	// check for balancing authority parameter
	//
	// version number for program
	//
	version := 0.0
	//
	// l - balancing authority
	// debug - debug flag
	//
	var loc string
	var account string
	var password string
	boolPtr := flag.Bool("debug", false, "Debug flag")
	boolPtr2 := flag.Bool("q", false, "Quiet flag")
	flag.StringVar(&loc, "l", "", "Balancing Authority")
	flag.StringVar(&account, "a", "", "Account to use")
	flag.StringVar(&password, "p", "", "Account Password")

	flag.Parse()
	debug := *boolPtr
	quiet := *boolPtr2
	if debug {
		fmt.Printf("Debug flag true - in debug mode.\n")
		fmt.Printf("Version: %1.2f\n", version)
	}
	//
	// get account and password from $HOME/.WattTime/account
	// should be set by makeacct
	//
	defaultloc := "CAISO_ZP26"
	homedir := os.Getenv("HOME")
	acctfile := homedir + "/.WattTime/account"
	bafile := homedir + "/.WattTime/ba"
	//
	// see if ba file created
	// command line takes precedent.
	// if loc is not "g" then check file
	// if file isn't there then default to defaultloc
	//
	var locate string
	if loc != "" {
		locate = loc
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
	//*
	//* if account not specified in command line, look in
	//* account file in $HOME/.WattTime/account
	//*
	if account == "" {
		accts, err := ioutil.ReadFile(acctfile)
		Check(err, "Accounts file not found or other read error")
		var macct MakeAcct
		err = json.Unmarshal(accts, &macct)
		Check(err, "Error unmarshalling accounts files")
		account = macct.Username
		password = macct.Password
	}
	//*
	//* sanity check
	//*
	if account == "" || password == "" {
		fmt.Printf("Error: account and password must be specified.\n")
		os.Exit(-1)
	}

	if debug {
		fmt.Printf("Account Name: %s\n", account)
		fmt.Printf("Password: %s\n", password)
	}
	if !quiet {
		fmt.Printf("WattTime Emissions Real Time Analysis for %s\n\n", locate)
	}
	//
	// had to increase time out
	//
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", "https://api2.watttime.org/v2/login", nil)
	req.SetBasicAuth(account, password)
	resp, err := client.Do(req)
	Check(err, "Error login request call")

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("Error: Status Code: %d\n", resp.StatusCode)
		panic(fmt.Sprintf("Status Error: %s\n", resp.Status))

	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if debug {
		fmt.Printf("body: %s\n", bodyText)
	}
	var wtoken map[string]interface{}
	err = json.Unmarshal(bodyText, &wtoken)
	Check(err, "Error unmarshalling first call for token")
	if debug {
		fmt.Printf("token: %s\n", wtoken["token"])
	}
	request := "https://api2.watttime.org/v2/index/?ba=" + locate + "&latitude=&longitude=&style=all"
	req, err = http.NewRequest("GET", request, nil)
	Check(err, "Error creating NewRequest")
	bearer := "Bearer " + wtoken["token"].(string)
	//*
	//* basic authorization header
	//*
	req.Header.Add("Authorization", bearer)
	resp, err = client.Do(req)
	Check(err, "Error getting NewRequest")
	if resp.StatusCode != 200 {
		fmt.Printf("Error: Status Code: %d\n", resp.StatusCode)
		panic(fmt.Sprintf("Status Error: %s\n", resp.Status))

	}

	response, err := ioutil.ReadAll(resp.Body)
	Check(err, "Error reading response")
	if debug {
		fmt.Printf("Response: %s\n", response)
	}
	//*
	//* Get response and pretty print
	//*
	var emisres Response
	err = json.Unmarshal(response, &emisres)
	Check(err, "Error unmarshalling response")
	if !quiet {
		if emisres.Green == "0" {
			fmt.Printf("Don't switch (not green)\n")
		} else {
			fmt.Printf("Switch! (green grid)\n")
		}

		timed, err := time.Parse(time.RFC3339, emisres.Validuntil)
		Check(err, "Error parsing Valid Until time")
		fmt.Printf("Valid Until: %02d:%02d:%02d UT\n", timed.Hour(), timed.Minute(), timed.Second())
		fmt.Printf("Rating (0=Extremely Clean, 5=Harmful): %s\n", emisres.Rating)
		fmt.Printf("Percent Dirty (0-100): %s\n", emisres.Percent)
	}
	//*
	//* set return
	//* return 0 if green is good, else return 1
	//*
	if emisres.Green == "1" {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
