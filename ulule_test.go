package ulule

import (
	"fmt"
	"os"
	"testing"
)

const (
	projectID = 31458  // comme-convenu-2
	userID    = 241660 // bloglaurel
)

var (
	client *Client
)

func TestMain(m *testing.M) {
	// initialize client
	username := os.Getenv("USERNAME")
	apikey := os.Getenv("APIKEY")
	accessToken := os.Getenv("ACCESS_TOKEN")

	if username != "" && apikey != "" {
		client = ClientWithUsernameAndApiKey(username, apikey)
	} else if accessToken != "" {
		client = ClientWithToken(accessToken)
	} else {
		fmt.Fprintf(os.Stderr, "USERNAME & APIKEY or ACCESS_TOKEN required")
		os.Exit(1)
	}

	// run tests
	os.Exit(m.Run())
}

func TestGetProjects(t *testing.T) {
	fmt.Println("list projects")
	_, err := client.GetProjects(userID, ProjectFilterAll)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOneProject(t *testing.T) {
	fmt.Println("get project")
	_, err := client.GetProject(projectID)
	if err != nil {
		t.Error(err)
	}
}

func TestGetProjectSupporters(t *testing.T) {
	fmt.Println("get project supporters")
	_, err, _ := client.GetProjectSupporters(projectID, 20, 0)
	if err != nil {
		t.Error(err)
	}
}

func TestGetProjectOrders(t *testing.T) {
	fmt.Println("get project orders")
	_, err, _ := client.GetProjectOrders(projectID, 20, 0)
	if err != nil {
		t.Error(err)
	}
}

// Users

func TestGetUser(t *testing.T) {
	fmt.Println("get user")
	usr, err := client.GetUser(userID)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", usr)
}
