package application

import (
	"github.com/atreugo/session"
	"github.com/atreugo/session/providers/memory"
	"github.com/atreugo/session/providers/mysql"
	"github.com/atreugo/session/providers/redis"
	session2 "github.com/fasthttp/session/v2"
	"github.com/orangebees/go-oneutils/Random"
	"github.com/savsgio/atreugo/v11"
	"os"
	"strconv"
	"time"
)

var atreugoSession *session.Session

func GetAtreugoSession() *session.Session {
	if atreugoSession != nil {
		return atreugoSession
	}
	return createAtreugoSession()
}

type sqlConfig struct {
	Host      string
	Port      string
	UserName  string
	Password  string
	DbName    string
	TableName string
}

func createAtreugoSession() *session.Session {
	var provider session.Provider
	var err error
	var sc sqlConfig
	sc.Host = "127.0.0.1"
	if v := os.Getenv("ATREUGO_SESSION_HOST"); v != "" {
		sc.Host = v
	}
	if v := os.Getenv("ATREUGO_SESSION_PORT"); v != "" {
		sc.Port = v
	}
	if v := os.Getenv("ATREUGO_SESSION_USERNAME"); v != "" {
		sc.UserName = v
	}
	if v := os.Getenv("ATREUGO_SESSION_PASSWORD"); v != "" {
		sc.Password = v
	}
	if v := os.Getenv("ATREUGO_SESSION_DBNAME"); v != "" {
		sc.DbName = v
	}
	if v := os.Getenv("ATREUGO_SESSION_TABLENAME"); v != "" {
		sc.TableName = v
	}
	encoder := session.MSGPEncode
	decoder := session.MSGPDecode
	var defaultProvider = os.Getenv("ATREUGO_SESSION_PROVIDER")
	switch defaultProvider {
	case "memory":
		provider, err = memory.New(memory.Config{})
	case "redis":
		if sc.Port == "" {
			sc.Port = "6379"
		}
		provider, err = redis.New(redis.Config{
			KeyPrefix:   sc.TableName,
			Addr:        sc.Host + ":" + sc.Port,
			Username:    sc.UserName,
			Password:    sc.Password,
			PoolSize:    8,
			IdleTimeout: 30 * time.Second,
		})
	case "mysql":
		encoder = session.Base64Encode
		decoder = session.Base64Decode
		port, err2 := strconv.ParseInt(sc.Port, 10, 64)
		if err2 != nil || (port < 0 || port > 65535) {
			port = 3306
		}
		cfg := mysql.NewConfigWith(sc.Host, 3306, sc.UserName, sc.Password, sc.DbName, sc.TableName)
		provider, err = mysql.New(cfg)
	default:
		provider, err = memory.New(memory.Config{})
	}

	if err != nil {
		panic(err)
	}

	cfg := session.NewDefaultConfig()
	cfg.EncodeFunc = encoder
	cfg.DecodeFunc = decoder
	cfg.SessionIDGeneratorFunc = Random.RandBytes32
	as := session.New(cfg)
	if err = as.SetProvider(provider); err != nil {
		panic(err)
	}

	return as
}

// AutoLoadSaveSessionStore 自动加载保存会话存储
func AutoLoadSaveSessionStore(a *atreugo.Router) {
	a.UseBefore(loadSessionStore).UseAfter(saveSessionStore)
}

// LoadSessionStore 加载会话存储
func loadSessionStore(ctx *atreugo.RequestCtx) error {
	store, err := GetAtreugoSession().Get(ctx)
	if err != nil {
		log.Err(err).Send()
		return nil
	}
	ctx.SetUserValue("store", store)
	return ctx.Next()
}
func saveSessionStore(ctx *atreugo.RequestCtx) error {
	storeI := ctx.UserValue("store")
	if storeI == nil {
		return ctx.Next()
	}
	err := GetAtreugoSession().Save(ctx, storeI.(*session2.Store))
	if err != nil {
		GetLogger().Err(err).Send()
		return nil
	}
	return ctx.Next()
}
