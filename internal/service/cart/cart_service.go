package cart

import (
	"context"
	"fmt"

	"github.com/crabrus/the-workshop/internal/domain/entity"
	"github.com/crabrus/the-workshop/internal/domain/repository"
	"github.com/crabrus/the-workshop/internal/service/product"
	"github.com/google/uuid"
)

type CartService interface {
	AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*CartItemDTO, error)
	UpdateItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID, quantity int) (*CartItemDTO, error)
	RemoveItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID) error
	GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	cartRepo       repository.CartItemRepository
	productService product.ProductService
}

func NewService(cartRepo repository.CartItemRepository, productService product.ProductService) CartService {
	return &service{
		cartRepo:       cartRepo,
		productService: productService,
	}
}

// AddItem добавляет товар в корзину или увеличивает количество если товар уже есть
func (s *service) AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*CartItemDTO, error) {
	// Валидация
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	// Проверить что товар существует и есть на складе
	prod, err := s.productService.GetByID(ctx, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if prod.Stock < quantity {
		return nil, ErrInsufficientStock
	}

	// Проверить есть ли уже этот товар в корзине пользователя
	existingItems, _, err := s.cartRepo.List(ctx, repository.CartItemFilter{
		UserID: userIDToString(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing items: %w", err)
	}

	// Ищем товар в корзине
	var existingItem *entity.CartItem
	for _, item := range existingItems {
		if item.ProductID == productID {
			existingItem = item
			break
		}
	}

	// Если товар уже есть - обновляем количество
	if existingItem != nil {
		newQuantity := existingItem.Quantity + quantity

		// Проверяем что всё поместится на складе
		if newQuantity > prod.Stock {
			return nil, ErrInsufficientStock
		}

		existingItem.Quantity = newQuantity
		if err := s.cartRepo.Update(ctx, existingItem); err != nil {
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}

		return s.buildCartItemDTO(existingItem, prod), nil
	}

	// Создаём новый CartItem
	cartItem := &entity.CartItem{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	if err := s.cartRepo.Create(ctx, cartItem); err != nil {
		return nil, fmt.Errorf("failed to create cart item: %w", err)
	}

	return s.buildCartItemDTO(cartItem, prod), nil
}

// UpdateItem обновляет количество товара в корзине
func (s *service) UpdateItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID, quantity int) (*CartItemDTO, error) {
	// Получить CartItem
	cartItem, err := s.cartRepo.GetByID(ctx, cartItemID)
	if err != nil {
		return nil, ErrCartItemNotFound
	}

	// Проверить что он принадлежит пользователю
	if cartItem.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Если quantity = 0, удалить товар
	if quantity == 0 {
		if err := s.cartRepo.Delete(ctx, cartItemID); err != nil {
			return nil, fmt.Errorf("failed to delete cart item: %w", err)
		}
		return nil, nil
	}

	// Валидация
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	// Проверить остаток на складе
	prod, err := s.productService.GetByID(ctx, cartItem.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if quantity > prod.Stock {
		return nil, ErrInsufficientStock
	}

	// Обновить количество
	cartItem.Quantity = quantity
	if err := s.cartRepo.Update(ctx, cartItem); err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	return s.buildCartItemDTO(cartItem, prod), nil
}

// RemoveItem удаляет товар из корзины
func (s *service) RemoveItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID) error {
	// Получить CartItem
	cartItem, err := s.cartRepo.GetByID(ctx, cartItemID)
	if err != nil {
		return ErrCartItemNotFound
	}

	// Проверить что он принадлежит пользователю
	if cartItem.UserID != userID {
		return ErrUnauthorized
	}

	// Удалить
	if err := s.cartRepo.Delete(ctx, cartItemID); err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	return nil
}

// GetCart получает полную информацию о корзине пользователя
func (s *service) GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	// Получить все CartItems пользователя
	cartItems, _, err := s.cartRepo.List(ctx, repository.CartItemFilter{
		UserID: userIDToString(userID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list cart items: %w", err)
	}

	response := &CartResponse{
		Items:       make([]*CartItemDTO, 0),
		TotalAmount: 0,
		TotalItems:  len(cartItems),
	}

	// Для каждого CartItem достать информацию о продукте
	for _, item := range cartItems {
		prod, err := s.productService.GetByID(ctx, item.ProductID)
		if err != nil {
			// Пропускаем товары которые удалены (или обработать ошибку иначе)
			continue
		}

		dto := s.buildCartItemDTO(item, prod)
		response.Items = append(response.Items, dto)
		response.TotalAmount += dto.Sum
	}

	return response, nil
}

// Clear очищает корзину пользователя
func (s *service) Clear(ctx context.Context, userID uuid.UUID) error {
	// Получить все CartItems пользователя
	cartItems, _, err := s.cartRepo.List(ctx, repository.CartItemFilter{
		UserID: userIDToString(userID),
	})
	if err != nil {
		return fmt.Errorf("failed to list cart items: %w", err)
	}

	// Удалить каждый CartItem
	for _, item := range cartItems {
		if err := s.cartRepo.Delete(ctx, item.ID); err != nil {
			return fmt.Errorf("failed to delete cart item: %w", err)
		}
	}

	return nil
}

// buildCartItemDTO вспомогательная функция для создания CartItemDTO
func (s *service) buildCartItemDTO(item *entity.CartItem, prod *entity.Product) *CartItemDTO {
	sum := float64(item.Quantity) * prod.Price
	return &CartItemDTO{
		ID:           item.ID,
		ProductID:    item.ProductID,
		ProductName:  prod.Name,
		ProductPrice: prod.Price,
		Quantity:     item.Quantity,
		Sum:          sum,
	}
}

// userIDToString вспомогательная функция для конвертации UUID в string для фильтра
func userIDToString(id uuid.UUID) *string {
	str := id.String()
	return &str
}
