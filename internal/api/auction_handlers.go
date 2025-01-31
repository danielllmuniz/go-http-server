package api

import (
	"errors"
	"net/http"

	"github.com/danielllmuniz/go-http-server/internal/jsonutils"
	"github.com/danielllmuniz/go-http-server/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	// Get the product id from the URL
	rawProductId := chi.URLParam(r, "product_id")

	// Validate the product id
	productId, err := uuid.Parse(rawProductId)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"message": "Invalid product id - must be a valid UUID",
		})
		return
	}

	// Check if the product exists
	_, err = api.ProductService.GetProductById(r.Context(), productId)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"message": "Product not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unexpected error, please try again later",
		})
		return
	}

	// Get the user id from the session
	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unexpected error, please try again later",
		})
		return
	}

	// Get the auction room
	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productId]
	api.AuctionLobby.Unlock()

	// Check if the auction has ended
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"message": "the auction has ended",
		})
	}

	// Upgrade the connection to a websocket protocol
	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "could not upgrade connection to a websocket protocol",
		})
		return
	}

	// Create a new client
	client := services.NewClient(room, conn, userId)

	// Register the client
	room.Register <- client
	go client.ReadEventLoop()
	go client.WriteEventLoop()
}
