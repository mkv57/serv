package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"serv/internal/domain"
	"serv/internal/logger"
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestServer_GetBook(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStore := NewMockStore(ctrl)

	server := Server{
		Database: mockStore,
	}

	expectedBook := &domain.Books{
		Id:     0,
		Title:  "book",
		Year:   2025,
		UserId: 0,
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ // создаём переменную log для добавления в context
		Level: slog.LevelDebug,
	}))

	ctx := context.Background()       // создаём переменную типа context с пустым значением Background
	ctx = logger.NewContext(ctx, log) // вкладываем в context 2е переменные ctx и log , как и в функции GetBookFromDatabaseByRAWSql
	fmt.Println("1")
	mockStore.EXPECT().GetBookDB(ctx, 0).Return(*expectedBook)
	fmt.Println("2")
	reg, err := http.NewRequest(http.MethodGet, "http://localhost:8086/book?id=0", nil) // оправляем запрос на обработчик GetBook,
	// который вызывает функцию GetBookFromDatabaseByRAWSql(ctx, uint(idint)), он достаёт  данные из структуры Book и эти данные отправляются в Request reg
	fmt.Println("3")
	require.NoError(t, err) // проверка используется вместо if t!= nil  но только для тестов
	fmt.Println("4")
	reg = reg.WithContext(ctx) // добавляем в запрос в context то же значение ctx
	fmt.Println("5")
	resp := httptest.NewRecorder() // ??? создаём переменную типа ResponseWrite, чтобы вызвать GetBook, вложив resp
	fmt.Println("6", resp.Body, "q ", reg, "p ", reg)

	server.GetBook(resp, reg) // вызываем функцию, которая должна записать в тело ответа resp.Body  json c данными структуры Book
	fmt.Println("7", resp.Body, reg)
	result := &domain.Books{} // создаём переменную типа структуры Book

	err = json.Unmarshal(resp.Body.Bytes(), &result) // resp.Body.Bytes() что это ???	распарсили json и вложили данные в result
	require.NoError(t, err)

	//expectedBook.CreatedAt = result.CreatedAt // уровняли зачения времени, так как будут разные значения
	//expectedBook.UpdatedAt = result.UpdatedAt // уровняли зачения времени

	require.Equal(t, expectedBook, result) // сравнили результаты

}
