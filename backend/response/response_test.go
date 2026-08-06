package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func decode(t *testing.T, body []byte) Response {
	t.Helper()
	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, body)
	}
	return r
}

func TestSuccessEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Success(c, "ok", gin.H{"hello": "world"})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	r := decode(t, w.Body.Bytes())
	if r.Status != "success" || r.Message != "ok" {
		t.Fatalf("bad status/message: %+v", r)
	}
	if r.Data == nil {
		t.Fatal("data is nil")
	}
}

func TestCreatedEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Created(c, "made", gin.H{"id": 1})
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}
	r := decode(t, w.Body.Bytes())
	if r.Status != "created" || r.Message != "made" {
		t.Fatalf("bad envelope: %+v", r)
	}
}

func TestErrorEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Error(c, http.StatusBadRequest, "bad")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	r := decode(t, w.Body.Bytes())
	if r.Status != "error" || r.Message != "bad" {
		t.Fatalf("bad envelope: %+v", r)
	}

	if r.Data != nil {
		t.Fatalf("data should be nil, got %+v", r.Data)
	}
}
