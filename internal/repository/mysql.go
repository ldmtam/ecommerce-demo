package repository

import (
	"errors"
	"time"

	"github.com/ldmtam/ecommerce-demo/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrProductNameIsEmpty = errors.New("product name is empty")
)

type MysqlRepo struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewMySQLRepo(logger *zap.Logger, db *gorm.DB) (*MysqlRepo, error) {
	return &MysqlRepo{
		logger: logger,
		db:     db,
	}, nil
}

func (repo *MysqlRepo) CreateProduct(name string, price uint) (*models.Product, error) {
	if name == "" {
		return nil, ErrProductNameIsEmpty
	}

	product := &models.Product{
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().UnixMilli(),
	}

	if err := repo.db.Create(product).Error; err != nil {
		repo.logger.Error("Insert new product to database failed", zap.Error(err))
		return nil, err
	}

	return product, nil
}

func (repo *MysqlRepo) GetProductByID(id uint) (*models.Product, error) {
	product := &models.Product{ID: id}

	if err := repo.db.First(product).Error; err != nil {
		repo.logger.Error("Get product from database failed", zap.Error(err))
		return nil, err
	}

	return product, nil
}

func (repo *MysqlRepo) GetProductByName(name string, limit uint) ([]*models.Product, error) {
	query := `
		SELECT *, MATCH (name) AGAINST (?) as score FROM products
		WHERE MATCH (name) AGAINST (?)
		ORDER BY score DESC
		LIMIT ?;
	`
	var products []*models.Product
	if err := repo.db.Raw(query, name, name, limit).Scan(&products).Error; err != nil {
		repo.logger.Error("Get product by name from database failed", zap.Error(err))
		return nil, err
	}

	return products, nil
}

func (repo *MysqlRepo) CreateCustomerActivity(userID uint, createdAt int64, action, data string) (*models.CustomerActivity, error) {
	customerActivity := &models.CustomerActivity{
		UserID:    userID,
		CreatedAt: createdAt,
		Action:    action,
		Data:      data,
	}

	if err := repo.db.Create(customerActivity).Error; err != nil {
		repo.logger.Error("Insert new customer activity to database failed", zap.Error(err))
		return nil, err
	}

	return customerActivity, nil
}

func (repo *MysqlRepo) GetCustomerActivities(id uint, limit uint) ([]*models.CustomerActivity, error) {
	var customerActivities []*models.CustomerActivity

	if err := repo.db.Where("user_id = ?", id).
		Order("created_at DESC").
		Limit(int(limit)).
		Find(&customerActivities).Error; err != nil {
		repo.logger.Error("Get customer activities failed", zap.Error(err))
		return nil, err
	}

	return customerActivities, nil
}

func (repo *MysqlRepo) GetCustomerActivitiesByAction(id uint, action string, limit uint) ([]*models.CustomerActivity, error) {
	var customerActivities []*models.CustomerActivity

	if err := repo.db.Where("user_id = ? AND action = ?", id, action).
		Order("created_at DESC").
		Limit(int(limit)).
		Find(&customerActivities).Error; err != nil {
		repo.logger.Error("Get customer activities by action failed", zap.Error(err))
		return nil, err
	}

	return customerActivities, nil
}
