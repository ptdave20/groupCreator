package main

import (
	"github.com/skratchdot/open-golang/open"
	"google.golang.org/api/admin/directory/v1"
	"golang.org/x/oauth2/google"
	"github.com/shibukawa/configdir"
	"os"
	"context"
	"fmt"
	"io/ioutil"
	"golang.org/x/oauth2"
	"encoding/json"
	"net/http"
	"path"
)

var saveDir string

func HasExistingToken() (*oauth2.Token,error) {
	b, e :=ioutil.ReadFile(path.Join(saveDir,"token.json"))
	if e != nil {
		return nil, e
	}
	var token *oauth2.Token
	token=new(oauth2.Token)

	if e:=json.Unmarshal(b,token); e!=nil {
		return nil, e
	}

	return token, nil
}

func GetClient() (*http.Client,error) {
	ctx:= context.Background()
	cfgBytes, err:= ioutil.ReadFile(path.Join(saveDir,"config.json"))
	o2, err := google.ConfigFromJSON(cfgBytes, admin.AdminDirectoryGroupScope)

	if err!=nil {
		return nil, err
	}

	token, err:= HasExistingToken()

	if err!=nil {
		url:= o2.AuthCodeURL("", oauth2.ApprovalForce)
		open.Start(url)

		fmt.Printf("If a browser window did not open, please visit:   \n%s\n", url)

		var code string
		fmt.Printf("Enter the given code: ")
		if _, err := fmt.Scan(&code); err!=nil {
			return nil, err
		}

		token, err := o2.Exchange(ctx, code)
		if err!=nil {
			return nil, err
		}

		b, _ := json.MarshalIndent(*token,"", "    ")
		ioutil.WriteFile(path.Join(saveDir,"token.json"), b,0644)

		return o2.Client(ctx, token), nil
	} else {
		return o2.Client(ctx, token), nil
	}

	return nil,nil
}

func main() {
	args:=os.Args

	cfgDir := configdir.New("ptdave","groupcreator")
	saveDir = cfgDir.LocalPath

	client, err := GetClient()

	if err!=nil {
		println(err.Error())
		return
	}

	groupService, err := admin.New(client)

	for _, v := range args[:0] {
		println(v)
		req:=groupService.Groups.Get(v)
		if group, err:= req.Do(); err!=nil {
			println(err.Error())
		} else {
			fmt.Println("Group %s exist", group.Email)
		}
	}
}
