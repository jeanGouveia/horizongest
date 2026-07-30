package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var _ ports.NotificationsRepository = (*GormNotificationsRepository)(nil)

type GormNotificationsRepository struct{ db *gorm.DB }

func NewGormNotificationsRepository(db *gorm.DB) *GormNotificationsRepository {
	return &GormNotificationsRepository{db: db}
}

func (r *GormNotificationsRepository) GetNotifications(ctx context.Context) (*domain.Notifications, error) {
	notifications := &domain.Notifications{}

	// Usar ApplyTenantFilter para isolamento de tenant
	query := ApplyTenantFilter(ctx, r.db)

	// Pedidos pendentes
	var pendingOrders int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("status = ? AND deleted_at IS NULL", "pending").
		Count(&pendingOrders)
	notifications.PendingOrders = int(pendingOrders)

	// Estoque baixo
	var lowStockCount int64
	query.WithContext(ctx).Model(&GormIngredient{}).
		Where("stock_quantity < min_stock AND deleted_at IS NULL").
		Count(&lowStockCount)
	notifications.LowStockCount = int(lowStockCount)

	// Produtos sem foto
	var productsWithoutPhoto int64
	query.WithContext(ctx).Model(&GormProduct{}).
		Where("(photo_url = '' OR photo_url IS NULL) AND deleted_at IS NULL").
		Count(&productsWithoutPhoto)
	notifications.ProductsWithoutPhoto = int(productsWithoutPhoto)

	// Promoções vencidas
	now := time.Now()
	var expiredPromotions int64
	query.WithContext(ctx).Model(&GormProduct{}).
		Where("promotion_end IS NOT NULL AND FROM_UNIXTIME(promotion_end) < ? AND deleted_at IS NULL", now).
		Count(&expiredPromotions)
	notifications.ExpiredPromotions = int(expiredPromotions)

	return notifications, nil
}
