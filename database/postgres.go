package database

import (
	"context"
	"log"

	"github.com/AndresLince/go-rest-websockets/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.Exec(ctx, "INSERT INTO users (id, email, password) VALUES($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.Query(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	var user = models.User{}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return &user, nil
}
func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.Query(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	var user = models.User{}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return &user, nil
}
