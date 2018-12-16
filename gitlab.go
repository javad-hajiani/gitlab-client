package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	host  *string
	token *string
	key   *string
)

type gitlabkeys struct {
	Id    int
	Name  string
	Title string
}

func gitlab(url string, datastring map[string]interface{}, token string) []gitlabkeys {
	var (
		req    *http.Request
		resp   *http.Response
		err    error
		result []gitlabkeys
	)
	if datastring != nil {
		var jsonStr, err = json.Marshal(datastring)
		if err != nil {
			log.Fatalln(err)
		}

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	} else {
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
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		panic(err)
	}
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
	case "deploykey":
		Enable_deploy_key()
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
	err := listcmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	if *host != "" && *token != "" {
		url := string(*host) + "/api/v4/projects"
		projects = gitlab(url, nil, *token)
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
	err := searchcmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	if *key != "" && *host != "" && *token != "" {

		url := string(*host) + "/api/v4/projects"
		projects = gitlab(url, nil, *token)
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

func Enable_deploy_key() string {
	var projects []gitlabkeys
	var key_id int
	deploycmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	host = deploycmd.String("host", "", "Enter Gitlab Host")
	token = deploycmd.String("token", "", "Enter Gitlab AccessToken")
	key = deploycmd.String("deploykey", "", "Enter DeployKey ;)")
	group := deploycmd.String("group", "", "Enter group for enable ;)")
	err := deploycmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	if *key != "" && *host != "" && *token != "" && *group != "" {
		url := string(*host) + "/api/v4/groups/" + *group + "/projects?page=2&per_page=100"
		deploykeysurl := string(*host) + "/api/v4/deploy_keys"
		deploykeys := gitlab(deploykeysurl, nil, *token)
		projects = gitlab(url, nil, *token)
		for deploykey := range deploykeys {
			if strings.Contains(*key, deploykeys[deploykey].Title) {
				key_id = deploykeys[deploykey].Id
			}
		}
		if key_id != 0 {

			for value := range projects {
				jsondata := map[string]interface{}{"key": strconv.Itoa(projects[value].Id), "key_id": strconv.Itoa(key_id)}
				enable_key_url := string(*host) + "/api/v4/projects/" + strconv.Itoa(projects[value].Id) + "/deploy_keys/" + strconv.Itoa(key_id) + "/enable"
				_ = gitlab(enable_key_url, jsondata, *token)
				fmt.Printf("we found this one : id:%d\tname:%s\n", projects[value].Id, projects[value].Name)
			}
		}
	} else {
		deploycmd.PrintDefaults()
	}

	return ""
}
func main() {
	getarg()

}
