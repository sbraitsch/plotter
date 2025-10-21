package service

import (
	"github.com/sbraitsch/plotter/internal/storage"
)

type AuthService interface {
}

type authServiceImpl struct {
	storage *storage.StorageClient
}

func NewAuthService(storage *storage.StorageClient) AuthService {
	return &authServiceImpl{storage: storage}
}
