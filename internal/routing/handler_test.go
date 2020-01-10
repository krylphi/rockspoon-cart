package routing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initRouting() http.Handler {
	repo := InitMockRepo()
	return RouterInit(repo, repo)
}

func Test_Handlers(t *testing.T) {
	type ciResponse struct {
		ID       string `json:"id,omitempty"`
		CartID   string `json:"cart_id,omitempty"`
		Product  string `json:"product,omitempty"`
		Quantity int    `json:"quantity,omitempty"`
	}

	type cartResponse struct {
		ID    string        `json:"id"`
		Items []*ciResponse `json:"items,omitempty"`
	}

	jsonVal := bytes.NewBuffer([]byte(`{"product": "Shoes","quantity": 10}`))

	t.Run("HandleNewCart", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/carts", nil)
		if err != nil {
			t.Fatalf("HandleNewCart() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			if err != nil {
				t.Logf("HandleNewCart() error, while closing request body: %v", err.Error())
			}
		}()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("HandleNewCart() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}

		decoder := json.NewDecoder(res.Body)
		var cartRes cartResponse
		err = decoder.Decode(&cartRes)
		if err != nil {
			t.Fatalf("HandleNewCart() Failed to deserialize response error %v", err)
		}
	})

	t.Run("HandleGetCart", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/carts/1", nil)
		if err != nil {
			t.Fatalf("HandleGetCart() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			if err != nil {
				t.Logf("HandleGetCart() error, while closing request body: %v", err.Error())
			}
		}()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("HandleGetCart() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}

		decoder := json.NewDecoder(res.Body)
		defer func() {
			err := res.Body.Close()
			if err != nil {
				t.Logf("HandleGetCart() error, while closing request body: %v", err.Error())
			}
		}()

		var response cartResponse

		err = decoder.Decode(&response)
		if err != nil {
			t.Fatalf("HandleGetCartFail() Failed to deserialize response error %v", err)
		}
	})

	t.Run("HandleGetCartFail", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/carts/99", nil)
		if err != nil {
			t.Fatalf("Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			if err != nil {
				t.Logf("HandleGetCart() error, while closing request body: %v", err.Error())
			}
		}()
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("HandleGetCartFail() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("HandleDeleteCart", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/1", nil)
		if err != nil {
			t.Fatalf("HandleDeleteCart() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			if err != nil {
				t.Logf("HandleDeleteCart() error, while closing request body: %v", err.Error())
			}
		}()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("HandleDeleteCart() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("HandleDeleteCartFail", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/99", nil)
		if err != nil {
			t.Fatalf("HandleDeleteCartFail() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("HandleDeleteCartFail() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("HandleAddItem", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/carts/1/items", jsonVal)
		if err != nil {
			t.Fatalf("HandleAddItem() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("HandleAddItem() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}

		decoder := json.NewDecoder(res.Body)

		var ciRes ciResponse
		err = decoder.Decode(&ciRes)
		if err != nil {
			t.Fatalf("HandleAddItem() Failed to deserialize response error %v", err)
		}
	})

	t.Run("HandleAddItemFail", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/carts/99/items", jsonVal)
		if err != nil {
			t.Fatalf("HandleAddItemFail() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("HandleAddItemFail() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("HandleRemoveItem", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/1/items/1", jsonVal)
		if err != nil {
			t.Fatalf("HandleRemoveItem() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("HandleRemoveItem() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("HandleRemoveItemFail1", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/99/items/1", jsonVal)
		if err != nil {
			t.Fatalf("HandleRemoveItemFail1() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("HandleRemoveItemFail1() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("HandleRemoveItemFail2", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/1/items/99", jsonVal)
		if err != nil {
			t.Fatalf("HandleRemoveItemFail1() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("HandleRemoveItemFail1() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})
}
