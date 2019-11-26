package routing

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initRouting() *mux.Router {
	repo := InitMockRepo()
	router := mux.NewRouter()
	return RouterInit(router, repo, repo)
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

	t.Run("Test_HandleNewCart", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/carts", nil)
		if err != nil {
			t.Errorf("Test_HandleNewCart() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Test_HandleNewCart() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}

		decoder := json.NewDecoder(res.Body)

		var cartRes cartResponse
		err = decoder.Decode(&cartRes)
		if err != nil {
			t.Errorf("Test_HandleNewCart() Failed to deserialize response error %v", err)
		}
	})

	t.Run("Test_HandleGetCart", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/carts/1", nil)
		if err != nil {
			t.Errorf("Test_HandleGetCart() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Test_HandleGetCart() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}

		decoder := json.NewDecoder(res.Body)

		var response cartResponse
		err = decoder.Decode(&response)
		if err != nil {
			t.Errorf("Test_HandleGetCartFail() Failed to deserialize response error %v", err)
		}
	})

	t.Run("Test_HandleGetCartFail", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/carts/99", nil)
		if err != nil {
			t.Errorf("Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Test_HandleGetCartFail() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("Test_HandleDeleteCart", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/1", nil)
		if err != nil {
			t.Errorf("Test_HandleDeleteCart() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Test_HandleDeleteCart() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("Test_HandleDeleteCartFail", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/99", nil)
		if err != nil {
			t.Errorf("Test_HandleDeleteCartFail() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Test_HandleDeleteCartFail() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("Test_HandleAddItem", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/carts/1/items", jsonVal)
		if err != nil {
			t.Errorf("Test_HandleAddItem() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("Test_HandleAddItem() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Test_HandleAddItem() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}

		decoder := json.NewDecoder(res.Body)

		var ciRes ciResponse
		err = decoder.Decode(&ciRes)
		if err != nil {
			t.Errorf("Test_HandleAddItem() Failed to deserialize response error %v", err)
		}
	})

	t.Run("Test_HandleAddItemFail", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/carts/99/items", jsonVal)
		if err != nil {
			t.Errorf("Test_HandleAddItemFail() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("Test_HandleAddItemFail() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Test_HandleAddItemFail() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("Test_HandleRemoveItem", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/1/items/1", jsonVal)
		if err != nil {
			t.Errorf("Test_HandleRemoveItem() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("Test_HandleRemoveItem() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Test_HandleRemoveItem() code=%v, expected %v", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("Test_HandleRemoveItemFail1", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/99/items/1", jsonVal)
		if err != nil {
			t.Errorf("Test_HandleRemoveItemFail1() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("Test_HandleRemoveItemFail1() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Test_HandleRemoveItemFail1() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("Test_HandleRemoveItemFail2", func(t *testing.T) {
		router := initRouting()
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodDelete, "/carts/1/items/99", jsonVal)
		if err != nil {
			t.Errorf("Test_HandleRemoveItemFail1() Request initialization error %v", err)
		}
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf("Test_HandleRemoveItemFail1() Request initialization error %v", err)
		}

		router.ServeHTTP(w, request)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Test_HandleRemoveItemFail1() code=%v, expected %v", res.StatusCode, http.StatusBadRequest)
		}
	})
}
