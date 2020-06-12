package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"

	"github.com/dxvgef/tsing-center/api"
	"github.com/dxvgef/tsing-center/global"
	"github.com/dxvgef/tsing-center/storage"

	"github.com/dxvgef/tsing"
)

func main() {
	var (
		configFile     string
		err            error
		apiHttpServer  *http.Server
		apiHttpsServer *http.Server
	)

	// 设置默认logger
	setDefaultLogger()

	// --------------------- 加载配置文件 ----------------------
	flag.StringVar(&configFile, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	err = global.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("加载配置文件失败")
		return
	}

	// --------------------- 配置snowflake id ----------------------
	snowflake.Epoch = time.Now().Unix()
	global.SnowflakeNode, err = snowflake.NewNode(int64(time.Now().Hour()))
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("配置snowflake失败")
		return
	}

	// --------------------- 根据配置构建存储器 ----------------------
	global.Storage, err = storage.Build(global.Config.Storage.Name, global.Config.Storage.Config)
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("构建存储器失败")
		return
	}
	// 从存储器中加载所有数据
	if err = global.Storage.LoadAll(); err != nil {
		log.Fatal().Err(err).Caller().Msg("加载数据失败")
		return
	}

	// 监听存储中的数据变更
	go func() {
		log.Info().Msg("开始监听数据变更")
		if err = global.Storage.Watch(); err != nil {
			log.Fatal().Err(err).Caller().Msg("启动存储器监听失败")
			return
		}
	}()

	// 启动API服务
	if global.Config.API.HTTP.Port > 0 || global.Config.API.HTTPS.Port > 0 {
		var (
			apiEngineConfig tsing.Config
			rootPath        string
		)
		apiEngineConfig.EventHandler = api.EventHandler
		apiEngineConfig.Recover = true
		apiEngineConfig.EventShortPath = true
		apiEngineConfig.EventSource = true
		apiEngineConfig.EventTrace = true
		apiEngineConfig.EventHandlerError = true
		rootPath, err = os.Getwd()
		if err == nil {
			apiEngineConfig.RootPath = rootPath
		}
		apiEngine := tsing.New(&apiEngineConfig)
		// 设置路由
		api.SetRouter(apiEngine)
		// 启动api http服务
		if global.Config.API.HTTP.Port > 0 {
			go func() {
				apiHttpServer = &http.Server{
					Addr:              global.Config.API.IP + ":" + strconv.FormatUint(uint64(global.Config.API.HTTP.Port), 10),
					Handler:           apiEngine,
					ReadTimeout:       global.Config.API.ReadTimeout,
					WriteTimeout:      global.Config.API.WriteTimeout,
					IdleTimeout:       global.Config.API.IdleTimeout,
					ReadHeaderTimeout: global.Config.API.ReadHeaderTimeout,
				}
				log.Info().Str("addr", apiHttpServer.Addr).Msg("API HTTP服务")
				if err = apiHttpServer.ListenAndServe(); err != nil {
					if err == http.ErrServerClosed {
						log.Info().Msg("API HTTP服务已关闭")
						return
					}
					log.Fatal().Err(err).Caller().Msg("启动API HTTP服务失败")
					return
				}
			}()
		}

		// 启动api https服务
		if global.Config.API.HTTPS.Port > 0 {
			go func() {
				apiHttpsServer = &http.Server{
					Addr:              global.Config.API.IP + ":" + strconv.FormatUint(uint64(global.Config.API.HTTPS.Port), 10),
					Handler:           apiEngine,
					ReadTimeout:       global.Config.API.ReadTimeout,
					WriteTimeout:      global.Config.API.WriteTimeout,
					IdleTimeout:       global.Config.API.IdleTimeout,
					ReadHeaderTimeout: global.Config.API.ReadHeaderTimeout,
				}
				if global.Config.API.HTTPS.HTTP2 {
					if err = http2.ConfigureServer(apiHttpsServer, &http2.Server{}); err != nil {
						log.Fatal().Err(err).Caller().Msg("启动API HTTP2支持失败")
						return
					}
				}
				log.Info().Bool("HTTP2", global.Config.API.HTTPS.HTTP2).Str("addr", apiHttpsServer.Addr).Msg("API HTTPS服务")
				if err = apiHttpsServer.ListenAndServeTLS("server.cert", "server.key"); err != nil {
					if err == http.ErrServerClosed {
						log.Info().Msg("API HTTPS服务已关闭")
						return
					}
					log.Fatal().Err(err).Caller().Msg("启动API HTTPS服务失败")
					return
				}
			}()
		}
	}

	// 阻塞并等待退出超时
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), global.Config.API.QuitWaitTimeout)
	defer cancel()

	// 关闭API HTTP服务
	if global.Config.API.HTTP.Port > 0 {
		if err := apiHttpServer.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Caller().Msg("启动API HTTP服务失败")
			return
		}
	}
	// 关闭API HTTPS服务
	if global.Config.API.HTTPS.Port > 0 {
		if err := apiHttpsServer.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Caller().Msg("启动API HTTPS服务失败")
			return
		}
	}

	log.Info().Msg("进程已退出")
}
