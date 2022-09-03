package server

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"valuation/api/valuation/v1"
	"valuation/internal/conf"
	"valuation/internal/service"
	"valuation/pkg/aop"
	"valuation/pkg/errorx"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, good *service.GoodService, logger log.Logger) *kratosHttp.Server {
	var opts = []kratosHttp.ServerOption{
		kratosHttp.ResponseEncoder(encodeResponseFunc),
		kratosHttp.ErrorEncoder(errorEncoder),
		kratosHttp.Middleware(
			aop.Recovery(),
			aop.Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratosHttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratosHttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratosHttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := kratosHttp.NewServer(opts...)
	srv.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		mySigningKey := []byte("AllYourBase")
		o := &struct {
			signingMethod jwt.SigningMethod
			claims        jwt.Claims
		}{
			signingMethod: jwt.SigningMethodHS256,
			claims: jwt.RegisteredClaims{
				ID: "1",
			},
		}
		token := jwt.NewWithClaims(o.signingMethod, o.claims)
		tokenString, err := token.SignedString(mySigningKey)
		b, err := json.Marshal(map[string]interface{}{
			"msg": "ok",
			"data": map[string]interface{}{
				"access_token": tokenString,
			},
			"code": 200,
		})
		_, err = writer.Write(b)
		if err != nil {
			return
		}
	})
	valuation.RegisterGoodHTTPServer(srv, good)
	return srv
}

func encodeResponseFunc(w http.ResponseWriter, r *http.Request, v interface{}) error {
	type response struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
		Msg  string      `json:"msg"`
	}
	reply := &response{
		Code: 200,
		Data: v,
		Msg:  "ok",
	}
	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func errorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errorx.FromError(err)
	codec, _ := kratosHttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	w.WriteHeader(se.Code)
	_, _ = w.Write(body)
}
