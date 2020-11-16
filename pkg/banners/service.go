package banners

import (
	"context"
	"errors"
	"sort"
	"sync"
)

type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
}

type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

//NewService construct
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) <= 0 {
		return nil, errors.New("banner not found")
	}
	return s.items, nil
}

func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("banner not found")
}

func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		sort.Slice(s.items[:], func(i, j int) bool {
			return s.items[i].ID > s.items[j].ID
		})
		item.ID = s.items[0].ID + 1
		s.items = append(s.items, item)
		return item, nil
	} else {
		for _, banner := range s.items {
			if banner.ID == item.ID {
				break
			}
			banner = &Banner{
				Title:   item.Title,
				Content: item.Content,
				Button:  item.Button,
				Link:    item.Link,
			}
			return banner, nil
		}
	}
	return nil, errors.New("banner not found")
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return banner, nil
		}
	}
	return nil, errors.New("banner not found")
}
