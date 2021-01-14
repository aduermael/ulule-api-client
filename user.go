package ulule

import (
	"errors"
	"strconv"
)

// GetUser returns a User for for given ID
func (c *Client) GetUser(userID int) (*User, error) {
	userIDStr := strconv.Itoa(userID)

	resp := &User{}
	err := c.apiget("/users/"+userIDStr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Me returns connected user
func (c *Client) Me() (*User, error) {
	resp := &User{}
	err := c.apiget("/me", resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserOrders lists orders for a user
// limit and offset stand for pagination
// the boolean returns indicates if it was the last
// page or not.
// This function only works for connected user.
func (c *Client) GetUserOrders(u *User, limit, offset int) ([]*Order, error, bool) {
	if u == nil {
		return nil, errors.New("user can't be nil"), false
	}

	userIDStr := strconv.Itoa(u.ID)
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)

	orders := &ListOrderResponse{}
	err := c.apiget("/users/"+userIDStr+"/orders?limit="+limitStr+"&offset="+offsetStr, orders)
	if err != nil {
		return nil, err, false
	}

	return orders.Orders, nil, orders.Meta.Next == ""
}
