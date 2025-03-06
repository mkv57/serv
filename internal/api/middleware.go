package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"serv/internal/logger"
	"strings"

	"github.com/gorilla/mux"
)

var Id int

func Logg(log *slog.Logger, t Server) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			path := r.URL.Path
			fmt.Println(path)
			if path == "/user" {
				fmt.Println("токен не проверяем")
			} else {
				fmt.Println("токен проверяем")
				tok := r.Header.Get("authorization")

				token := strings.TrimPrefix(tok, "Bearer ") // после Bearer ставим пробел, иначе токен будет начинаться с пробела и не будет принят

				//fmt.Println(token)

				id, err := t.Database.ControlTokenDB(ctx, token)
				Id = id
				if err != nil {

					fmt.Println("токен не подошёл", token)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				fmt.Println(id, Id, err)
			}
			log1 := log

			log = log.With(
				slog.String("ip", r.RemoteAddr), // например добавляем в лог ip адрес
				slog.String("url_path", path),   // добавляем в лог точку входа
			)
			ctx = logger.NewContext(ctx, log)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			log = log1

		})
	}

}
