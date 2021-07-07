package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"demo_casbin/api"
	"demo_casbin/bll"
	"demo_casbin/casbin"
	"demo_casbin/config"
	"demo_casbin/middleware"
	"demo_casbin/model/gorm"
	"demo_casbin/model/gorm/model"
	"demo_casbin/pkg/logger"
	"demo_casbin/router"
)

type options struct {
	ConfigFile string
	ModelFile  string
	MenuFile   string
	WWWDir     string
	Version    string
}

// Option 定义配置项
type Option func(*options)

// SetConfigFile 设定配置文件
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetModelFile 设定casbin模型配置文件
func SetModelFile(s string) Option {
	return func(o *options) {
		o.ModelFile = s
	}
}

// SetWWWDir 设定静态站点目录
func SetWWWDir(s string) Option {
	return func(o *options) {
		o.WWWDir = s
	}
}

// SetMenuFile 设定菜单数据文件
func SetMenuFile(s string) Option {
	return func(o *options) {
		o.MenuFile = s
	}
}

// SetVersion 设定版本号
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

// Init 应用初始化
func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	config.MustLoad(o.ConfigFile)
	if v := o.ModelFile; v != "" {
		config.C.Casbin.Model = v
	}
	if v := o.WWWDir; v != "" {
		config.C.WWW = v
	}
	if v := o.MenuFile; v != "" {
		config.C.Menu.Data = v
	}
	config.PrintWithJSON()

	logger.WithContext(ctx).Printf("服务启动，运行模式：%s，版本号：%s，进程号：%d", config.C.RunMode, o.Version, os.Getpid())

	db, cleanup, err := gorm.InitGormDB()
	if err != nil {
		return cleanup, err
	}
	user := &model.User{
		DB: db,
	}
	project := &model.Project{
		DB: db,
	}
	projectUser := &model.ProjectUser{
		DB: db,
	}
	policySource := &model.PolicySource{
		DB: db,
	}
	casbinAdapter := &casbin.CasbinAdapter{
		UserModel:         user,
		ProjectUserModel:  projectUser,
		ProjectModel:      project,
		PolicySourceModel: policySource,
	}

	syncedEnforcer, cleanup2, err := casbin.InitCasbin(casbinAdapter)
	if err != nil {
		return func() {
			cleanup()
			cleanup2()
		}, err
	}

	trans := &model.Trans{
		DB: db,
	}

	bllUser := &bll.User{
		Enforcer:   syncedEnforcer,
		TransModel: trans,
		UserModel:  user,
	}

	apiUser := &api.User{
		UserBll: bllUser,
	}

	login := &bll.Login{
		UserModel: user,
	}
	apiLogin := &api.Login{
		LoginBll: login,
	}

	bllProject := &bll.Project{
		Enforcer:         syncedEnforcer,
		TransModel:       trans,
		ProjectModel:     project,
		ProjectUserModel: projectUser,
	}

	apiProject := &api.Project{
		ProjectBll: bllProject,
	}
	routerRouter := &router.Router{
		CasbinEnforcer: syncedEnforcer,
		LoginAPI:       apiLogin,
		ProjectAPI:     apiProject,
		UserAPI:        apiUser,
	}
	engine := InitGinEngine(routerRouter)

	// 初始化菜单数据
	policyM := model.PolicySource{
		DB: db,
	}
	bllPolicySource := bll.PolicySource{
		PolicySourceModel: policyM,
	}

	err = bllPolicySource.InitData(ctx, "D:/Mark/go/src/demo_casbin/config/policy_source.yaml")
	if err != nil {
		return nil, err
	}

	// 初始化HTTP服务
	httpServerCleanFunc := InitHTTPServer(ctx, engine)
	httpDebugCleanFunc := InitDebugSever(ctx, engine)

	return func() {
		cleanup()
		cleanup2()
		httpServerCleanFunc()
		httpDebugCleanFunc()
	}, nil
}

// InitHTTPServer 初始化http服务
func InitHTTPServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.WithContext(ctx).Printf("HTTP server is running at %s.", addr)

		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Errorf("主服务退出: " + err.Error())
		}
	}
}

func InitDebugSever(ctx context.Context, handler http.Handler) func() {

	srv := &http.Server{
		Addr: ":6060",
	}
	srv.ListenAndServe()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*60)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Errorf("pprof 退出", err.Error())
		}
	}

}

// Run 运行服务
func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("接收到信号[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.WithContext(ctx).Infof("服务退出")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}

// InitGinEngine 初始化gin引擎
func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	// Router register
	r.Register(app)

	return app
}

func main() {
	ctx := logger.NewTagContext(context.Background(), "__main__")
	err := Run(ctx, SetConfigFile("/config/config.toml"))
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}
