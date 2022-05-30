package ggm

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func (m *model) cacheKeyPrefix(id int64) string {
	return fmt.Sprintf("%s-%s-%d", m.connName, m.modelInfo.Name, id)
}

func (m *model) FindBy(id int64, dest interface{}) (err error) {
	var kv []interface{}
	defer dbLog("FindBy", time.Now(), &err, &kv)
	pk := m.pk()
	if pk == "" {
		return ErrPrimaryKeyNotDefined
	}
	if cache == nil {
		return m.SelectOne(WhereEq(pk, id))
	}

	key := m.cacheKeyPrefix(id)
	c, err := cache.Get(key)
	if err != nil {
		return errors.Wrap(err, "get cache error")
	}
	if c != "" {
		err := json.Unmarshal([]byte(c), dest)
		if err != nil {
			return err
		}
		kv = append(kv, "load from cache", key)
		return nil
	}

	kv = append(kv, "load from db", key)
	err = m.SelectOne(dest, WhereEq(pk, id))
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(dest)
	if err != nil {
		return err
	}
	err = cache.Set(key, string(bytes))
	if err != nil {
		return err
	}
	return nil
}

func (m *model) UpdateBy(id int64, row interface{}) (int64, error) {
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
	key := m.cacheKeyPrefix(id)
	err = cache.Set(key, string(bytes))
	if err != nil {
		return 0, err
	}
	return effect, nil
}
