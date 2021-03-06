package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Order holds a subset of the fields of an order.
type Order struct {
	AccessEndDate    time.Time
	EndDate          time.Time
	ID               string
	ProductName      string
	ProductPaymentID string
	StartDate        time.Time
	UserID           string
}

// Order returns information about an order.
func (c *Client) Order(ctx context.Context, platform, orderID string) (*Order, error) {
	path := fmt.Sprintf("/api/%s/order/%s", platform, orderID)

	resp, err := c.get(ctx, path, url.Values{})
	if err != nil {
		return nil, err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 64)
		resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, ErrUnknown
	}

	return parseOrder(resp.Body)
}

// CurrentOrders returns information about a user's currently active orders.
func (c *Client) CurrentOrders(ctx context.Context, platform, userID string) ([]*Order, error) {
	path := fmt.Sprintf("/api/%s/user/%s/orders/current", platform, userID)

	resp, err := c.get(ctx, path, url.Values{})
	if err != nil {
		return nil, err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 64)
		resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, ErrUnknown
	}

	return parseOrders(resp.Body)
}

// CreateOrder creates an order.
func (c *Client) CreateOrder(ctx context.Context, platform, userID, productPaymentID string) (*Order, error) {
	path := fmt.Sprintf("/api/%s/order/%s/create", platform, userID)

	body, err := json.Marshal(struct {
		ProductPaymentID string `json:"productPaymentId"`
	}{
		ProductPaymentID: productPaymentID,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.post(ctx, path, url.Values{}, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 64)
		resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, ErrUnknown
	}

	return parseOrder(resp.Body)
}

// SetOrderEndDates updates the endData and accessEndDate fields of an order to
// the given end date. This method fetches the given order, strips null values
// and nested objects (the Vimond API explodes on them), sets the dates, and
// PUTs the resulting object back. This may result in data loss. Use with
// caution.
func (c *Client) SetOrderEndDates(ctx context.Context, platform, orderID string, endDate time.Time) (*Order, error) {
	return c.updateOrder(ctx, platform, orderID, map[string]interface{}{
		"accessEndDate": endDate.Unix() * 1000,
		"endDate":       endDate.Unix() * 1000,
	})
}

// updateOrder updates an order by overwriting the given field values. This
// method skips nested objects.
func (c *Client) updateOrder(ctx context.Context, platform, orderID string, values map[string]interface{}) (*Order, error) {
	getRawOrder := func(ctx context.Context, platform, orderID string) (map[string]interface{}, error) {
		path := fmt.Sprintf("/api/%s/order/%s", platform, orderID)

		resp, err := c.get(ctx, path, url.Values{})
		if err != nil {
			return nil, err
		}
		defer func() {
			io.CopyN(ioutil.Discard, resp.Body, 64)
			resp.Body.Close()
		}()

		switch resp.StatusCode {
		case http.StatusOK:
			break
		case http.StatusNotFound:
			return nil, ErrNotFound
		default:
			return nil, ErrUnknown
		}

		var m map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			return nil, err
		}

		return m, nil
	}

	rawOrder, err := getRawOrder(ctx, platform, orderID)
	if err != nil {
		return nil, err
	}

	// The Vimond API seems to explode on null values and nested objects, so
	// remove these.
	for k, v := range rawOrder {
		if v == nil {
			delete(rawOrder, k)
			continue
		}
		if _, ok := v.(map[string]interface{}); ok {
			delete(rawOrder, k)
		}
	}

	for k, v := range values {
		rawOrder[k] = v
	}

	body, err := json.Marshal(&rawOrder)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/%s/order/%s", platform, orderID)

	resp, err := c.put(ctx, path, url.Values{}, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 64)
		resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, ErrUnknown
	}

	return parseOrder(resp.Body)
}

func parseOrder(r io.Reader) (*Order, error) {
	var o struct {
		AccessEndDate    time.Time `json:"accessEndDate"`
		EndDate          time.Time `json:"endDate"`
		ID               int       `json:"id"`
		ProductName      string    `json:"productName"`
		ProductPaymentID int       `json:"productPaymentID"`
		StartDate        time.Time `json:"startDate"`
		UserID           int       `json:"userId"`
	}

	if err := json.NewDecoder(r).Decode(&o); err != nil {
		return nil, err
	}

	return &Order{
		AccessEndDate:    o.AccessEndDate,
		EndDate:          o.EndDate,
		ID:               strconv.Itoa(o.ID),
		ProductName:      o.ProductName,
		ProductPaymentID: strconv.Itoa(o.ProductPaymentID),
		StartDate:        o.StartDate,
		UserID:           strconv.Itoa(o.UserID),
	}, nil
}

func parseOrders(r io.Reader) ([]*Order, error) {
	var resp []struct {
		AccessEndDate    time.Time `json:"accessEndDate"`
		EndDate          time.Time `json:"endDate"`
		ID               int       `json:"id"`
		ProductName      string    `json:"productName"`
		ProductPaymentID int       `json:"productPaymentID"`
		StartDate        time.Time `json:"startDate"`
		UserID           int       `json:"userId"`
	}

	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return nil, err
	}

	orders := make([]*Order, 0, len(resp))

	for n := range resp {
		orders = append(orders, &Order{
			AccessEndDate:    resp[n].AccessEndDate,
			EndDate:          resp[n].EndDate,
			ID:               strconv.Itoa(resp[n].ID),
			ProductName:      resp[n].ProductName,
			ProductPaymentID: strconv.Itoa(resp[n].ProductPaymentID),
			StartDate:        resp[n].StartDate,
			UserID:           strconv.Itoa(resp[n].UserID),
		})
	}

	return orders, nil
}
