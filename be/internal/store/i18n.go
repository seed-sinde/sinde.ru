package store

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"sinde.ru/internal/media"
)

type LocaleCode string

const (
	LocaleEN LocaleCode = "en"
	LocaleRU LocaleCode = "ru"
	LocaleZH LocaleCode = "zh"
	LocaleJA LocaleCode = "ja"
)

type I18nRawSection map[string]map[LocaleCode]string
type I18nCommonSection map[string]string
type I18nResolvedSection map[string]string

type sectionCacheItem struct {
	data I18nRawSection
	exp  time.Time
}

type commonCacheItem struct {
	data I18nCommonSection
	exp  time.Time
}

var ErrI18nNamespaceNotFound = errors.New("i18n namespace not found")

type I18nStore struct {
	bucket   string
	ttl      time.Duration
	mu       sync.RWMutex
	sections map[string]sectionCacheItem
	common   commonCacheItem
}

func NewI18nStore() *I18nStore {
	return &I18nStore{
		bucket:   "i18n",
		ttl:      5 * time.Minute,
		sections: make(map[string]sectionCacheItem),
	}
}
func (s *I18nStore) getSectionCached(ns string) (I18nRawSection, bool) {
	s.mu.RLock()
	it, ok := s.sections[ns]
	s.mu.RUnlock()
	return it.data, ok && time.Now().Before(it.exp)
}
func (s *I18nStore) getCommonCached() (I18nCommonSection, bool) {
	s.mu.RLock()
	it := s.common
	s.mu.RUnlock()
	return it.data, it.data != nil && time.Now().Before(it.exp)
}
func (s *I18nStore) setSectionCache(ns string, data I18nRawSection) {
	s.mu.Lock()
	s.sections[ns] = sectionCacheItem{
		data: data,
		exp:  time.Now().Add(s.ttl),
	}
	s.mu.Unlock()
}
func (s *I18nStore) setCommonCache(data I18nCommonSection) {
	s.mu.Lock()
	s.common = commonCacheItem{
		data: data,
		exp:  time.Now().Add(s.ttl),
	}
	s.mu.Unlock()
}
func (s *I18nStore) ResolveSection(ctx context.Context, locale LocaleCode, namespace string) (I18nResolvedSection, error) {
	ns := strings.TrimSpace(namespace)
	if ns == "" {
		return nil, ErrI18nNamespaceNotFound
	}
	if ns == "common" {
		common, err := s.getCommon(ctx)
		if err != nil {
			return nil, err
		}
		out := make(I18nResolvedSection, len(common))
		for key, val := range common {
			trimmed := strings.TrimSpace(val)
			out[key] = key
			if trimmed != "" {
				out[key] = trimmed
			}
		}
		return out, nil
	}
	raw, err := s.getSection(ctx, ns)
	if err != nil {
		return nil, err
	}
	common, err := s.getCommon(ctx)
	if err != nil {
		return nil, err
	}
	return s.resolve(raw, common, locale), nil
}

func (s *I18nStore) getSection(ctx context.Context, namespace string) (I18nRawSection, error) {
	ns := strings.TrimSpace(namespace)
	if ns == "" || ns == "common" {
		return nil, ErrI18nNamespaceNotFound
	}
	if data, ok := s.getSectionCached(ns); ok {
		return data, nil
	}
	body, err := media.ReadObjectBytesFromBucket(ctx, s.bucket, ns+".json")
	if err != nil {
		return nil, ErrI18nNamespaceNotFound
	}
	var data I18nRawSection
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	s.setSectionCache(ns, data)
	return data, nil
}

func (s *I18nStore) getCommon(ctx context.Context) (I18nCommonSection, error) {
	if data, ok := s.getCommonCached(); ok {
		return data, nil
	}
	body, err := media.ReadObjectBytesFromBucket(ctx, s.bucket, "common.json")
	if err != nil {
		return nil, ErrI18nNamespaceNotFound
	}
	var data I18nCommonSection
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	s.setCommonCache(data)
	return data, nil
}

func (s *I18nStore) resolve(raw I18nRawSection, common I18nCommonSection, locale LocaleCode) I18nResolvedSection {
	out := make(I18nResolvedSection, len(common)+len(raw))
	for key, val := range common {
		out[key] = key
		if trimmed := strings.TrimSpace(val); trimmed != "" {
			out[key] = trimmed
		}
	}
	for key, item := range raw {
		if val := strings.TrimSpace(item[locale]); val != "" {
			out[key] = val
			continue
		}
		if val := strings.TrimSpace(item[LocaleEN]); val != "" {
			out[key] = val
			continue
		}
		if val := strings.TrimSpace(common[key]); val != "" {
			out[key] = val
			continue
		}
		out[key] = key
	}
	return out
}
