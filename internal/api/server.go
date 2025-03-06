package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"serv/internal/domain"
	"serv/internal/logger"
	"strconv"

	"github.com/google/uuid"
)

type Store interface {

	// logica
	ControlTokenDB(ctx context.Context, token string) (int, error)
	DeleteDB(ctx context.Context, c int)
	GetBookDB(cxt context.Context, id int) domain.Books
	GetBooksDB(ctx context.Context, n int) ([]domain.Books, error)
	PasswordDB(ctx context.Context, email string) (string, int, error)
	SaveBook(ctx context.Context, c domain.Books, r int) domain.Books
	Session(ctx context.Context, s domain.Session) error
	UpdateDb(ctx context.Context, id int, c domain.Books) domain.Books
	UserDB(ctx context.Context, user domain.User) (int, error)
}

type Server struct {
	//Database *db.Repository `json:"database"`
	Database Store
}

var (
	ErrInvalidPassword = errors.New("invalid password")
)

//const authScheme = "Bearer"

func (p Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jsong, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("problem LoginUser")
	}
	var user domain.User
	err = json.Unmarshal(jsong, &user)
	if err != nil {
		fmt.Println("problem AddUserUnmarshal")
	}
	password, id, err := p.Database.PasswordDB(ctx, user.Email)
	if err != nil {
		fmt.Println("problem AddUserDB")
	}
	resp, err := json.Marshal(id)
	if err != nil {
		fmt.Println("problem Marshal", resp)
	}

	if user.Password != password {
		fmt.Println(ErrInvalidPassword)
		fmt.Println(user.Email, "пороль error")
	} else {
		user_agent := r.UserAgent() // это  будет (PostmanRuntime/7.43.0)
		//g1 := r.URL		//  это  будет точка входа (/user)
		g2 := r.URL.User   // пусто?
		ip := r.RemoteAddr //  ip адрес
		//g := r.URL.Path    // это  будет точка входа (/user)
		fmt.Println(g2)
		token := uuid.New().String()
		session := domain.Session{
			UserID:     id,
			Token:      token,
			Ip:         ip,
			User_Agent: user_agent,
		}
		err = p.Database.Session(ctx, session)
		//fmt.Println(session)
		if err != nil {
			fmt.Println("problem Session")
		}

		w.Header().Set("token", token)
		// w.Header().Add("token1", token)
		log, found := logger.FromContext(ctx)
		if !found {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("problem loginUser log")
			return
		}
		log.Debug("токен отправлен")

		_, err := w.Write(resp)
		if err != nil {
			fmt.Println("problem Session")
		}
	}
}

func (p Server) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log, found := logger.FromContext(ctx)
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsong, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("problem AddUser")
	}
	var user domain.User
	err = json.Unmarshal(jsong, &user)
	if err != nil {
		fmt.Println("problem JsonAddUser")
	}
	id, err := p.Database.UserDB(ctx, user)
	if err != nil {
		fmt.Println("problem AddUserUserDB")
	}
	IdUser, err := json.Marshal(id)
	if err != nil {
		fmt.Println("problem Marshal AddUser")
	}
	log.Debug("user добавлен")
	_, err = w.Write(IdUser)
	if err != nil {
		fmt.Println("problem Session")
	}
}

func (p Server) AddBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	/*tok := r.Header.Get("authorization")

	token := strings.TrimPrefix(tok, "Bearer ") // после Bearer ставим пробел, иначе токен будет начинаться с пробела и не будет принят

	fmt.Println(token)

	UserID, err := p.Database.ControlTokenDB(ctx, token)
	if err != nil {
		fmt.Errorf("токен не подошёл")
	} else {*/
	//ooo := Id
	//fmt.Println(ooo)
	jsong, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("problem")
	}
	var newBook domain.Books
	err = json.Unmarshal(jsong, &newBook)
	if err != nil {
		fmt.Println("problem1")
	}
	save := p.Database.SaveBook(ctx, newBook, Id) //UserID)

	f, err := json.Marshal(save)
	//t, err := json.Marshal(id)
	if err != nil {
		fmt.Println("problem2")
	}
	log, found := logger.FromContext(ctx)
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Проблемы у нас")
		return
	}
	log.Debug("книга добавлена")
	_, err = w.Write(f)
	if err != nil {
		fmt.Println("problem Session")
	}
}

