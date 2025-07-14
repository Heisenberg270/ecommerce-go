package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestProductWorkflow(t *testing.T) {
	// 1) Create
	payload := map[string]interface{}{
		"name":        "IntProd",
		"description": "from integration test",
		"price":       3.21,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:8080/products", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}
	var created map[string]interface{}
	decode(t, resp.Body, &created)
	id := int(created["id"].(float64))

	// 2) Get
	resp, err = http.Get(fmt.Sprintf("http://localhost:8080/products/%d", id))
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	// 3) Delete
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/products/%d", id), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", resp.StatusCode)
	}
}
