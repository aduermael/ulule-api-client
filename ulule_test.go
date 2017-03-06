package ulule

import (
	"fmt"
	"os"
	"testing"
)

var (
	client *Client
)

func TestMain(m *testing.M) {
	// initialize client
	username := os.Getenv("USERNAME")
	userid := os.Getenv("USERID")
	apikey := os.Getenv("APIKEY")
	client = New(username, userid, apikey)

	// run tests
	os.Exit(m.Run())
}

func TestGetProjects(t *testing.T) {
	fmt.Println("list created projects")
	projects, err := client.GetProjects(ProjectFilterCreated)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("found %d projects\n", len(projects))

	fmt.Println("list followed projects")
	projects, err = client.GetProjects(ProjectFilterFollowed)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("found %d projects\n", len(projects))

	fmt.Println("list supported projects")
	projects, err = client.GetProjects(ProjectFilterSupported)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("found %d projects\n", len(projects))

	fmt.Println("list all projects")
	projects, err = client.GetProjects(ProjectFilterAll)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("found %d projects\n", len(projects))
}
