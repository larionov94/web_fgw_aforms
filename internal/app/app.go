package app

import (
	"context"
	"fgw_web_aforms/internal/config"
	"fgw_web_aforms/internal/config/db"
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/handler/http_web"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr = ":7777"
const fileEnv = ".env"

func StartApp() {
	config.InitSessionStore()

	logger, err := common.NewLogger("")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	authMiddleware := handler.NewAuthMiddleware(config.Store, logger)

	configDB, err := config.NewMSSQLCfg(logger, fileEnv)
	if err != nil {
		logger.LogE(msg.E3000, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mssqlDB, err := db.NewConnMSSQL(ctx, configDB, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(mssqlDB)

	repoRole := repository.NewRoleRepo(mssqlDB, logger)
	serviceRole := service.NewRoleService(repoRole, logger)

	repoPerformer := repository.NewPerformerRepo(mssqlDB, logger)
	servicePerformer := service.NewPerformerService(repoPerformer, logger)

	handlerAuthHTML := http_web.NewAuthHandlerHTML(servicePerformer, serviceRole, logger, authMiddleware)

	mux := http.NewServeMux()

	handlerAuthHTML.ServerHTTPRouter(mux)

	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))

	server := config.NewServer(addr, mux, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = server.StartServer(ctx); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()
	<-quit
	log.Println("Получен сигнал остановки сервера...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.LogE(msg.E3102, err)
	}
}
