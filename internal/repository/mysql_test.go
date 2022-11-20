package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ldmtam/ecommerce-demo/internal/models"
	"github.com/ldmtam/ecommerce-demo/internal/repository"
	"github.com/ldmtam/ecommerce-demo/utils"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	repo *repository.MysqlRepo
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "latest", []string{"MYSQL_ROOT_PASSWORD=example"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = gorm.Open(mysql.Open(fmt.Sprintf("root:example@(localhost:%s)/mysql", resource.GetPort("3306/tcp"))), &gorm.Config{})
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		return db.WithContext(ctx).Exec("SELECT 1").Error
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	db.AutoMigrate(models.Product{}, models.CustomerActivity{})
	logger := utils.NewLogger("./logs")

	repo, err = repository.NewMySQLRepo(logger, db)
	if err != nil {
		log.Fatalf("Could not create mysql repo: %v", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.RemoveAll("./logs")

	os.Exit(code)
}

func TestCreateProduct(t *testing.T) {
	type inputStruct struct {
		name  string
		price uint
	}
	tests := map[string]struct {
		input          inputStruct
		expectedOutput uint
		expectedError  error
	}{
		"happy case": {
			input: inputStruct{
				name:  "Ultraboost 2022 shoes",
				price: 300,
			},
			expectedOutput: 1,
			expectedError:  nil,
		},
		"name is empty": {
			input: inputStruct{
				name:  "",
				price: 1000,
			},
			expectedOutput: 0,
			expectedError:  repository.ErrProductNameIsEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out, err := repo.CreateProduct(test.input.name, test.input.price)
			if out != nil {
				assert.EqualValues(t, test.expectedOutput, out.ID)
			}
			assert.EqualValues(t, test.expectedError, err)
		})
	}
}

func TestGetProductByID(t *testing.T) {
	type inputStruct struct {
		name  string
		price uint
	}
	tests := map[string]struct {
		input          inputStruct
		expectedOutput *models.Product
		expectedError  error
	}{
		"happy case": {
			input: inputStruct{
				name:  "Ultraboost 2022 shoes",
				price: 300,
			},
			expectedOutput: &models.Product{
				ID:    1,
				Name:  "Ultraboost 2022 shoes",
				Price: 300,
			},
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out, err := repo.CreateProduct(test.input.name, test.input.price)
			assert.Nil(t, err)

			createdProduct, err := repo.GetProductByID(out.ID)
			assert.EqualValues(t, test.expectedError, err)
			assert.EqualValues(t, test.expectedOutput.ID, createdProduct.ID)
			assert.EqualValues(t, test.expectedOutput.Name, createdProduct.Name)
			assert.EqualValues(t, test.expectedOutput.Price, createdProduct.Price)
		})
	}
}
