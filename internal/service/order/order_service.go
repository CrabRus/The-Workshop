package order

import (
	"context"
	"fmt"
	"time"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/crabrus/the-workshop/internal/service/cart"
	"github.com/crabrus/the-workshop/internal/service/product"
	"github.com/google/uuid"
)

type OrderService interface {
	// Public methods
	CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*OrderDTO, error)
	GetOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*OrderDTO, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID, filter repository.OrderFilter) (*OrderListResponse, error)
	CancelOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) error

	// Admin methods
	GetAllOrders(ctx context.Context, filter repository.OrderFilter) (*OrderListResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error
}

type service struct {
	orderRepo      repository.OrderRepository
	orderItemRepo  repository.OrderItemRepository
	cartService    cart.CartService
	productService product.ProductService
	productRepo    repository.ProductRepository
}

func NewService(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	cartService cart.CartService,
	productService product.ProductService,
	productRepo repository.ProductRepository,
) OrderService {
	return &service{
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
		cartService:    cartService,
		productService: productService,
		productRepo:    productRepo,
	}
}

// ---------- PUBLIC ----------

// CreateOrder создаёт новый заказ из корзины пользователя
func (s *service) CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*OrderDTO, error) {
	// Валидация данных
	if err := ValidateCreateOrder(req); err != nil {
		return nil, err
	}

	// Получаем корзину пользователя
	cartResp, err := s.cartService.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Проверяем что корзина не пустая
	if len(cartResp.Items) == 0 {
		return nil, ErrEmptyCart
	}

	// Создаём заказ
	shippingAddr := entity.ShippingAddress{
		FullName:    *req.ShippingAddress.FullName,
		PhoneNumber: *req.ShippingAddress.PhoneNumber,
		Email:       *req.ShippingAddress.Email,
		Country:     *req.ShippingAddress.Country,
		City:        *req.ShippingAddress.City,
		PostalCode:  *req.ShippingAddress.PostalCode,
		AddressLine: *req.ShippingAddress.AddressLine,
	}

	order := &entity.Order{
		ID:              uuid.New(),
		UserID:          userID,
		Status:          entity.OrderStatusPending,
		TotalAmount:     cartResp.TotalAmount,
		ShippingAddress: shippingAddr,
		PaymentMethod:   *req.PaymentMethod,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Сохраняем заказ в БД
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Добавляем товары из корзины в заказ и уменьшаем остатки
	for _, cartItem := range cartResp.Items {
		// Создаём OrderItem
		orderItem := &entity.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.ProductPrice,
			CreatedAt: time.Now(),
		}

		if err := s.orderItemRepo.Create(ctx, orderItem); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		// Уменьшаем остатки товара
		if err := s.productRepo.DecreaseStock(ctx, cartItem.ProductID, cartItem.Quantity); err != nil {
			return nil, fmt.Errorf("failed to decrease stock: %w", err)
		}
	}

	// Очищаем корзину после успешного создания заказа
	if err := s.cartService.Clear(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Получаем полную информацию о заказе
	orderDTO, err := s.buildOrderDTO(ctx, order)
	if err != nil {
		return nil, err
	}

	return orderDTO, nil
}

// GetOrder получает заказ по ID (для пользователя - только свои заказы)
func (s *service) GetOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*OrderDTO, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// Проверяем что заказ принадлежит пользователю
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}

	orderDTO, err := s.buildOrderDTO(ctx, order)
	if err != nil {
		return nil, err
	}

	return orderDTO, nil
}

// GetUserOrders получает все заказы пользователя
func (s *service) GetUserOrders(ctx context.Context, userID uuid.UUID, filter repository.OrderFilter) (*OrderListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	orders, total, err := s.orderRepo.GetByUserID(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	orderDTOs := make([]*OrderDTO, len(orders))
	for i, order := range orders {
		items, err := s.orderItemRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}

		itemDTOs := s.buildOrderItemDTOs(items)
		orderDTOs[i] = &OrderDTO{
			ID:              order.ID,
			UserID:          order.UserID,
			Status:          order.Status,
			TotalAmount:     order.TotalAmount,
			ShippingAddress: order.ShippingAddress,
			PaymentMethod:   order.PaymentMethod,
			Items:           itemDTOs,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
		}
	}

	return &OrderListResponse{
		Orders: orderDTOs,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}

// CancelOrder отменяет заказ и возвращает товары на склад
func (s *service) CancelOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	// Проверяем что заказ принадлежит пользователю
	if order.UserID != userID {
		return ErrUnauthorized
	}

	// Проверяем что заказ можно отменить
	if order.Status == entity.OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}

	if order.Status == entity.OrderStatusDelivered || order.Status == entity.OrderStatusShipped {
		return fmt.Errorf("cannot cancel delivered or shipped order")
	}

	// Получаем все товары в заказе
	items, err := s.orderItemRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}

	// Возвращаем товары на склад
	for _, item := range items {
		// Увеличиваем остатки товара в БД
		// В реальном приложении нужна транзакция
		prod, err := s.productService.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}

		newStock := prod.Stock + item.Quantity
		prod.Stock = newStock

		if err := s.productRepo.Update(ctx, prod); err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}
	}

	// Меняем статус заказа на отменённый
	if err := s.orderRepo.UpdateStatus(ctx, orderID, entity.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// ---------- ADMIN ----------

// GetAllOrders получает все заказы в системе
func (s *service) GetAllOrders(ctx context.Context, filter repository.OrderFilter) (*OrderListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	orders, total, err := s.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	orderDTOs := make([]*OrderDTO, len(orders))
	for i, order := range orders {
		items, err := s.orderItemRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}

		itemDTOs := s.buildOrderItemDTOs(items)
		orderDTOs[i] = &OrderDTO{
			ID:              order.ID,
			UserID:          order.UserID,
			Status:          order.Status,
			TotalAmount:     order.TotalAmount,
			ShippingAddress: order.ShippingAddress,
			PaymentMethod:   order.PaymentMethod,
			Items:           itemDTOs,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
		}
	}

	return &OrderListResponse{
		Orders: orderDTOs,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}

// UpdateOrderStatus изменяет статус заказа (только для админа)
func (s *service) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	// Валидируем статус
	validStatuses := map[string]bool{
		entity.OrderStatusPending:   true,
		entity.OrderStatusConfirmed: true,
		entity.OrderStatusShipped:   true,
		entity.OrderStatusDelivered: true,
		entity.OrderStatusCancelled: true,
	}

	if !validStatuses[status] {
		return ErrInvalidStatus
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

// ---------- HELPERS ----------

func (s *service) buildOrderDTO(ctx context.Context, order *entity.Order) (*OrderDTO, error) {
	items, err := s.orderItemRepo.GetByOrderID(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	itemDTOs := s.buildOrderItemDTOs(items)

	return &OrderDTO{
		ID:              order.ID,
		UserID:          order.UserID,
		Status:          order.Status,
		TotalAmount:     order.TotalAmount,
		ShippingAddress: order.ShippingAddress,
		PaymentMethod:   order.PaymentMethod,
		Items:           itemDTOs,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}, nil
}

func (s *service) buildOrderItemDTOs(items []*entity.OrderItem) []*OrderItemDTO {
	itemDTOs := make([]*OrderItemDTO, len(items))
	for i, item := range items {
		itemDTOs[i] = &OrderItemDTO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		}
	}
	return itemDTOs
}
