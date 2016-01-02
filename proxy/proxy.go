// Package proxy contains all proxies used by V2Ray.

package proxy

import (
	"errors"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

var (
	inboundFactories  = make(map[string]internal.InboundConnectionHandlerCreator)
	outboundFactories = make(map[string]internal.OutboundConnectionHandlerCreator)

	ErrorProxyNotFound    = errors.New("Proxy not found.")
	ErrorNameExists       = errors.New("Proxy with the same name already exists.")
	ErrorBadConfiguration = errors.New("Bad proxy configuration.")
)

func RegisterInboundConnectionHandlerFactory(name string, creator internal.InboundConnectionHandlerCreator) error {
	if _, found := inboundFactories[name]; found {
		return ErrorNameExists
	}
	inboundFactories[name] = creator
	return nil
}

func RegisterOutboundConnectionHandlerFactory(name string, creator internal.OutboundConnectionHandlerCreator) error {
	if _, found := outboundFactories[name]; found {
		return ErrorNameExists
	}
	outboundFactories[name] = creator
	return nil
}

func CreateInboundConnectionHandler(name string, space app.Space, rawConfig []byte) (connhandler.InboundConnectionHandler, error) {
	creator, found := inboundFactories[name]
	if !found {
		return nil, ErrorProxyNotFound
	}
	if len(rawConfig) > 0 {
		proxyConfig, err := config.CreateInboundConnectionConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator(space, proxyConfig)
	}
	return creator(space, nil)
}

func CreateOutboundConnectionHandler(name string, space app.Space, rawConfig []byte) (connhandler.OutboundConnectionHandler, error) {
	creator, found := outboundFactories[name]
	if !found {
		return nil, ErrorNameExists
	}

	if len(rawConfig) > 0 {
		proxyConfig, err := config.CreateOutboundConnectionConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator(space, proxyConfig)
	}

	return creator(space, nil)
}
