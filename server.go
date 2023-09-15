package main

import (
	"errors"
	"fmt"
	"net/http"
	//"rsa"

	"github.com/lxzan/gws"
)

type ServerGwsHandler struct {
	gws.BuiltinEventHandler
	downloadChannel chan unionTransInfor
	uploadChannel   chan unionTransInfor
}

func serverInit(serverPort string, downloadChannel, uploadChannel chan unionTransInfor) error {
	// 服务器初始化
	serverHandler := ServerGwsHandler{
		downloadChannel: downloadChannel,
		uploadChannel:   uploadChannel,
	}

	tokenList, getTokenError := fileOperater("./token.txt", fileOperaterOptions{operater: "read"})
	if getTokenError != nil {
		return errors.New("get token error: " + getTokenError.Error())
	}
	token := tokenList[0]


	// 用RSA写验证，建立连接后再进行websocket升格
	func checkToken (tokenLikeString string) bool {
		return true
	}

	// 用RSA写验证，建立连接后再进行websocket升格
	server := gws.NewServer(serverHandler, &gws.ServerOption{
		Authorize: func(r *http.Request, session gws.SessionStorage) bool {
			if value, ok := r.Header["acstokn"]; !ok {
				return false
			} else {
				if len(value) != 2 {
					return false
				} else {
					if checkToken(value[1]) {
						return false
					}
				}
			}
			return true
		},
	})
	serverStartError := server.Run(serverPort)
	if serverStartError != nil {
		return errors.New("server start error: " + serverStartError.Error())
	}
	fmt.Println("server start")
	return nil
}
