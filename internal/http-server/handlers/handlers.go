package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wbTechL0/internal/storage/cache"
)

func GetOrderHandler(w http.ResponseWriter, r *http.Request, cache *cache.Cache) {
	orderIDStr := r.URL.Path[len("/order/"):]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	order, exist := cache.GetOrder(orderID)
	if !exist {
		http.Error(w, "Order with this ID does not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/index.html")
}
