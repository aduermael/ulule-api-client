package ulule

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

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
func (c *Client) GetProjects(userID int, filter ProjectFilter) ([]*Project, error) {
	userIDStr := strconv.Itoa(userID)

	resp := &ListProjectResponse{}
	err := c.apiget("/users/"+userIDStr+"/projects?state="+string(filter), resp)
	if err != nil {
		return nil, err
	}
	return resp.Projects, nil
}

// GetProject returns one specific ClientAPI user's
// project identified by its Id or Slug.
func (c *Client) GetProject(projectID int) (*Project, error) {
	projectIDStr := strconv.Itoa(projectID)

	project := &Project{}
	err := c.apiget("/projects/"+projectIDStr, project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// GetProjectSupporters lists supporters for a project
// limit and offset stand for pagination
// the boolean returns indicates if it was the last
// page of supporters or not.
func (c *Client) GetProjectSupporters(projectID, limit, offset int) ([]*Supporter, error, bool) {
	projectIDStr := strconv.Itoa(projectID)
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)

	supporters := &ListSupporterResponse{}
	err := c.apiget("/projects/"+projectIDStr+"/supporters?limit="+limitStr+"&offset="+offsetStr, supporters)
	if err != nil {
		return nil, err, false
	}

	return supporters.Supporters, nil, supporters.Meta.Next == ""
}

// GetProjectOrders lists orders for a project
// limit and offset stand for pagination
// the boolean returns indicates if it was the last
// page or not.
func (c *Client) GetProjectOrders(projectID, limit, offset int) ([]*Order, error, bool) {
	projectIDStr := strconv.Itoa(projectID)
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)

	orders := &ListOrderResponse{}
	err := c.apiget("/projects/"+projectIDStr+"/orders?limit="+limitStr+"&offset="+offsetStr, orders)
	if err != nil {
		return nil, err, false
	}

	return orders.Orders, nil, orders.Meta.Next == ""
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
