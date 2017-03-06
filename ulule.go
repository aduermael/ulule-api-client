package ulule

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// ClientAPI is a structure that can be used to
// communicate with Ulule's API
type Client struct {
	username string
	// TODO: get userid from username, if that's possible
	userid     string
	apikey     string
	httpClient *http.Client
}

// New returns a Client initialized with given credentials
func New(username, userid, apikey string) *Client {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{Transport: transport}

	clientAPI := &Client{
		username:   username,
		userid:     userid,
		apikey:     apikey,
		httpClient: httpClient,
	}

	return clientAPI
}

// ProjectFilter defines filters that can be used when listing projects
type ProjectFilter string

const (
	ProjectFilterCreated   ProjectFilter = "created"
	ProjectFilterFollowed  ProjectFilter = "followed"
	ProjectFilterSupported ProjectFilter = "supported"
	ProjectFilterAll       ProjectFilter = ""
)

// GetProjects returns ClientAPI user's projects.
// Supported string filters: "created", "followed", "supported", "" (no filter)
func (c *Client) GetProjects(filter ProjectFilter) ([]*Project, error) {
	req, err := http.NewRequest("GET", "https://api.ulule.com/v1/users/"+c.username+"/projects?state="+string(filter), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "ApiKey "+c.username+":"+c.apikey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	listProjectResp := &ListProjectResponse{}
	bodyBytes, err := rawHTMLBody(resp)

	if resp.StatusCode != 200 {
		if err != nil {
			return nil, fmt.Errorf("%d : %s", resp.StatusCode, err.Error())
		}
		return nil, fmt.Errorf("%d : %s", resp.StatusCode, string(bodyBytes))
	}

	err = json.Unmarshal(bodyBytes, listProjectResp)
	if err != nil {
		return nil, err
	}

	return listProjectResp.Projects, nil
}

// GetProject returns one specific ClientAPI user's
// project identified by its Id or Slug.
func (c *Client) GetProject(identifier string) (*Project, error) {
	identifier = strings.Trim(identifier, " ")
	projects, err := c.GetProjects("")
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if identifier == project.Slug || identifier == strconv.FormatFloat(project.ID, 'f', 0, 64) {
			return project, nil
		}
	}
	return nil, errors.New("error: project not found (" + identifier + ")")
}

// GetProjectSupporters lists supporters for a project
// limit and offset stand for pagination
// the boolean returns indicates if it was the last
// page of supporters or not.
func (c *Client) GetProjectSupporters(projectID, limit, offset int) ([]*Supporter, error, bool) {

	projectIDStr := strconv.Itoa(projectID)

	req, err := http.NewRequest("GET", "https://api.ulule.com/v1/projects/"+projectIDStr+"/supporters?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil)
	if err != nil {
		return nil, err, false
	}

	req.Header.Add("Authorization", "ApiKey "+c.username+":"+c.apikey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err, false
	}

	listSupporterResp := &ListSupporterResponse{}
	decodeHTMLBody(resp, listSupporterResp)
	lastPage := listSupporterResp.Meta.Next == ""

	return listSupporterResp.Supporters, nil, lastPage
}

// GetProjectOrders lists orders for a project
// limit and offset stand for pagination
// the boolean returns indicates if it was the last
// page or not.
func (c *Client) GetProjectOrders(projectID, limit, offset int) ([]*Order, error, bool) {

	projectIDStr := strconv.Itoa(projectID)

	req, err := http.NewRequest("GET", "https://api.ulule.com/v1/projects/"+projectIDStr+"/orders?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil)
	if err != nil {
		return nil, err, false
	}

	req.Header.Add("Authorization", "ApiKey "+c.username+":"+c.apikey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err, false
	}

	listOrderResp := &ListOrderResponse{}
	decodeHTMLBody(resp, listOrderResp)
	lastPage := listOrderResp.Meta.Next == ""

	return listOrderResp.Orders, nil, lastPage
}

// HTML utils

func decodeHTMLBody(response *http.Response, i interface{}) error {
	decoder := json.NewDecoder(response.Body)
	for {
		err := decoder.Decode(i)
		if err != nil && err != io.EOF {
			return err
		}
		if err != nil && err == io.EOF {
			break
		}
	}
	return nil
}

func rawHTMLBody(response *http.Response) ([]byte, error) {
	resp := make([]byte, 0)
	p := make([]byte, 64)
	for {
		n, err := response.Body.Read(p)
		if err != nil && err != io.EOF {
			return resp, err
		}
		resp = append(resp, p[:n]...)
		if err != nil {
			break
		}
	}
	return resp, nil
}
