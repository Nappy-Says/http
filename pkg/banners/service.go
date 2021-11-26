package banners

import (
	"context"
	"sync"
)


type Service struct {
	mu 		sync.RWMutex
	items 	[]*Banner
}


type Banner struct {
	ID 		int64
	Title 	string
	Content string
	Button 	string
	Link	string
}

var panicNotEmplemented = "not emplemented"

func NewService() *Service  {
	return &Service{items: make([]*Banner, 0)}
}



func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	// s.mu.RLock()
	// defer s.mu.RUnlock()

	return s.items , nil	
}


func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	panic(panicNotEmplemented)	
}


func (s *Service) Save(ctx context.Context, item *Banner) ([]*Banner, error) {
	panic(panicNotEmplemented)	
}


func (s *Service) RemoveByID(ctx context.Context, id int64) ([]*Banner, error) {
	panic(panicNotEmplemented)	
}
