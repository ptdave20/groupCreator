package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shibukawa/configdir"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"encoding/csv"
	"io"
	"bytes"
	"strings"
)

var saveDir string

func HasExistingToken() (*oauth2.Token, error) {
	b, e := ioutil.ReadFile(path.Join(saveDir, "token.json"))
	if e != nil {
		return nil, e
	}
	var token *oauth2.Token
	token = new(oauth2.Token)

	if e := json.Unmarshal(b, token); e != nil {
		return nil, e
	}

	return token, nil
}

func GetClient() (*http.Client, error) {
	ctx := context.Background()
	cfgBytes, err := ioutil.ReadFile(path.Join(saveDir, "config.json"))
	if err != nil {
		return nil, err
	}

	o2, err := google.ConfigFromJSON(cfgBytes, admin.AdminDirectoryGroupScope)

	if err != nil {
		return nil, err
	}

	token, err := HasExistingToken()

	if err != nil {
		url := o2.AuthCodeURL("", oauth2.ApprovalForce)
		open.Start(url)

		fmt.Printf("If a browser window did not open, please visit:   \n%s\n", url)

		var code string
		fmt.Printf("Enter the given code: ")
		if _, err := fmt.Scan(&code); err != nil {
			return nil, err
		}

		token, err := o2.Exchange(ctx, code)
		if err != nil {
			return nil, err
		}

		b, _ := json.MarshalIndent(*token, "", "    ")
		ioutil.WriteFile(path.Join(saveDir, "token.json"), b, 0644)

		return o2.Client(ctx, token), nil
	} else {
		return o2.Client(ctx, token), nil
	}

	return nil, nil
}

func main() {
	args := os.Args

	cfgDir := configdir.New("ptdave", "groupcreator")
	saveDir = cfgDir.LocalPath

	client, err := GetClient()

	if err != nil {
		println(err.Error())
		return
	}

	groupService, err := admin.New(client)

	for i, v := range args[1:] {
		fmt.Printf("Reading file #%d: %s\n",i,v)
		b, e := ioutil.ReadFile(v)
		if e!=nil {
			fmt.Printf("Error opening file %s : %s\n", v,e.Error())
			continue
		}

		r := csv.NewReader(bytes.NewReader(b))

		line:=0
		for {
			record, err := r.Read()

			// Header
			if line == 0 {
				line++
				continue
			}

			if err == io.EOF {
				break
			}
			if err!=nil {
				fmt.Printf("Error reading file: %s\n",err.Error())
				break
			}

			// address, name, member of
			address := record[0]
			name := record[1]
			var memberOf string
			if len(record) > 2 {
				memberOf = record[2]
			}

			// Does the group exist?
			group, err := groupService.Groups.Get(address).Do()
			if err!=nil {
				group = new(admin.Group)
				group.Name = name
				group.Email = address

				group, err = groupService.Groups.Insert(group).Do()
				if err!= nil {
					fmt.Printf("Error creating group %s: %s\n",address, err.Error())
					continue
				}
			}

			if group.Name != name {
				group.Name = name

				group, err = groupService.Groups.Patch(group.Id, group).Do()
				if err!= nil {
					fmt.Printf("Error updating group %s: %s\n",address, err.Error())
					continue
				}
			}

			if len(strings.TrimSpace(memberOf))>0 {
				parentGroups := strings.Split(memberOf,",")
				for _, parent := range parentGroups {
					// Is it a member of the parent group?
					member, err := groupService.Members.Get(parent,address).Do()
					if err != nil {
						// Not a member, add it
						member = new(admin.Member)
						member.Email = address

						member.Role = "MEMBER"

						member, err = groupService.Members.Insert(parent,member).Do()
						if err!= nil {
							fmt.Printf("Error making %s a member of %s: %s\n",address,parent,err.Error())
							continue
						}
					}
				}
			}


			line++
		}
	}
}
