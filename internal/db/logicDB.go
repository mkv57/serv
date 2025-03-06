package db

import (
	"context"
	"database/sql"
	"fmt"
	"serv/internal/domain"
	"serv/internal/logger"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(rawDB *sql.DB) *Repository {
	return &Repository{db: rawDB}
}

func (p Repository) SaveBook(ctx context.Context, c domain.Books, r int) domain.Books {
	book := &domain.Books{}
	qvery := "INSERT INTO books (title, year, user_id) VALUES ($1, $2, $3) RETURNING id, title, year, user_id"
	err := p.db.QueryRow(qvery, c.Title, c.Year, r).Scan(&book.Id, &book.Title, &book.Year, &book.UserId)
	if err != nil {
		fmt.Println("error при добавлении книги", err)
	}
	return *book
}
func (p Repository) GetBookDB(ctx context.Context, c int) domain.Books {
	book := &domain.Books{}
	qvery := "SELECT id, title, year FROM books WHERE id = $1"
	err := p.db.QueryRow(qvery, c).Scan(&book.Id, &book.Title, &book.Year)
	if err != nil {
		fmt.Println("error такой книги нет", err)
	}
	return *book
}
func (p Repository) DeleteDB(ctx context.Context, c int) {
	_, err := p.db.Exec("DELETE FROM books WHERE id = $1", c)
	if err != nil {
		fmt.Println("книга удалена")
	}
}
func (p Repository) UpdateDb(ctx context.Context, id int, c domain.Books) domain.Books {
	var newbook domain.Books

	query := "UPDATE books SET title = $1, year = $2, user_id = $3 WHERE  id = $4"
	_, err := p.db.Exec(query, c.Title, c.Year, c.UserId, id)
	if err != nil {
		fmt.Println("книга обновлена")
	}
	query1 := "SELECT id, title, year FROM books WHERE id = $1"
	err = p.db.QueryRow(query1, id).Scan(&newbook.Id, &newbook.Title, &newbook.Year)
	if err != nil {
		fmt.Println("книга обновлена", newbook, c.Id)
	}
	return newbook
}
func (p Repository) GetBooksDB(ctx context.Context, n int) ([]domain.Books, error) {
	log, found := logger.FromContext(ctx)
	if !found {
		log.Debug("problem GetBookDB log")
	}
	books := []domain.Books{}
	query := "SELECT id, title, year FROM books"
	resp, err := p.db.Query(query)
	if err != nil {
		fmt.Println("problem DB")
	}
	if n > 0 {
		for i := 1; i <= n; resp.Next() {
			i++
			book := &domain.Books{}
			err = resp.Scan(&book.Id, &book.Title, &book.Year)
			if err != nil {
				fmt.Println("problem DB")
			}
			books = append(books, *book)
		}
	} else {
		for resp.Next() {
			book := &domain.Books{}
			err = resp.Scan(&book.Id, &book.Title, &book.Year)
			if err != nil {
				fmt.Println("problem Session")
			}
			books = append(books, *book)
		}
	}
	return books, nil
}
func (p Repository) Session(ctx context.Context, s domain.Session) error {

	query := "insert into session (user_id, token, ip, user_agent) values ($1, $2, $3, $4)"
	_, err := (p.db.Exec(query, s.UserID, s.Token, s.Ip, s.User_Agent))
	if err != nil {
		fmt.Println("problem sessionDB")
	}
	return nil
}
func (p Repository) UserDB(ctx context.Context, user domain.User) (int, error) {

	var newUser domain.User

	query := "insert into users (password, email) values ($1, $2) returning user_id"
	resp := p.db.QueryRow(query, user.Password, user.Email)
	err := resp.Scan(&newUser.User_id)
	if err != nil {
		fmt.Println("problem UserDB")
	}

	return newUser.User_id, nil
}
func (p Repository) PasswordDB(ctx context.Context, email string) (string, int, error) {
	var Pas domain.User
	query := "select password, user_id from users where email=$1"
	password := p.db.QueryRow(query, email)
	err := password.Scan(&Pas.Password, &Pas.User_id)
	if err != nil {
		fmt.Println("problem paswordQuery")
	}
	log, found := logger.FromContext(ctx)
	if !found {
		fmt.Println("problem PasswordDB log")
	}
	log.Debug("клиент авторизован")

	return Pas.Password, Pas.User_id, nil
}
func (p Repository) ControlTokenDB(ctx context.Context, token string) (int, error) {
	var User domain.Session

	fmt.Println(token)
	//token1 := "2474ea20-e5ca-4fb7-997a-4275117523f8"
	query := "select user_id from session where token=$1"
	err := p.db.QueryRow(query, token).Scan(&User.UserID)
	Userid := User.UserID
	if err != nil {
		fmt.Println("problem ControlToken")

		return Userid, err
	} else {

		return Userid, nil
	}

}
