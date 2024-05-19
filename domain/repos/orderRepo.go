package repos

import (
	"fmt"
	"kurs-server/domain/entities"

	"gorm.io/gorm"
)

type OrderRepo struct {
	Storage *gorm.DB
}

func (r *OrderRepo) GetPendingOrderGroup(userID uint) entities.OrderGroup {
	var orderGroup entities.OrderGroup
	if err := r.Storage.Where("user_id = ? AND status = ?", userID, "pending").First(&orderGroup).Error; err == nil {
		return orderGroup
	}

	orderGroup = entities.OrderGroup{UserID: userID, Status: "pending"}

	r.Storage.Create(&orderGroup)
	return orderGroup
}

// Creates a new order in case if it does not exist otherwise increments order's quntity by 1
func (r *OrderRepo) CreateOrder(order *entities.Order) {
	var existingOrder entities.Order
	r.Storage.Where("order_group_id = ? AND product_id = ?", order.OrderGroupID, order.ProductID).First(&existingOrder)
	if existingOrder.ID == 0 {
		r.Storage.Create(order)
		return
	}

	existingOrder.Quantity++
	r.Storage.Save(&existingOrder)

}

// Deletes an order if it's quantity is 1 otherwise decrements quantity by 1
func (r *OrderRepo) DeleteOrder(orderGroupID uint, productID uint) {
	var order entities.Order
	r.Storage.Where("order_group_id = ? AND product_id = ?", orderGroupID, productID).First(&order)

	if order.ID == 0 {
		return
	}

	if order.Quantity == 1 {
		r.Storage.Delete(&order)
		return
	}

	order.Quantity--

	r.Storage.Save(&order)
}

func (r *OrderRepo) GetOrdersByOrderGroupID(orderGroupID uint) []entities.Order {
	var orders []entities.Order
	r.Storage.Where("order_group_id = ?", orderGroupID).Find(&orders)
	return orders
}

func (r *OrderRepo) ProcessOrderGroup(orderGroupID uint, name string, surname string, city string, address string) {
	r.Storage.Model(&entities.OrderGroup{}).Where("id = ?", orderGroupID).Update("status", "processing").Update("name", name).Update("surname", surname).Update("city", city).Update("address", address)
}

func (r *OrderRepo) GetOrderGroupsByUserID(userID uint) []entities.OrderGroup {
	var orderGroups []entities.OrderGroup
	r.Storage.Where("user_id = ?", userID).Order("id desc").Find(&orderGroups)
	return orderGroups
}

func (r *OrderRepo) GetAllOrderGroups(offset int, limit int) ([]entities.OrderGroup, int64) {
	fmt.Println(limit, offset)
	var orderGroups []entities.OrderGroup
	query := r.Storage.Model(&entities.OrderGroup{}).Where("status IN ('processing', 'ready', 'completed')")
	var count int64
	query.Count(&count)
	query.Limit(limit).Offset(offset).Order("id desc").Find(&orderGroups)
	return orderGroups, count
}

func (r *OrderRepo) MarkOrderGroup(mark string, orderGroupID uint) {
	r.Storage.Model(&entities.OrderGroup{}).Where("id = ?", orderGroupID).Update("status", mark)
}

func (r *OrderRepo) DeleteOrderGroup(orderGroupID uint) {
	r.Storage.Delete(&entities.OrderGroup{}, orderGroupID)
}
