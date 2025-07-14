package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func decode(t *testing.T, r io.Reader, v interface{}) {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		t.Fatalf("decode error: %v", err)
	}
}

func TestOrderWorkflow(t *testing.T) {
	// 1) Create a product
	prodPayload := map[string]interface{}{
		"name":        "IntegrationProduct",
		"description": "for orders",
		"price":       9.99,
	}
	buf, _ := json.Marshal(prodPayload)
	resp, err := http.Post("http://localhost:8080/products", "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("product create failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("product create status = %d; want %d", resp.StatusCode, http.StatusCreated)
	}
	var prod map[string]interface{}
	decode(t, resp.Body, &prod)
	prodID := int(prod["id"].(float64))

	// 2) Sign up a user
	userPayload := map[string]string{"email": "order@int.test", "password": "secret"}
	buf, _ = json.Marshal(userPayload)
	resp, err = http.Post("http://localhost:8080/users/signup", "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("signup failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("signup status = %d; want %d", resp.StatusCode, http.StatusCreated)
	}

	// 3) Log in and capture JWT
	resp, err = http.Post("http://localhost:8080/users/login", "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d; want %d", resp.StatusCode, http.StatusOK)
	}
	var login struct {
		Token string `json:"token"`
	}
	decode(t, resp.Body, &login)
	token := login.Token
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	client := &http.Client{}

	// 4) Create a cart
	req, _ := http.NewRequest("POST", "http://localhost:8080/carts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("create cart failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("cart status = %d; want %d", resp.StatusCode, http.StatusCreated)
	}
	var cart map[string]interface{}
	decode(t, resp.Body, &cart)
	cartID := int(cart["id"].(float64))

	// 5) Add item to cart
	itemPayload := map[string]interface{}{"product_id": prodID, "quantity": 2}
	buf, _ = json.Marshal(itemPayload)
	req, _ = http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/carts/%d/items", cartID), bytes.NewReader(buf))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("add item failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("add item status = %d; want %d", resp.StatusCode, http.StatusNoContent)
	}

	// 6) Place an order
	orderPayload := map[string]int{"cart_id": cartID}
	buf, _ = json.Marshal(orderPayload)
	req, _ = http.NewRequest("POST", "http://localhost:8080/orders", bytes.NewReader(buf))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("order status = %d; want %d", resp.StatusCode, http.StatusCreated)
	}
	var order map[string]interface{}
	decode(t, resp.Body, &order)
	orderID := int(order["id"].(float64))

	// 7) Get the order
	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/orders/%d", orderID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("get order failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get order status = %d; want %d", resp.StatusCode, http.StatusOK)
	}

	// 8) List orders
	req, _ = http.NewRequest("GET", "http://localhost:8080/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("list orders failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list orders status = %d; want %d", resp.StatusCode, http.StatusOK)
	}
}
