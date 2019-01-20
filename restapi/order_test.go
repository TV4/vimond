package restapi

import (
	"strings"
	"testing"
	"time"
)

func TestParseOrder(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const orderJSON = `
		{
			"accessEndDate":    946782245234,
			"endDate":          946782245345,
			"id":               123456,
			"productName":      "foo-product-name",
			"productPaymentID": 234567,
			"startDate":        946782245123,
			"userId":           345678
		}
	`
		order, err := parseOrder(strings.NewReader(orderJSON))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantOrder := &Order{
			AccessEndDate:    time.Date(2000, 1, 2, 3, 4, 5, 234000000, time.UTC),
			EndDate:          time.Date(2000, 1, 2, 3, 4, 5, 345000000, time.UTC),
			ID:               "123456",
			ProductName:      "foo-product-name",
			ProductPaymentID: "234567",
			StartDate:        time.Date(2000, 1, 2, 3, 4, 5, 123000000, time.UTC),
			UserID:           "345678",
		}

		if got, want := order.AccessEndDate, wantOrder.AccessEndDate; !got.Equal(want) {
			t.Errorf("order.AccessEndDate = %s, want %s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
		}

		if got, want := order.EndDate, wantOrder.EndDate; !got.Equal(want) {
			t.Errorf("order.EndDate = %s, want %s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
		}

		if got, want := order.ID, wantOrder.ID; got != want {
			t.Errorf("order.ID = %q, want %q", got, want)
		}

		if got, want := order.ProductName, wantOrder.ProductName; got != want {
			t.Errorf("order.ProductName = %q, want %q", got, want)
		}

		if got, want := order.ProductPaymentID, wantOrder.ProductPaymentID; got != want {
			t.Errorf("order.ProductPaymentID = %q, want %q", got, want)
		}

		if got, want := order.StartDate, wantOrder.StartDate; !got.Equal(want) {
			t.Errorf("order.StartDate = %s, want %s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
		}

		if got, want := order.UserID, wantOrder.UserID; got != want {
			t.Errorf("order.UserID = %q, want %q", got, want)
		}
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		_, err := parseOrder(strings.NewReader("not-json"))

		if err == nil {
			t.Fatal("err is nil")
		}
	})
}

func TestParseOrders(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const orderJSON = `
		[
			{
				"accessEndDate":    946782245234,
				"endDate":          946782245345,
				"id":               123456,
				"productName":      "foo-product-name",
				"productPaymentID": 234567,
				"startDate":        946782245123,
				"userId":           345678
			},
			{
				"accessEndDate":    946782245456,
				"endDate":          946782245567,
				"id":               234567,
				"productName":      "bar-product-name",
				"productPaymentID": 345678,
				"startDate":        946782245234,
				"userId":           456789
			}
		]
	`
		orders, err := parseOrders(strings.NewReader(orderJSON))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantOrders := []*Order{
			{
				AccessEndDate:    time.Date(2000, 1, 2, 3, 4, 5, 234000000, time.UTC),
				EndDate:          time.Date(2000, 1, 2, 3, 4, 5, 345000000, time.UTC),
				ID:               "123456",
				ProductName:      "foo-product-name",
				ProductPaymentID: "234567",
				StartDate:        time.Date(2000, 1, 2, 3, 4, 5, 123000000, time.UTC),
				UserID:           "345678",
			},
			{
				AccessEndDate:    time.Date(2000, 1, 2, 3, 4, 5, 456000000, time.UTC),
				EndDate:          time.Date(2000, 1, 2, 3, 4, 5, 567000000, time.UTC),
				ID:               "234567",
				ProductName:      "bar-product-name",
				ProductPaymentID: "345678",
				StartDate:        time.Date(2000, 1, 2, 3, 4, 5, 234000000, time.UTC),
				UserID:           "456789",
			},
		}

		if got, want := len(orders), len(wantOrders); got != want {
			t.Fatalf("got %d orders, want %d", got, want)
		}

		for n := range wantOrders {
			order := orders[n]
			wantOrder := wantOrders[n]

			if got, want := order.AccessEndDate, wantOrder.AccessEndDate; !got.Equal(want) {
				t.Errorf("order[%d].AccessEndDate = %s, want %s", n, got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
			}

			if got, want := order.EndDate, wantOrder.EndDate; !got.Equal(want) {
				t.Errorf("order[%d].EndDate = %s, want %s", n, got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
			}

			if got, want := order.ID, wantOrder.ID; got != want {
				t.Errorf("order[%d].ID = %q, want %q", n, got, want)
			}

			if got, want := order.ProductName, wantOrder.ProductName; got != want {
				t.Errorf("order[%d].ProductName = %q, want %q", n, got, want)
			}

			if got, want := order.ProductPaymentID, wantOrder.ProductPaymentID; got != want {
				t.Errorf("order[%d].ProductPaymentID = %q, want %q", n, got, want)
			}

			if got, want := order.StartDate, wantOrder.StartDate; !got.Equal(want) {
				t.Errorf("order[%d].StartDate = %s, want %s", n, got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
			}

			if got, want := order.UserID, wantOrder.UserID; got != want {
				t.Errorf("order[%d].UserID = %q, want %q", n, got, want)
			}
		}
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		_, err := parseOrder(strings.NewReader("not-json"))

		if err == nil {
			t.Fatal("err is nil")
		}
	})
}
