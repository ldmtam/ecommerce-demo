package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/ldmtam/ecommerce-demo/internal/consumers"
	"github.com/ldmtam/ecommerce-demo/internal/handlers"
	"github.com/ldmtam/ecommerce-demo/internal/models"
	"github.com/ldmtam/ecommerce-demo/internal/repository"
	"github.com/ldmtam/ecommerce-demo/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func initDB(logger *zap.Logger, dsn string) (*gorm.DB, error) {
	logger.Info("Connecting to database...")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to database")

	db.AutoMigrate(models.Product{}, models.CustomerActivity{})

	return db, nil
}

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		logger := utils.NewLogger(viper.GetString("setting.log_path"))

		db, err := initDB(logger, viper.GetString("mysql.dsn"))
		if err != nil {
			panic(err)
		}

		mysqlRepo, err := repository.NewMySQLRepo(logger, db)
		if err != nil {
			panic(err)
		}

		h, err := handlers.New(logger, mysqlRepo)
		if err != nil {
			panic(err)
		}

		activityConsumer, err := consumers.NewActivityConsumer(logger, mysqlRepo)
		if err != nil {
			panic(err)
		}
		if err := activityConsumer.Start(); err != nil {
			panic(err)
		}

		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		router.Use(ginzap.RecoveryWithZap(logger, true))

		v1 := router.Group("/api/v1")
		{
			v1.GET("/ping", h.Ping) // for health check

			v1.POST("/products", h.CreateProduct)
			v1.GET("/products/:id", h.GetProduct)
			v1.GET("/products/seachByName/:name", h.SearchProductByName)
			v1.GET("/customer_activities/:id", h.GetCustomerActivites)
			v1.GET("/customer_activities/:id/actions/:action_type", h.GetCustomerActivitesByAction)
		}

		go func() {
			port := viper.GetInt("setting.port")
			logger.Info("Start listening http...", zap.Int("port", port))
			if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
				panic(err)
			}
		}()

		logger.Info("Starting service...")

		sigs := make(chan os.Signal, 1)
		done := make(chan bool)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		logger.Info("Server is now listening...")

		go func() {
			<-sigs
			activityConsumer.Stop()
			done <- true
		}()

		logger.Info("Ctrl-C to interrupt...")
		<-done
		logger.Info("Exiting...")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
