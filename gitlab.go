package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	host  *string
	token *string
	key   *string
)

type gitlabkeys struct {
	Id   int
	Name string
}

func gitlab(url string, datastring string, token string) []gitlabkeys {
	var (
		req    *http.Request
		resp   *http.Response
		err    error
		result []gitlabkeys
	)
	if datastring != "" {
		var jsonStr = []byte(`{"name":` + `"` + datastring + `"` + `}`)
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	} else
	{
		var jsonStr = []byte(``)
		req, err = http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	}

	req.Header.Set("Private-Token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &result)

	return result
}

func getarg() bool {
	if len(os.Args) == 1 {
		fmt.Println("usage: %s Command [<args>]", os.Args[0])
		fmt.Printf("<<Command>>\n")
		fmt.Printf("list:\n-host \"gitlab basic hostname\"\n-token \"Gitlab Access Token\"\n")
		fmt.Printf("search:\n-host \"gitlab basic hostname\"\n-token \"Gitlab Access Token\"\n-key \"key for finding in project names\"\n")
		fmt.Printf("-----\n")
		fmt.Printf("Examples:\n\n./myapp list -host \"gitlab.com\" -token \"myaccesstoken\"\n")
		fmt.Printf("./myapp search -key \"myprojectname\" -host \"gitlab.com\" -token \"myaccesstoken\"\n\n")
		fmt.Println("Have Fun :)")
		return true
	}
	switch os.Args[1] {
	case "list":
		Listprojects()
	case "search":
		Searchprojects()
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)

	}
	return true
}

func Listprojects() bool {
	var projects []gitlabkeys
	listcmd := flag.NewFlagSet("list", flag.ExitOnError)
	host = listcmd.String("host", "", "Enter Gitlab Host")
	token = listcmd.String("token", "", "Enter Gitlab AccessToken")
	listcmd.Parse(os.Args[2:])
	if *host != "" && *token != "" {
		url := string(*host) + "/api/v4/projects"
		projects = gitlab(url, "", *token)
		for value := range projects {
			fmt.Printf("repository with id %d and name :%s\n", projects[value].Id, projects[value].Name)
		}
	} else {
		listcmd.PrintDefaults()
	}
	return true
}
func Searchprojects() string {
	var projects []gitlabkeys
	searchcmd := flag.NewFlagSet("search", flag.ExitOnError)
	host = searchcmd.String("host", "", "Enter Gitlab Host")
	token = searchcmd.String("token", "", "Enter Gitlab AccessToken")
	key = searchcmd.String("key", "", "Enter what you'r looking for ;)")
	searchcmd.Parse(os.Args[2:])
	if *key != "" && *host != "" && *token != "" {

		url := string(*host) + "/api/v4/projects"
		projects = gitlab(url, "", *token)
		for value := range projects {
			if strings.Contains(*key, projects[value].Name) == true {
				fmt.Printf("we found this one : id:%d\tname:%s\n", projects[value].Id, projects[value].Name)
			}
		}
	} else {
		searchcmd.PrintDefaults()
	}

	return ""
}

func main() {
	getarg()

}
