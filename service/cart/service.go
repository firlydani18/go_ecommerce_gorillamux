package cart

import (
	"errors"
	"fmt"
	"go-ecommerce/model"
)

func getCartItemsIds(item []model.CartItem) ([]int, error) {
	productIds := make([]int, len(item))
	for i, item := range item {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than zero")
		}

		productIds[i] = item.ProductID
	}
	return productIds, nil
}

func (h *Handler) createOrder(ps []model.Product, items []model.CartItem, userId int) (int, float64, error) {
	productMap := make(map[int]model.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	// check if all products are in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, nil
	}
	// calculate total price
	totalPrice := calculateTotalPrice(items, productMap)
	// reduce quantity products in db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}
	// create order
	orderId, err := h.orderStore.CreateOrder(model.Order{
		UserID:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}

	// create order items
	for _, item := range items {
		h.orderStore.CreateOrderItem(model.OrderItem{
			OrderID:   orderId,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}
	return orderId, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []model.CartItem, products map[int]model.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product does not exist")
		}
		if product.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than zero")
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []model.CartItem, products map[int]model.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}
