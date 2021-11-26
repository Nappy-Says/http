package banners

import (
	"context"
	"errors"
	"log"
	"sync"
)

var tempBannerID int64
// var panicNotEmplemented = "not emplemented"
var errorBannerNotFound = errors.New("banner not found")


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

func NewService() *Service  {
	return &Service{items: make([]*Banner, 0)}
}



func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

	return nil, errorBannerNotFound
}


func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	if item == nil {
		tempBannerID++

		createdBanner := &Banner{
			Link:    item.Link,
			Title:   item.Title,
			Button:  item.Button,
			Content: item.Content,
			ID:      tempBannerID,
		}

		s.items = append(s.items, createdBanner)

		return createdBanner, nil
	}

	banner, err := s.ByID(ctx, item.ID)

	if err != nil {
		log.Print(err)
		return nil, errorBannerNotFound
	}

	banner.Link = item.Link
	banner.Title = item.Title
	banner.Button = item.Button
	banner.Content = item.Content

	return banner, nil

}


func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	banner, err := s.ByID(ctx, id)
	if err != nil {
		log.Print(err)
		return nil, errorBannerNotFound
	}

	for i, j := range s.items {
		if j.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			break
		}
	}

	return banner, nil
}
