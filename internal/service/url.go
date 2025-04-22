package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bobbybaiOuO/BShortUrl/internal/model"
	"github.com/bobbybaiOuO/BShortUrl/internal/repo"
)

// ShortCodeGenerator .
type ShortCodeGenerator interface {
	GenerateShortCode() string
}

// Cacher .
type Cacher interface {
	SetURL(ctx context.Context, url repo.Url) error
	GetURL(ctx context.Context, shortCode string) (*repo.Url, error)
}

// URLService .
type URLService struct {
	querier            repo.Querier
	shortCodeGenerator ShortCodeGenerator
	defaultDuration    time.Duration
	cache              Cacher
	baseURL            string
}

// NewURLService ..
func NewURLService(db *sql.DB, shortCodeGenerator ShortCodeGenerator,
duration time.Duration, cache Cacher, baseURL string) *URLService {
	return &URLService{
		querier: repo.New(db),
		shortCodeGenerator: shortCodeGenerator,
		defaultDuration: duration,
		cache: cache,
		baseURL: baseURL,
	}
}

// CreateURL .
func (s *URLService) CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error) {
	var shortCode string
	var isCustom bool
	var expiredAt time.Time

	if req.CustomCode != "" {
		isAvailable, err := s.querier.IsShortCodeAvailable(ctx, req.CustomCode)
		if err != nil {
			return nil, err
		}
		if !isAvailable {
			return nil, fmt.Errorf("CustomCode already exists")
		}
		shortCode = req.CustomCode
		isCustom = true
	} else {
		code, err := s.getShortCode(ctx, 0)
		if err != nil {
			return nil, err
		}
		shortCode = code
	}

	if req.Duration == nil {
		expiredAt = time.Now().Add(s.defaultDuration)
	} else {
		expiredAt = time.Now().Add(time.Hour * time.Duration(*req.Duration))
	}

	// 插入数据库
	url, err := s.querier.CreateURL(ctx, repo.CreateURLParams{
		OriginalUrl: req.OriginalURL,
		ShortCode:   shortCode,
		IsCustom:    isCustom,
		ExpiredAt:   expiredAt,
	})
	if err != nil {
		return nil, err
	}

	// 存入缓存
	if err := s.cache.SetURL(ctx, url); err != nil {
		return nil, err
	}

	return &model.CreateURLResponse{
		ShortURL: s.baseURL + "/" + url.ShortCode,
		ExpiredAt: url.ExpiredAt,
	}, nil
}

// GetURL .
func (s *URLService) GetURL(ctx context.Context, shortCode string) (string, error) {
	// 先访问缓存
	url, err := s.cache.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if url != nil {
		return url.OriginalUrl, nil
	}

	// 访问数据库
	url2, err := s.querier.GetUrlByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	// 存入缓存
	if err := s.cache.SetURL(ctx, url2); err != nil {
		return "", err
	}

	return url2.OriginalUrl, nil
}


func (s *URLService) getShortCode(ctx context.Context, n int) (string, error) {
	if n > 5 {
		return "", errors.New("Too many retries")
	}
	shortCode := s.shortCodeGenerator.GenerateShortCode()

	isAvailable, err := s.querier.IsShortCodeAvailable(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if isAvailable {
		return shortCode, nil
	}

	return s.getShortCode(ctx, n+1)
}

// DeleteURL .
func (s *URLService) DeleteURL(ctx context.Context) error {
	return s.querier.DeleteURLExpired(ctx)
}
