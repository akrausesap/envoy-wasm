// Copyright 2020 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type rootContextConfig struct {
	HeaderName  string
	HeaderValue string
}

type httpHeaders struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultHttpContext
	contextID uint32
	config    rootContextConfig
}

type rootContext struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultRootContext
	config rootContextConfig
}

func main() {
	proxywasm.SetNewRootContext(newRootContext)
	proxywasm.SetNewHttpContext(newHTTPContext)
}

func newRootContext(contextID uint32) proxywasm.RootContext {
	return &rootContext{}
}

func newHTTPContext(rootContextID, contextID uint32) proxywasm.HttpContext {
	ctx := &httpHeaders{}

	rootCtx, err := proxywasm.GetRootContextByID(rootContextID)
	if err != nil {
		proxywasm.LogErrorf("unable to get root context: %v", err)

		return ctx
	}

	receivedRootCtx, ok := rootCtx.(*rootContext)
	if !ok {
		proxywasm.LogError("could not cast root context")
	}

	ctx.config = receivedRootCtx.config

	proxywasm.LogInfof("plugin config from root context: %v\n", ctx.config)
	return ctx
}

//override
func (ctx *rootContext) OnPluginStart(pluginConfigurationSize int) bool {
	data, err := proxywasm.GetPluginConfiguration(pluginConfigurationSize)
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
		return false
	}

	configString := string(data)
	configStringArray := strings.Split(configString, ":")
	if len(configStringArray) != 2 {
		proxywasm.LogCriticalf("error extracting plugin configuration from %s. format headerName: headerValue expected", configString)
		return false
	}

	ctx.config = rootContextConfig{
		HeaderName:  strings.TrimSpace(configStringArray[0]),
		HeaderValue: strings.TrimSpace(configStringArray[1]),
	}

	return true
}

// override
func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {

	err := proxywasm.SetHttpResponseHeader(ctx.config.HeaderName, ctx.config.HeaderValue)
	if err != nil {
		proxywasm.LogCriticalf("failed to set response headers: %v", err)
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogInfof("request header --> %s: %s", h[0], h[1])
	}

	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
