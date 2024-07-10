package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore() (*MySQLStore, error) {
	dsn := "root:pass$@tcp(127.0.0.1:3306)/bank_db?parseTime=true"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MySQLStore{
		db: db,
	}, nil
}

func (s *MySQLStore) Init() error {
	return s.CreateAccountTable()
}

func (s *MySQLStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS user (
				id INT PRIMARY KEY AUTO_INCREMENT,
				first_name VARCHAR(50),
				last_name VARCHAR(50),
				number INT UNIQUE,
				encrypted_password VARCHAR(50),
				balance DECIMAL(10,2),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) CreateAccount(acc *Account) error {
	query := `
	insert into user
	(first_name, last_name, number, encrypted_password, balance, created_at)
	values
	(?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Query(query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) UpdateAccount(*Account) error {
	return nil
}

func (s *MySQLStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM user WHERE id = ?", id)

	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE number = ?", number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", number)
}

func (s *MySQLStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *MySQLStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
