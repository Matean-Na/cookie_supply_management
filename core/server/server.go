package server

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/core/connect"
	"cookie_supply_management/core/database"
	"cookie_supply_management/internal/repositories"
	"cookie_supply_management/internal/services"
	"cookie_supply_management/pkg/logger"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Server struct {
	Srv  *http.Server
	Test bool
}

func (s *Server) Init() (*gin.Engine, error) {
	conf, err := config.Load()
	if err != nil {
		fmt.Printf("err config.Load() %s\n", err)
		return nil, err
	}

	//postgres
	dbase, err := database.Connect(conf.Database)
	if err != nil {
		fmt.Printf("err db.Connect() %s\n", err)
		return nil, err
	}

	//init entities
	repo := repositories.NewRepository(dbase)
	service := services.NewService(repo, conf)

	//init logger
	defer logger.LogFile.Close()
	if err = logger.Init("logfile.log"); err != nil {
		fmt.Printf("err logger.Init() %s\n", err)
		return nil, err
	}

	connect.DB = dbase

	//routing api
	r := SetupRoutesWithDeps(service)

	return r, nil
}

func (s *Server) Run(r *gin.Engine) {
	if r == nil {
		return
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run.seeder: %s\n", err)
		}
	}()

	s.Srv = srv
}

func CloseDb(db2 *gorm.DB) {
	sql, err := db2.DB()
	if err != nil {
		fmt.Printf("error get sql db %s\n", err)
		return
	}

	err = sql.Close()
	if err != nil {
		fmt.Printf("error sql db close %s\n", err)
		return
	}
}

func (s *Server) CloseAll() {
	d := connect.DB
	CloseDb(d)
}
