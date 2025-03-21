package middleware

import (
	"context"
	"net/http"
	"strconv"
)

// PaginationParams — структура для хранения параметров пагинации
type PaginationParams struct {
	Limit  int
	Offset int
}

// PaginationMiddleware — middleware для обработки параметров пагинации
func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Значения по умолчанию
		offset := 0
		limit := 10

		// Читаем `page` из запроса
		if o := r.URL.Query().Get("offset"); o != "" {
			if pInt, err := strconv.Atoi(o); err == nil && pInt > 0 {
				offset = pInt
			}
		}

		// Читаем `limit` из запроса
		if l := r.URL.Query().Get("limit"); l != "" {
			if lInt, err := strconv.Atoi(l); err == nil && lInt > 0 {
				limit = lInt
			}
		}

		// Добавляем параметры в контекст запроса
		ctx := context.WithValue(r.Context(), "pagination", &PaginationParams{Offset: offset, Limit: limit})

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPaginationParamsFromCtx(ctx context.Context) (*PaginationParams, bool) {
	p, ok := ctx.Value("pagination").(*PaginationParams)
	return p, ok
}
