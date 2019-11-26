package routing

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"rockspoon-cart/internal/domain"
	"rockspoon-cart/internal/repository"
)

type EndpointFactory interface {
	HandleHeartbeat() HttpEndpoint
	HandleNewCart() HttpEndpoint
	HandleAddItem(cartIDParam string) HttpEndpoint
	HandleRemoveItem(cartIDParam, itemIDParam string) HttpEndpoint
	HandleDeleteCart(cartIDParam string) HttpEndpoint
	HandleGetCart(cartIDParam string) HttpEndpoint
}

func RouterInit(r *mux.Router, write repository.CartWriteRepository, read repository.CartReadRepository) *mux.Router {
	fac := NewEndpointFactory(write, read)

	// Heartbeat
	r.HandleFunc(
		"/heartbeat", Json(fac.HandleHeartbeat()),
	).Methods(http.MethodGet)

	//Cart commands
	r.HandleFunc(
		"/carts", Json(fac.HandleNewCart()),
	).Methods(http.MethodPost)

	r.HandleFunc(
		"/carts/{id}", Json(fac.HandleGetCart("id")),
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/carts/{id}", Json(fac.HandleDeleteCart("id")),
	).Methods(http.MethodDelete)

	// Items commands
	r.HandleFunc(
		"/carts/{id}/items", Json(fac.HandleAddItem("id")),
	).Methods(http.MethodPost)

	r.HandleFunc(
		"/carts/{cart}/items/{item}", Json(fac.HandleRemoveItem("cart", "item")),
	).Methods(http.MethodDelete)

	return r
}

type endpointFactory struct {
	write repository.CartWriteRepository
	read  repository.CartReadRepository
}

func NewEndpointFactory(write repository.CartWriteRepository, read repository.CartReadRepository) EndpointFactory {
	return &endpointFactory{
		write: write,
		read:  read,
	}
}

func (f *endpointFactory) HandleHeartbeat() HttpEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HttpResponse {
		return OK(struct {
			Result string `json:"result"`
		}{
			Result: "OK",
		})
	}
}

func (f *endpointFactory) HandleNewCart() HttpEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HttpResponse {
		res, err := f.write.CreateCart(r.Context())
		if err != nil {
			return BadRequestErrResp(err)
		}
		return OK(res)
	}
}

func (f *endpointFactory) HandleAddItem(cartIDParam string) HttpEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HttpResponse {
		decoder := json.NewDecoder(r.Body)

		var cartItem domain.CartItem
		err := decoder.Decode(&cartItem)

		vars := mux.Vars(r)
		cartID := vars[cartIDParam]

		if err != nil {
			return BadRequestErrResp(err)
		}

		res, err := f.write.AddItem(r.Context(), cartID, cartItem.Product, cartItem.Quantity)

		if err != nil {
			return BadRequestErrResp(err)
		}

		return OK(res)
	}
}

func (f *endpointFactory) HandleRemoveItem(cartIDParam, itemIDParam string) HttpEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HttpResponse {

		vars := mux.Vars(r)
		cartID := vars[cartIDParam]
		itemID := vars[itemIDParam]

		err := f.write.RemoveItem(r.Context(), cartID, itemID)

		if err != nil {
			return BadRequestErrResp(err)
		}

		return OK(struct{}{})
	}
}

func (f *endpointFactory) HandleDeleteCart(cartIDParam string) HttpEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HttpResponse {

		vars := mux.Vars(r)
		cartID := vars[cartIDParam]

		err := f.write.DeleteCart(r.Context(), cartID)

		if err != nil {
			return BadRequestErrResp(err)
		}

		return OK(struct{}{})
	}
}

func (f *endpointFactory) HandleGetCart(cartIDParam string) HttpEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HttpResponse {

		vars := mux.Vars(r)
		id := vars[cartIDParam]

		res, err := f.read.Cart(r.Context(), id)

		if err != nil {
			return BadRequestErrResp(err)
		}

		return OK(res)
	}
}
