package postgres

import (
	pb "auth-service/genproto/user"
	"auth-service/models"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Add(ctx context.Context, u *models.RegisterRequest) (*models.RegisterResponse, error) {
	query := `
	insert into
		users (email, hashed_password, first_name, last_name, phone_number, role)
	values
		($1, $2, $3, $4, $5, $6)
	returning id, created_at
	`

	var resp models.RegisterResponse
	err := r.DB.QueryRowContext(ctx, query, u.Email, u.Password, u.FirstName, u.LastName,
		u.PhoneNumber, u.Role).Scan(&resp.ID, &resp.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "query execution failed")
	}

	return &resp, nil
}

func (r *UserRepo) Read(ctx context.Context, u *pb.ID) (*pb.Profile, error) {
	query := `
	select
		email, first_name, last_name, phone_number, created_at, updated_at
	from
		users
	where
		id = $1
	`

	var resp pb.Profile
	err := r.DB.QueryRowContext(ctx, query, u.Id).Scan(&resp.Email, &resp.FirstName, &resp.LastName,
		&resp.PhoneNumber, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.Wrap(err, "query execution failed")
	}

	return &resp, nil
}

func (r *UserRepo) Update(ctx context.Context, u *pb.NewData) (*pb.UpdateResp, error) {
	query := `
	update
		users
	set
		email = $1, first_name = $2, last_name = $3, phone_number = $4, updated_at = NOW()
	where
		id = $5
	returning
		id, updated_at
	`

	var resp pb.UpdateResp
	err := r.DB.QueryRowContext(ctx, query, u.Email, u.FirstName, u.LastName,
		u.PhoneNumber, u.Id).Scan(&resp.Id, &resp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.Wrap(err, "query execution failed")
	}

	return &resp, nil
}

func (r *UserRepo) GetDetails(ctx context.Context, email string) (*models.UserDetails, error) {
	query := `
	select
		id, hashed_password, role
	from
		users
	where
		email = $1
	`

	var resp models.UserDetails
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&resp.Id, &resp.Password, &resp.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.Wrap(err, "query execution failed")
	}

	return &resp, nil
}

func (r *UserRepo) GetRole(ctx context.Context, id string) (string, error) {
	query := `
	select
		role
	from
		users
	where
		id = $1
	`

	var role string
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", errors.Wrap(err, "query execution failed")
	}

	return role, nil
}
