//******************************************************
//* makeacct - make an account in the WattTime system
//*
//* Written by Maurice Bizzarri, January, 2019
//*
//* Version 0.0 - initial functionaility
//*
//* takes two (or three) arguments on command line:
//* -debug  displays debug information
//* -account account name to create
//* -password if you don't specify password, a random 14 digit password witll
//* be created for you.
//* This software stores the account and password in a text file
//* in ~/home/.watttime/account which is used by the rest of the system
//*
//*
//*****************************************************

package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"flag"
	"bytes"
	"math/rand"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!"

func RandStringBytesRmndr(n int) string {
	b := make([]byte,n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}


type Wtoken struct {
     token string `json:"token"`

}

type Response struct {
	Ok string `json:"ok"`
	User string `json:"user"`
}
	
type MakeAcct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Org      string `json:"org"`
	
}
func Check(val error, explain string)  {
	if val != nil {
		fmt.Printf("Error: %s\n%v\n",explain,val)
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
	flag.StringVar(&account, "a","default","Account name to create")
	flag.StringVar(&password, "p", "", "Password to use")
	flag.StringVar(&email, "e", "", "Email to use")
	flag.StringVar(&org, "o", "", "(optional) organization name")

	flag.Parse()
	debug := *boolPtr
	if debug {
		fmt.Printf("Debug flag true - in debug mode.\n")
		fmt.Printf("Version: %1.2f\n", version)
	}

	fmt.Printf("WattTime Account Creation - create account: %s\n",account)
	if password == "" {
		fmt.Printf("Password not specified, will create 14 character random password\n")

                password = RandStringBytesRmndr(14)
		fmt.Printf("password: %s\n",password)
	}

     timeout := time.Duration(5 * time.Second)
     client := &http.Client{
     	    Timeout: timeout,
	    }
        var regvals MakeAcct
	regvals.Username = account
	regvals.Password  = password
	regvals.Email = email
	regvals.Org = org
	//
	// make into []byte for NewRequest
	//
	reqbytes := new(bytes.Buffer)
	json.NewEncoder(reqbytes).Encode(regvals)
        req, err := http.NewRequest("POST","https://api2.watttime.org/v2/register",bytes.NewBuffer(reqbytes.Bytes()))
	req.Header.Set("Content-Type","application/json")
	//
	// don't actually do anything if debugging for the moment
	//
        var bodyText []byte
        if !debug {
	resp,err := client.Do(req)
	Check (err,"Error Account request call")
        defer resp.Body.Close()
        bodyText, err = ioutil.ReadAll(resp.Body)
	Check (err,"Error reading returned json")
	} else {
                m := Response{"User Created!", account}
                bodyText, err = json.Marshal(m)
                Check(err,"Error marshalling")

	   fmt.Printf("returned body: %s\n",bodyText)
	}


     var respdata Response
     err = json.Unmarshal(bodyText,&respdata)
	Check (err,"Error unmarshalling first call for token")

        fmt.Printf("confirmation: %s\n",respdata.Ok)
     	fmt.Printf("confirm account: %s\n",respdata.User)
        err = ioutil.WriteFile(acctfile, reqbytes.Bytes(),0755)
        Check (err, "Error writing account file")



}