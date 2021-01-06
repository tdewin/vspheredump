/*
Copyright (c) 2020 tdewin

Original :
Copyright (c) 2017 VMware, Inc. All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/


package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"errors"
	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"time"
)


// getEnvString returns string from environment variable.
func getEnvString(v string, def string) string {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	return r
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(v string, def bool) bool {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	switch strings.ToLower(r[0:1]) {
	case "t", "y", "1":
		return true
	}

	return false
}

const (
	envURL      = "GOVMOMI_URL"
	envUserName = "GOVMOMI_USERNAME"
	envPassword = "GOVMOMI_PASSWORD"
	envInsecure = "GOVMOMI_INSECURE"
)


/*
	vSphere Client/API creation 
*/

//Used by NewClient if the username and password are not in the url
func processOverride(u *url.URL,username string,password string) {
	pwd := password
	if pwd == "" {
		fmt.Print("Password :")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		pwd = string(bytePassword)
		fmt.Println("")
	}
	un := u.User.Username()
	if username != "" {
		un = username
	}
	if un == "" {
		panic("Please provide a user")
	}

	u.User = url.UserPassword(un, pwd)
}

// NewClient creates a vim25.Client for use in the examples
func NewClient(ctx context.Context,url string,insecure bool,username string,password string) (*vim25.Client, error) {
	// Parse URL from string
	u, err := soap.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// Override username and/or password as required
	processOverride(u,username,password)

	// Share govc's session cache
	s := &cache.Session{
		URL:      u,
		Insecure: insecure,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {



	var urlDescription = fmt.Sprintf("ESX or vCenter URL [%s]", envURL)
	var urlFlag = flag.String("url", getEnvString(envURL, ""), urlDescription)
	
	var insecureDescription = fmt.Sprintf("Don't verify the server's certificate chain [%s]", envInsecure)
	var insecureFlag = flag.Bool("insecure", getEnvBool(envInsecure, false), insecureDescription)
	
	var usernameDescription = fmt.Sprintf("Username [%s]", envUserName)
	var usernameFlag = flag.String("username", getEnvString(envUserName, ""), usernameDescription )
	
	var passwordDescription = fmt.Sprintf("Username [%s]", envPassword)
	var passwordFlag = flag.String("password", getEnvString(envPassword, ""), passwordDescription )
	
	var fileName = fmt.Sprintf("vspheredump_%s.json",time.Now().Format("20060201-150405"))
	var fileFlag = flag.String("file",fileName,fmt.Sprintf("Name of the file to dump to %s",fileName))

	var shadowFileFlag = flag.Bool("shadow", false, "detailed shadow file")

	flag.Parse()

	var err error
	if *urlFlag == "" {
		err = errors.New("Please provide url")
	} else {
		//Background context means there is no timeout
		//and for this command line, we don't have any deadlines so that's fine
		ctx := context.Background()
		c, rerr := NewClient(ctx,*urlFlag,*insecureFlag,*usernameFlag,*passwordFlag)
		if rerr == nil {
			err = VsphereDump(ctx,c,*fileFlag,*shadowFileFlag)
		} else {
			err = rerr
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	
}