func (p Server) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log, found := logger.FromContext(ctx)
	if !found {
		handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		return
	}
	/*tok := r.Header.Get("authorization")

	token := strings.TrimPrefix(tok, "Bearer ") // после Bearer ставим пробел, иначе токен будет начинаться с пробела и не будет принят

	fmt.Println(token)

	_, err := p.Database.ControlTokenDB(token)
	if err != nil {
		fmt.Errorf("токен не подошёл")
	} else {*/
	idstring := r.URL.Query().Get("id")
	int, err := strconv.Atoi(idstring)
	if err != nil {
		fmt.Println("problem3")
	}
	book := p.Database.GetBookDB(ctx, int)

	getbook, err := json.Marshal(book)
	if err != nil {
		fmt.Println("problem4")
	}
	log.Debug("ответ с книгой отправлен")
	_, err = w.Write(getbook)
	if err != nil {
		fmt.Println("problem Session")
	}
}

func (p Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	/*tok := r.Header.Get("authorization")

	token := strings.TrimPrefix(tok, "Bearer ") // после Bearer ставим пробел, иначе токен будет начинаться с пробела и не будет принят

	fmt.Println(token)

	_, err := p.Database.ControlTokenDB(token)
	if err != nil {
		fmt.Errorf("токен %s не подошёл", token)
	} else {*/
	idstring := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstring)
	if err != nil {
		fmt.Println("problem5")
	}
	p.Database.DeleteDB(ctx, idint)

	result := "книга удалена"
	resp, err := json.Marshal(result)
	if err != nil {
		fmt.Println("problem7")
	}
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("problem Session")
	}
}

func (p Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	/*tok := r.Header.Get("authorization")

	token := strings.TrimPrefix(tok, "Bearer ") // после Bearer ставим пробел, иначе токен будет начинаться с пробела и не будет принят

	fmt.Println(token)

	UserID, err := p.Database.ControlTokenDB(token)
	if err != nil {
		fmt.Errorf("токен не подошёл")
	} else {*/
	idstring := r.URL.Query().Get("id")
	idint, err := strconv.Atoi(idstring)
	if err != nil {
		fmt.Println("problem5")
	}
	jsong, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("problem8")
	}
	var book domain.Books
	err = json.Unmarshal(jsong, &book)
	if err != nil {
		fmt.Println("problem5")
	}
	book.UserId = idint //UserID
	newBook := p.Database.UpdateDb(ctx, idint, book)

	resp, err := json.Marshal(newBook)
	if err != nil {
		fmt.Println("problem8")
	}
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("problem Session")
	}
}

func (p Server) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log, found := logger.FromContext(ctx)
	if !found {
		handleError(w, http.StatusInternalServerError, errors.New("Проблемы у нас"))
		return
	}
	/*tok := r.Header.Get("authorization")

	token := strings.TrimPrefix(tok, "Bearer ") // после Bearer ставим пробел, иначе токен будет начинаться с пробела и не будет принят

	fmt.Println(token)

	_, err := p.Database.ControlTokenDB(token)
	if err != nil {
		fmt.Errorf("токен не подошёл")
	} else {*/
	query := r.URL.Query() // можно одной строкой query := r.URL.Query().Get("limit")
	limit := query.Get("limit")
	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println("лимит не указан")
	}

	books, err := p.Database.GetBooksDB(ctx, limit_int)
	if err != nil {
		fmt.Println("problem9")
	}
	resp, err := json.Marshal(books)
	if err != nil {
		fmt.Println("problem10")
	}
	log.Debug("ответ с книгами отправлен")
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("problem Session")
	}

}
