package postgresdb

import (
	"auth-service/internal/config"
	"auth-service/internal/repositories/user"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DataBase struct {
	DB *sql.DB
}

const (
	RoleUser   = "user"
	ROLE_ADMIN = "admin"
)

const (
	MODE_DROP_ON_START = "drop-on-start"

	DROP_TABLE = "DROP TABLE IF EXISTS users;"

	CREATE_TABLE = `
	CREATE TABLE IF NOT EXISTS users
	(
		user_id serial PRIMARY KEY,
		username varchar(80) NOT NULL UNIQUE,
		password text NOT NULL,
		role text
	);`

	InsertUser = "INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING user_id;"

	SelectUserByUsername = "SELECT * FROM users WHERE username = $1;"

	CHECK_IF_USER_EXISTS = "SELECT 1 FROM users WHERE username = $1;"
)

func New(dbConfig config.DB) (*DataBase, error) {
	op := "db.psql.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if dbConfig.Mode == MODE_DROP_ON_START {
		statement, err := db.Prepare(DROP_TABLE)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		_, err = statement.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	statement, err := db.Prepare(CREATE_TABLE)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = statement.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &DataBase{DB: db}, nil
}

func (db *DataBase) SaveUser(username, password string) (int, error) {
	op := "db.psql.SaveUser"

	statement, err := db.DB.Prepare(InsertUser)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	var result int
	err = statement.QueryRow(username, password, RoleUser).Scan(&result)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (db *DataBase) GetUser(username string) (*user.User, error) {
	op := "db.psql.GetUser"

	statement, err := db.DB.Prepare(SelectUserByUsername)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var p user.User
	err = statement.QueryRow(username).Scan(&p.UserId, &p.Username, &p.Password, &p.Role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &p, nil
}

/*func (db *DataBase) CheckIfUserExists(alias string) (bool, error) {
	op := "db.psql.CheckIfUserExists"

	statement, err := db.DB.Prepare(CHECK_IF_USER_EXISTS)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var result bool
	err = statement.QueryRow(alias).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}*/
