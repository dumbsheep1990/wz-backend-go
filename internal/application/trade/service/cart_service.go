package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yourusername/wz-backend-go/internal/application/trade/dto"
	"github.com/yourusername/wz-backend-go/internal/domain/shared/event"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/entity"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/repository"
	"github.com/yourusername/wz-backend-go/internal/domain/trade/valueobject"
	"github.com/yourusername/wz-backend-go/internal/infrastructure/database"
)

// CartService 购物车应用服务
type CartService struct {
	cartRepository repository.CartRepository
	eventBus       event.EventBus
	validator      *validator.Validate
	unitOfWork     database.UnitOfWork
}

// NewCartService 创建购物车应用服务
func NewCartService(
	cartRepository repository.CartRepository,
	eventBus event.EventBus,
	unitOfWork database.UnitOfWork,
) *CartService {
	return &CartService{
		cartRepository: cartRepository,
		eventBus:       eventBus,
		validator:      validator.New(),
		unitOfWork:     unitOfWork,
	}
}

// GetCart 获取用户的购物车
func (s *CartService) GetCart(ctx context.Context, userIDStr string) (*dto.CartDTO, error) {
	// 验证用户ID
	userID, err := valueobject.NewUserID(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 查找用户的购物车
	cart, err := s.cartRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询购物车失败: %w", err)
	}

	// 如果用户没有购物车，创建一个新的
	if cart == nil {
		newCart, err := entity.NewCart(userID)
		if err != nil {
			return nil, fmt.Errorf("创建购物车失败: %w", err)
		}

		// 保存新购物车
		err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
			if err := s.cartRepository.Save(ctx, newCart); err != nil {
				return fmt.Errorf("保存购物车失败: %w", err)
			}

			// 发布领域事件
			for _, event := range newCart.GetDomainEvents() {
				if err := s.eventBus.Publish(ctx, event); err != nil {
					return fmt.Errorf("发布购物车事件失败: %w", err)
				}
			}

			// 清除已处理的事件
			newCart.ClearDomainEvents()

			return nil
		})

		if err != nil {
			return nil, err
		}

		cart = newCart
	}

	// 转换为DTO
	return s.toCartDTO(cart)
}

// AddItem 添加商品到购物车
func (s *CartService) AddItem(ctx context.Context, req dto.AddCartItemRequest) (*dto.CartDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证添加购物车项请求失败: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 转换商品ID
	productID, err := valueobject.NewProductID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("无效的商品ID: %w", err)
	}

	// 转换价格
	price, err := valueobject.NewMoney(req.Price, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("无效的商品价格: %w", err)
	}

	// 转换数量
	quantity, err := valueobject.NewQuantity(req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("无效的商品数量: %w", err)
	}

	// 查找用户的购物车
	cart, err := s.cartRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询购物车失败: %w", err)
	}

	// 如果用户没有购物车，创建一个新的
	if cart == nil {
		cart, err = entity.NewCart(userID)
		if err != nil {
			return nil, fmt.Errorf("创建购物车失败: %w", err)
		}
	}

	// 创建购物车项
	item, err := entity.NewCartItem(productID, req.Name, price, quantity)
	if err != nil {
		return nil, fmt.Errorf("创建购物车项失败: %w", err)
	}

	// 添加到购物车
	if err := cart.AddItem(item); err != nil {
		return nil, fmt.Errorf("添加商品到购物车失败: %w", err)
	}

	// 使用工作单元保存购物车并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.cartRepository.Save(ctx, cart); err != nil {
			return fmt.Errorf("保存购物车失败: %w", err)
		}

		// 发布领域事件
		for _, event := range cart.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布购物车事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		cart.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toCartDTO(cart)
}

