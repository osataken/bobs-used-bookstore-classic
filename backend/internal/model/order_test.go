package model

import "testing"

func TestOrderSubTotal(t *testing.T) {
	// SubTotal is sum of Book.Price per OrderItem (NOT multiplied by Quantity - preserving source bug)
	order := &Order{
		OrderItems: []OrderItem{
			{Book: Book{Price: 10.00}, Quantity: 2},
			{Book: Book{Price: 5.00}, Quantity: 3},
		},
	}

	subtotal := order.SubTotal()
	expected := 15.00 // 10 + 5 (not 10*2 + 5*3)
	if subtotal != expected {
		t.Errorf("expected subtotal %.2f, got %.2f", expected, subtotal)
	}
}

func TestOrderTax(t *testing.T) {
	order := &Order{
		OrderItems: []OrderItem{
			{Book: Book{Price: 100.00}, Quantity: 1},
		},
	}

	tax := order.Tax()
	expected := 10.00 // 10% of 100
	if tax != expected {
		t.Errorf("expected tax %.2f, got %.2f", expected, tax)
	}
}

func TestOrderTotal(t *testing.T) {
	order := &Order{
		OrderItems: []OrderItem{
			{Book: Book{Price: 100.00}, Quantity: 1},
		},
	}

	total := order.Total()
	expected := 110.00 // 100 + 10
	if total != expected {
		t.Errorf("expected total %.2f, got %.2f", expected, total)
	}
}

func TestOrderStatusString(t *testing.T) {
	tests := []struct {
		status   OrderStatus
		expected string
	}{
		{OrderStatusPending, "Pending"},
		{OrderStatusOrdered, "Ordered"},
		{OrderStatusShipped, "Shipped"},
		{OrderStatusDelivered, "Delivered"},
		{OrderStatusCancelled, "Cancelled"},
	}

	for _, tt := range tests {
		if tt.status.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.status.String())
		}
	}
}
