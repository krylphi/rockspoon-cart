package routing

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Krylphi/rockspoon-cart/internal/domain"
	"github.com/Krylphi/rockspoon-cart/internal/repository"

	"github.com/gorilla/mux"
)

// EndpointFactory is meant to produce handlers for endpoints.
type EndpointFactory interface {
	HandleHeartbeat() HTTPEndpoint
	HandleNewCart() HTTPEndpoint
	HandleAddItem(cartIDParam string) HTTPEndpoint
	HandleRemoveItem(cartIDParam, itemIDParam string) HTTPEndpoint
	HandleDeleteCart(cartIDParam string) HTTPEndpoint
	HandleGetCart(cartIDParam string) HTTPEndpoint
}

// RouterInit initialize routing and assigning endpoints with handlers.
func RouterInit(write repository.CartWriteRepository, read repository.CartReadRepository) http.Handler {
	fac := NewEndpointFactory(write, read)
	r := mux.NewRouter()

	// Heartbeat
	r.HandleFunc(
		"/heartbeat", JSON(fac.HandleHeartbeat()),
	).Methods(http.MethodGet)

	//Cart commands
	r.HandleFunc(
		"/carts", JSON(fac.HandleNewCart()),
	).Methods(http.MethodPost)

	r.HandleFunc(
		"/carts/{id}", JSON(fac.HandleGetCart("id")),
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/carts/{id}", JSON(fac.HandleDeleteCart("id")),
	).Methods(http.MethodDelete)

	// Items commands
	r.HandleFunc(
		"/carts/{id}/items", JSON(fac.HandleAddItem("id")),
	).Methods(http.MethodPost)

	r.HandleFunc(
		"/carts/{cart}/items/{item}", JSON(fac.HandleRemoveItem("cart", "item")),
	).Methods(http.MethodDelete)

	return r
}

type endpointFactory struct {
	write repository.CartWriteRepository
	read  repository.CartReadRepository
}

// NewEndpointFactory produces instance of EndpointFactory implementation.
func NewEndpointFactory(write repository.CartWriteRepository, read repository.CartReadRepository) EndpointFactory {
	return &endpointFactory{
		write: write,
		read:  read,
	}
}

func (f *endpointFactory) HandleHeartbeat() HTTPEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HTTPResponse {
		return OK(struct {
			Result string `json:"result"`
		}{
			Result: "OK",
		})
	}
}

func (f *endpointFactory) HandleNewCart() HTTPEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HTTPResponse {
		res, err := f.write.CreateCart(r.Context())
		if err != nil {
			return BadRequestErrResp(err)
		}
		return OK(res)
	}
}

func (f *endpointFactory) HandleAddItem(cartIDParam string) HTTPEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HTTPResponse {
		decoder := json.NewDecoder(r.Body)
		defer func() {
			err := r.Body.Close()
			if err != nil {
				log.Printf("HandleAddItem() error, while closing request body: %v", err.Error())
			}
		}()

		var cartItem domain.CartItem
		err := decoder.Decode(&cartItem)
		if err != nil {
			return BadRequestErrResp(err)
		}

		vars := mux.Vars(r)
		cartID := vars[cartIDParam]

		res, err := f.write.AddItem(r.Context(), cartID, cartItem.Product, cartItem.Quantity)
		if err != nil {
			return BadRequestErrResp(err)
		}

		return OK(res)
	}
}

func (f *endpointFactory) HandleRemoveItem(cartIDParam, itemIDParam string) HTTPEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HTTPResponse {
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

func (f *endpointFactory) HandleDeleteCart(cartIDParam string) HTTPEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HTTPResponse {
		vars := mux.Vars(r)
		cartID := vars[cartIDParam]

		err := f.write.DeleteCart(r.Context(), cartID)
		if err != nil {
			return BadRequestErrResp(err)
		}

		return OK(struct{}{})
	}
}

func (f *endpointFactory) HandleGetCart(cartIDParam string) HTTPEndpoint {
	return func(w http.ResponseWriter, r *http.Request) HTTPResponse {
		vars := mux.Vars(r)
		id := vars[cartIDParam]

		res, err := f.read.Cart(r.Context(), id)
		if err != nil {
			return BadRequestErrResp(err)
		}

		return OK(res)
	}
}
