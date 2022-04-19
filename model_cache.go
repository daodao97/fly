package ggm

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

func asSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *model[T]) cacheKeyPrefix() string {
	key := m.modelInfo.Name
	t := reflectNew[T]()
	return key + "_" + asSha256(t)
}

func (m *model[T]) FindBy(id int) (T, error) {
	t := reflectNew[T]().(T)
	pk := m.pk()
	if pk == "" {
		return t, ErrPrimaryKeyNotDefined
	}
	if cache == nil {
		return m.SelectOne(WhereEq(pk, id))
	}

	key := m.cacheKeyPrefix()
	c, err := cache.Get(key)
	if err != nil {
		return t, errors.Wrap(err, "get cache error")
	}
	if c != "" {
		err := json.Unmarshal([]byte(c), t)
		if err != nil {
			return t, err
		}
		return t, nil
	}
	row, err := m.SelectOne(WhereEq(pk, id))
	if err != nil {
		return t, err
	}
	bytes, err := json.Marshal(row)
	if err != nil {
		return t, err
	}
	err = cache.Set(key, string(bytes))
	if err != nil {
		return t, err
	}
	return row, nil
}

func (m *model[T]) UpdateBy(id int, row T) (int64, error) {
	pk := m.pk()
	if pk == "" {
		return 0, ErrPrimaryKeyNotDefined
	}
	effect, err := m.Update(row, WhereEq(pk, id))
	if err != nil {
		return 0, err
	}
	if cache == nil {
		return effect, nil
	}
	bytes, err := json.Marshal(row)
	if err != nil {
		return 0, err
	}
	key := m.cacheKeyPrefix()
	err = cache.Set(key, string(bytes))
	if err != nil {
		return 0, err
	}
	return effect, nil
}
