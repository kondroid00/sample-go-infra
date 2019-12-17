package alert

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"io"
	"sync"
	"time"
)

type Config struct {
	Dsn string
}

// sentryに送信するClientを指定するキー
type ClientKey string

var (
	FlushTimeout = 10 * time.Second
)

const (
	clientKeyDefault ClientKey = "default"
)

type (
	Sentry struct {
		options   sentry.ClientOptions
		debugMode bool
		mutex     sync.RWMutex
		// key: 同一クライアントを選択するためのclientKey, value: *sentry.Client
		clientPool map[ClientKey]*sentry.Client
	}
)

func GetSentry() *Sentry {
	return instance
}

var instance *Sentry

func Init(config *Config, serverName string, debugMode bool) error {
	var debugLogger io.Writer
	if debugMode {
		debugLogger = &DebugLogger{}
	}

	options := sentry.ClientOptions{
		Dsn:         config.Dsn,
		Debug:       debugMode,
		DebugWriter: debugLogger,
		ServerName:  serverName,
	}

	instance = &Sentry{
		options:    options,
		debugMode:  debugMode,
		clientPool: make(map[ClientKey]*sentry.Client, 0),
	}

	return sentry.Init(options)
}

func Flush() error {
	if !sentry.Flush(FlushTimeout) {
		return &FlushTimeoutError{}
	}
	return nil
}

func (s *Sentry) Send(e error, info *Info) {
	if s.debugMode {
		// デバッグモードの時はsentryに送信しない。デバッグモードで送信したい時はコメントアウト。
		return
	}

	clientKey := clientKeyDefault
	var target AlertError
	if errors.As(e, &target) {
		if value := target.ClientKey(); value != "" {
			clientKey = value
		}
	}

	client, err := s.getClient(clientKey)
	if err != nil {
		return
	}

	scope := sentry.NewScope()
	if info != nil {
		if info.UserInfo != nil {
			scope.SetUser(sentry.User{
				Email:     info.UserInfo.Email,
				ID:        info.UserInfo.ID,
				IPAddress: info.UserInfo.IPAddress,
				Username:  info.UserInfo.Username,
			})
		}
		if info.Request != nil {
			scope.SetRequest(sentry.Request{
				URL:         info.Request.URL,
				Method:      info.Request.Method,
				Data:        info.Request.Data,
				QueryString: info.Request.QueryString,
				Cookies:     info.Request.Cookies,
				Headers:     info.Request.Headers,
				Env:         info.Request.Env,
			})
		}
		for k, v := range info.Tag {
			scope.SetTag(k, v)
		}
		for k, v := range info.Context {
			scope.SetContext(k, v)
		}
	}

	sentry.NewHub(client, scope).CaptureException(e)
}

func (s *Sentry) getClient(clientKey ClientKey) (*sentry.Client, error) {
	if c, ok := func() (*sentry.Client, bool) {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
		c, ok := s.clientPool[clientKey]
		return c, ok
	}(); ok {
		// poolにクライアントが存在すればそのクライアントを返却する
		return c, nil
	}

	// poolにクライアントがなければ新規でクライアントを作成してpoolに登録する
	client, err := sentry.NewClient(s.options)
	if err != nil {
		return nil, err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 複数のgoroutineがmutex.Lockでブロックされていた場合にclientが書き換わってしまうのを防ぐためpoolにclientが存在するかどうかを確認
	if c, ok := s.clientPool[clientKey]; ok {
		return c, nil
	}
	s.clientPool[clientKey] = client

	return client, nil
}
