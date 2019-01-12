//******************************************************
//* makeacct - make an account in the WattTime system
//*
//* Written by Maurice Bizzarri, January, 2019
//*
//* Version 0.0 - initial functionaility
//*
//* takes two (or three) arguments on command line:
//* -debug  displays debug information
//* -account (-a) account name to create
//* -password (-p) if you don't specify one, a random 14 digit password witll
//* be created for you.
//* -org (-o)
//* This software stores the account, password, email and org in a file
//* in ~/home/.watttime/account which is used by the rest of the system
//*
//*
//*****************************************************

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!"
//*
//* simple random string creator stolen from
//* somewhere on the internet
//* I did add numbers to the string and the exclam to make it 64 chars long
//*
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

type Wtoken struct {
	token string `json:"token"`
}

type Response struct {
	Ok   string `json:"ok"`
	User string `json:"user"`
}

type MakeAcct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Org      string `json:"org"`
}

//*
//* standard err check function
//*

func Check(val error, explain string) {
	if val != nil {
		fmt.Printf("Error: %s\n%v\n", explain, val)
		os.Exit(-1)
	}
}

func main() {
	//
	// check for debug option
	//
	// version number for program
	//
	// create directory if not created already
	// account will be written in JSON into
	// $HOME/.WattTime/account
	//
	version := 0.0
	homedir := os.Getenv("HOME")
	wattdir := homedir + "/.WattTime"
	acctfile := wattdir + "/account"
	//
	// create directory if not created
	//
	os.MkdirAll(wattdir, os.ModePerm)

	//
	// account (a) - account name to create
	// password (p)  - (optional) password to use, else a 14 letter/number
	//          password will be created for you
	// email (e) - email to use as recovery email account
	// org (o) - organization name (optional)
	// debug - debug flag
	//
	var account string
	var password string
	var email string
	var org string
	boolPtr := flag.Bool("debug", false, "Debug flag")
	flag.StringVar(&account, "a", "default", "Account name to create")
	flag.StringVar(&password, "p", "", "Password to use or one will be made for you")
	flag.StringVar(&email, "e", "", "Email to use")
	flag.StringVar(&org, "o", "", "(optional) organization name")

	flag.Parse()
	debug := *boolPtr
	if debug {
		fmt.Printf("Debug flag true - in debug mode.\n")
		fmt.Printf("Version: %1.2f\n", version)
	}

        if account == "default" {
		fmt.Printf("Please specify account and email - see -h for help\n")
		os.Exit(-1)
	}
	fmt.Printf("WattTime Account Creation - create account: %s\n", account)
	if password == "" {
		fmt.Printf("Password not specified, will create 14 character random password\n")

		password = RandStringBytesRmndr(14)
		fmt.Printf("password: %s\n", password)
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	var regvals MakeAcct
	regvals.Username = account
	regvals.Password = password
	regvals.Email = email
	regvals.Org = org
	//
	// make into []byte for NewRequest
	//
	reqbytes := new(bytes.Buffer)
	json.NewEncoder(reqbytes).Encode(regvals)
	req, err := http.NewRequest("POST", "https://api2.watttime.org/v2/register", bytes.NewBuffer(reqbytes.Bytes()))
	req.Header.Set("Content-Type", "application/json")
	//
	// don't actually do anything if debugging for the moment
	//
	var bodyText []byte
	if !debug {
		resp, err := client.Do(req)
		Check(err, "Error Account request call")
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Printf("Error: Status Code: %d\n",resp.StatusCode)
			fmt.Printf("Status Error: %s\n",resp.Status)
			os.Exit(-1)
		}

		bodyText, err = ioutil.ReadAll(resp.Body)
		Check(err, "Error reading returned json")
	} else {
		m := Response{"User Created!", account}
		bodyText, err = json.Marshal(m)
		Check(err, "Error marshalling")

		fmt.Printf("returned body: %s\n", bodyText)
	}

	var respdata Response
	err = json.Unmarshal(bodyText, &respdata)
	Check(err, "Error unmarshalling first call for token")

	fmt.Printf("Confirmation (should be OK): %s\n", respdata.Ok)
	fmt.Printf("Confirm account: %s\n", respdata.User)
	err = ioutil.WriteFile(acctfile, reqbytes.Bytes(), 0644)
	Check(err, "Error writing account file")

}
