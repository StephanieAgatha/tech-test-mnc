package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mnc-test/config"
	"mnc-test/delivery/controller"
	"mnc-test/delivery/middleware"
	"mnc-test/manager"
)

type Server struct {
	um   manager.UsecaseManager
	gin  *gin.Engine
	log  *zap.Logger
	host string
}

// middleware goes here
func (s *Server) InitMiddleware() {
	s.gin.Use(middleware.ZapLogger(s.log))
}

// controller
func (s *Server) InitController() {
	controller.NewUserCredentialController(s.um.UserCredUsecase(), s.gin).Route()
	controller.NewMerchantController(s.um.MerchantUsecase(), s.gin).Route()
	controller.NewTransactionController(s.um.TransactionUsecase(), s.gin).Route()
	controller.NewTransferController(s.um.TransferUsecase(), s.gin).Route()
}

// run server
func (s *Server) Run() {
	s.InitMiddleware()
	s.InitController()
	err := s.gin.Run(s.host)
	if err != nil {
		s.log.Error("Failed to run server", zap.String("error", err.Error()))
	}
}

func NewServer() *Server {
	// zap logger goes here baby
	zapconf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"activity.log"}, //activity from merchants,customer,etc goes here. we don't want to mix the log output to log.log
		ErrorOutputPaths: []string{"activity.log"}, //sameeee
	}
	logger, err := zapconf.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	//define contrusctor from config
	cfg, err := config.NewDbConfig()
	if err != nil {
		logger.Error("Failed on config server", zap.String("error", err.Error()))
	}

	//constructor from infra
	im, err := manager.NewInfraManager(cfg)
	if err != nil {
		logger.Error("Failed on construct infra", zap.String("error", err.Error()))
	}

	//constructor from repomanager
	rm := manager.NewRepoManager(im)
	//contructor from usecase manager
	um := manager.NewUsecaseManager(rm, logger)

	//set host for gin server
	host := fmt.Sprintf("%s:%s", cfg.ApiConfig.Host, cfg.ApiConfig.Port)
	//return gin instance
	g := gin.Default()
	return &Server{
		um:   um,
		gin:  g,
		log:  logger,
		host: host,
	}
}
