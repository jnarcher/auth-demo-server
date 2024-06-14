package database

import (
	"auth-demo/internal/model"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("DB_HOST")
	user     = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
	port     = os.Getenv("DB_PORT")
)

type PostgresStore struct {
	mux sync.RWMutex
	db  *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	query := `CREATE TABLE IF NOT EXISTS account  (
        id SERIAL PRIMARY KEY,
        username VARCHAR(255) UNIQUE NOT NULL,
        pwd_hash VARCHAR(255) NOT NULL,
        first_name VARCHAR(255),
        last_name VARCHAR(255),
        email VARCHAR(255),
        phone VARCHAR(255),
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ NULL
    );

    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ language 'plpgsql';

    CREATE OR REPLACE TRIGGER update_account_updated_at
    BEFORE UPDATE ON account
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *model.Account) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	query := `
    INSERT INTO account
        (username, pwd_hash, first_name, last_name, email, phone)
    VALUES
        ($1, $2, $3, $4, $5, $6)
    ;`

	if _, err := s.db.Query(query,
		acc.User,
		acc.PwdHash,
		acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Phone,
	); err != nil {
		return err
	}

	log.Printf("New account created: `%s`", acc.User)

	query = `SELECT id FROM account WHERE username = $1 LIMIT 1;`
	row := s.db.QueryRow(query, acc.User)
	if row == nil {
		return fmt.Errorf("Failed to insert new account")
	}

	id := &struct{ Id int64 }{}
	if err := row.Scan(&id.Id); err != nil {
		return err
	}
	acc.Id = id.Id
	return nil
}
func (s *PostgresStore) UpdateAccount(acc *model.Account) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	query := `UPDTATE account
    SET
        pwd_hash    = $1,
        first_name  = $2,
        last_name   = $3,
        email       = $4, 
        phone       = $5
    WHERE
        id = $6
    ;`

	_, err := s.db.Exec(query,
		acc.PwdHash,
		acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Phone,
	)
	return err
}
func (s *PostgresStore) DeleteAccount(id int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	query := `UPDATE account
    SET 
        deleted_at = NOW() 
    WHERE
        id = $1 AND deleted_at IS NULL
    ;`
	res, err := s.db.Exec(query, id)
	rows, err := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account doesn't exist")
	}
	return err
}

func (s *PostgresStore) GetAccounts() ([]*model.Account, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	rows, err := s.db.Query("SELECT * FROM account WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	accounts := []*model.Account{}
	for rows.Next() {
		acc, err := scanRowsIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id int) (*model.Account, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	query := `SELECT * FROM account WHERE id = $1 AND deleted_at IS NULL;`
	row := s.db.QueryRow(query, id)
	if row == nil {
		return nil, fmt.Errorf("Account %d not found", id)
	}

	acc, err := scanRowIntoAccount(row)
	if err != nil {
		return nil, fmt.Errorf("Account %d not found", id)
	}

	return acc, nil
}

func (s *PostgresStore) GetAccountByUser(user string) (*model.Account, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	query := `SELECT * FROM account WHERE username = $1 AND deleted_at IS NULL;`
	row := s.db.QueryRow(query, user)
	if row == nil {
		return nil, fmt.Errorf("Account `%s` not found", user)
	}

	acc, err := scanRowIntoAccount(row)
	if err != nil {
		return nil, fmt.Errorf("Account `%s` not found", user)
	}

	return acc, nil
}

func scanRowsIntoAccount(rows *sql.Rows) (*model.Account, error) {
	acc := &model.Account{}
	err := rows.Scan(
		&acc.Id,
		&acc.User,
		&acc.PwdHash,
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Phone,
		&acc.CreatedAt,
		&acc.UpdatedAt,
		&acc.DeletedAt,
	)
	return acc, err
}

func scanRowIntoAccount(row *sql.Row) (*model.Account, error) {
	acc := &model.Account{}
	err := row.Scan(
		&acc.Id,
		&acc.User,
		&acc.PwdHash,
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Phone,
		&acc.CreatedAt,
		&acc.UpdatedAt,
		&acc.DeletedAt,
	)
	return acc, err
}