// UpdateItemQuantity 更新购物车中商品的数量
func (s *CartService) UpdateItemQuantity(ctx context.Context, req dto.UpdateCartItemQuantityRequest) (*dto.CartDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证更新购物车项数量请求失败: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 转换商品ID
	productID, err := valueobject.NewProductID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("无效的商品ID: %w", err)
	}

	// 转换数量
	quantity, err := valueobject.NewQuantity(req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("无效的商品数量: %w", err)
	}

	// 查找用户的购物车
	cart, err := s.cartRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询购物车失败: %w", err)
	}

	if cart == nil {
		return nil, errors.New("未找到购物车")
	}

	// 更新商品数量
	if err := cart.UpdateItemQuantity(productID, quantity); err != nil {
		return nil, fmt.Errorf("更新购物车商品数量失败: %w", err)
	}

	// 使用工作单元保存购物车并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.cartRepository.Save(ctx, cart); err != nil {
			return fmt.Errorf("保存购物车失败: %w", err)
		}

		// 发布领域事件
		for _, event := range cart.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布购物车事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		cart.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toCartDTO(cart)
}

// RemoveItem 从购物车中移除商品
func (s *CartService) RemoveItem(ctx context.Context, req dto.RemoveCartItemRequest) (*dto.CartDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证移除购物车项请求失败: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 转换商品ID
	productID, err := valueobject.NewProductID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("无效的商品ID: %w", err)
	}

	// 查找用户的购物车
	cart, err := s.cartRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询购物车失败: %w", err)
	}

	if cart == nil {
		return nil, errors.New("未找到购物车")
	}

	// 从购物车中移除商品
	if err := cart.RemoveItem(productID); err != nil {
		return nil, fmt.Errorf("移除购物车商品失败: %w", err)
	}

	// 使用工作单元保存购物车并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.cartRepository.Save(ctx, cart); err != nil {
			return fmt.Errorf("保存购物车失败: %w", err)
		}

		// 发布领域事件
		for _, event := range cart.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布购物车事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		cart.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toCartDTO(cart)
}

// ClearCart 清空购物车
func (s *CartService) ClearCart(ctx context.Context, req dto.ClearCartRequest) (*dto.CartDTO, error) {
	// 验证请求
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("验证清空购物车请求失败: %w", err)
	}

	// 转换用户ID
	userID, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID: %w", err)
	}

	// 查找用户的购物车
	cart, err := s.cartRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询购物车失败: %w", err)
	}

	if cart == nil {
		return nil, errors.New("未找到购物车")
	}

	// 清空购物车
	cart.Clear()

	// 使用工作单元保存购物车并发布事件
	err = s.unitOfWork.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.cartRepository.Save(ctx, cart); err != nil {
			return fmt.Errorf("保存购物车失败: %w", err)
		}

		// 发布领域事件
		for _, event := range cart.GetDomainEvents() {
			if err := s.eventBus.Publish(ctx, event); err != nil {
				return fmt.Errorf("发布购物车事件失败: %w", err)
			}
		}

		// 清除已处理的事件
		cart.ClearDomainEvents()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return s.toCartDTO(cart)
}

// 将购物车实体转换为DTO
func (s *CartService) toCartDTO(cart *entity.Cart) (*dto.CartDTO, error) {
	// 转换购物车项
	cartItems := cart.Items()
	items := make([]dto.CartItemDTO, 0, len(cartItems))
	
	for _, item := range cartItems {
		subtotal, err := item.Subtotal()
		if err != nil {
			return nil, fmt.Errorf("计算购物车项小计金额失败: %w", err)
		}
		
		items = append(items, dto.CartItemDTO{
			ProductID: item.ProductID().String(),
			Name:      item.Name(),
			Price:     item.Price().Amount(),
			Currency:  item.Price().Currency(),
			Quantity:  item.Quantity().Value(),
			AddedAt:   item.AddedAt(),
			Subtotal:  subtotal.Amount(),
		})
	}

	// 计算总金额
	var totalAmount float64
	var currency string
	
	if !cart.IsEmpty() {
		total, err := cart.TotalAmount()
		if err != nil {
			return nil, fmt.Errorf("计算购物车总金额失败: %w", err)
		}
		totalAmount = total.Amount()
		currency = total.Currency()
	}

	// 构造DTO
	return &dto.CartDTO{
		ID:            cart.ID().String(),
		UserID:        cart.UserID().String(),
		Items:         items,
		TotalItems:    cart.ItemCount(),
		TotalQuantity: cart.TotalQuantity(),
		TotalAmount:   totalAmount,
		Currency:      currency,
		CreatedAt:     cart.CreatedAt(),
		UpdatedAt:     cart.UpdatedAt(),
	}, nil
}
