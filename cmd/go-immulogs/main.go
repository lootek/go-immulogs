package main

import (
	"context"

	"github.com/lootek/go-immulogs"
	"github.com/lootek/go-immulogs/pkg/service"
	"github.com/lootek/go-immulogs/pkg/storage"
)

func main() {
	s := immulogs.NewService(service.NewREST(), storage.NewMemory())

	ctx := context.Background()
	if err := s.Run(ctx); err != nil {
		return
	}
}
