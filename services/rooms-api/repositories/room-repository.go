package repositories

import (
	"context"
	"rooms-api/domain"
	"rooms-api/utils"

	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	// Verificar si ya existe una habitación con el mismo número
	var existingRoom domain.Room
	result := r.db.WithContext(ctx).Where("number = ?", room.Number).First(&existingRoom)

	if result.Error == nil {
		return utils.ErrRoomAlreadyExists
	}

	if result.Error != gorm.ErrRecordNotFound {
		return utils.ErrDatabaseError
	}

	// Crear la habitación
	if err := r.db.WithContext(ctx).Create(room).Error; err != nil {
		return utils.ErrDatabaseError
	}

	return nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id uint) (*domain.Room, error) {
	var room domain.Room
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&room)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, utils.ErrRoomNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	return &room, nil
}

func (r *RoomRepository) GetByNumber(ctx context.Context, number string) (*domain.Room, error) {
	var room domain.Room
	result := r.db.WithContext(ctx).Where("number = ?", number).First(&room)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, utils.ErrRoomNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	return &room, nil
}

func (r *RoomRepository) GetAll(ctx context.Context, filter domain.RoomFilter, page, limit int) ([]domain.Room, int64, error) {
	var rooms []domain.Room
	var total int64

	// Construir query base
	query := r.db.WithContext(ctx).Model(&domain.Room{})

	// Aplicar filtros
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Floor != nil {
		query = query.Where("floor = ?", *filter.Floor)
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.HasWifi != nil {
		query = query.Where("has_wifi = ?", *filter.HasWifi)
	}
	if filter.HasAC != nil {
		query = query.Where("has_ac = ?", *filter.HasAC)
	}
	if filter.HasTV != nil {
		query = query.Where("has_tv = ?", *filter.HasTV)
	}
	if filter.HasMinibar != nil {
		query = query.Where("has_minibar = ?", *filter.HasMinibar)
	}

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	// Aplicar paginación y ordenamiento
	offset := (page - 1) * limit
	if err := query.Order("number ASC").Offset(offset).Limit(limit).Find(&rooms).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	return rooms, total, nil
}

func (r *RoomRepository) Update(ctx context.Context, id uint, updateData domain.UpdateRoomRequest) (*domain.Room, error) {
	// Verificar que la habitación existe
	var room domain.Room
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&room)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, utils.ErrRoomNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	// Si se está actualizando el número, verificar que no exista otra habitación con ese número
	if updateData.Number != nil && *updateData.Number != "" {
		var existingRoom domain.Room
		result := r.db.WithContext(ctx).Where("number = ? AND id != ?", *updateData.Number, id).First(&existingRoom)

		if result.Error == nil {
			return nil, utils.ErrRoomAlreadyExists
		}

		if result.Error != gorm.ErrRecordNotFound {
			return nil, utils.ErrDatabaseError
		}
	}

	// Preparar los datos a actualizar
	updates := make(map[string]interface{})

	if updateData.Number != nil {
		updates["number"] = *updateData.Number
	}
	if updateData.Type != nil {
		updates["type"] = *updateData.Type
	}
	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}
	if updateData.Price != nil {
		updates["price"] = *updateData.Price
	}
	if updateData.Description != nil {
		updates["description"] = *updateData.Description
	}
	if updateData.Capacity != nil {
		updates["capacity"] = *updateData.Capacity
	}
	if updateData.Floor != nil {
		updates["floor"] = *updateData.Floor
	}
	if updateData.HasWifi != nil {
		updates["has_wifi"] = *updateData.HasWifi
	}
	if updateData.HasAC != nil {
		updates["has_ac"] = *updateData.HasAC
	}
	if updateData.HasTV != nil {
		updates["has_tv"] = *updateData.HasTV
	}
	if updateData.HasMinibar != nil {
		updates["has_minibar"] = *updateData.HasMinibar
	}

	// Actualizar la habitación
	if err := r.db.WithContext(ctx).Model(&room).Updates(updates).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	// Obtener la habitación actualizada
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&room).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return &room, nil
}

func (r *RoomRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Room{})

	if result.Error != nil {
		return utils.ErrDatabaseError
	}

	if result.RowsAffected == 0 {
		return utils.ErrRoomNotFound
	}

	return nil
}
