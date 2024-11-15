package book

import (
	"context"
	"github.com/TariqueNasrullah/otel-practice/proto"
)

type Service struct {
	proto.UnimplementedBookServiceServer
}

func (s *Service) Create(ctx context.Context, book *proto.Book) (*proto.Book, error) {
	return book, nil
}

func (s *Service) Update(ctx context.Context, book *proto.Book) (*proto.Book, error) {
	return book, nil
}

func NewService() proto.BookServiceServer {
	return &Service{}
}
