package postgres

import (
	"auth-service/config"
	pb "auth-service/genproto/user"
	"auth-service/models"
	"auth-service/storage"
	"context"
	"log"
	"testing"
)

func Repo() storage.IUserStorage {
	db, err := ConnectDB(&config.Config{
		DB_HOST:     "postgres",
		DB_PORT:     "5432",
		DB_USER:     "postgres",
		DB_PASSWORD: "root",
		DB_NAME:     "car_wash_auth",
	})
	if err != nil {
		log.Fatalf("error while connecting to postgres: %v", err)
	}

	return db.User()
}

func TestAdd(t *testing.T) {
	r := Repo()

	_, err := r.Add(context.Background(), &models.RegisterRequest{
		Email:     "test",
		Password:  "test",
		FirstName: "test",
		LastName:  "test",
		Role:      "customer",
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	r := Repo()

	_, err := r.Read(context.Background(), &pb.ID{Id: "84084758-519e-4651-a84e-2bee6a95564c"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	r := Repo()

	_, err := r.Update(context.Background(), &pb.NewData{
		Email:       "test_email",
		FirstName:   "test",
		LastName:    "test",
		PhoneNumber: "test_phone",
		Id:          "84084758-519e-4651-a84e-2bee6a95564c",
	})

	if err != nil {
		t.Fatal(err)
	}
}
