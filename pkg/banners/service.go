package banners

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"sync"
)

var tempBannerID int64
// var panicNotEmplemented = "not emplemented"
var errorBannerNotFound = errors.New("banner not found")
var errorReadFile = errors.New("an error occurred while reading the file")

type Service struct {
	mu 		sync.RWMutex
	items 	[]*Banner
}

type Banner struct {
	ID 		int64
	Link	string
	Image	string
	Title 	string
	Button 	string
	Content string
}

func NewService() *Service  {
	return &Service{items: make([]*Banner, 0)}
}



func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items != nil {
		return s.items, nil
	}

	return nil, errorBannerNotFound
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


// func (s *Service) Save(ctx context.Context, item *Banner, file *multipart.File) (*Banner, error) {
func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	
	if item.ID == 0 {
		tempBannerID++
		item.ID = tempBannerID

		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image

			img, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, errorReadFile
			}

			err = ioutil.WriteFile("./web/banners/" + item.Image, img, 0666)
			if err != nil {
				return nil, err
			}
		}

		s.items = append(s.items, item)
		return item, nil
	}

	for k, v := range s.items {
		if v.ID == item.ID {
			if item.Image != "" {
				item.Image = fmt.Sprint(item.ID) + "." + item.Image
				var data, err = ioutil.ReadAll(file)
				if err != nil {
					return nil, err
				}

				err = ioutil.WriteFile("./web/banners/"+item.Image, data, 0666)

				if err != nil {
					return nil, err
				}
			} else {
				item.Image = s.items[k].Image
			}

			s.items[k] = item
			return item, nil
		}
	}

	return nil, errorBannerNotFound
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
