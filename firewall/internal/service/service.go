package service

import "firewall/internal/storage"

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) Service {
	return Service{storage}
}
