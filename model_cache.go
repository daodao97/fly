package fly

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func (m *model) cacheKeyPrefix(id int64) string {
	return fmt.Sprintf("%s-%s-%d", m.connection, m.table, id)
}

func (m *model) FindBy(id int64) *Row {
	if cache == nil {
		return &Row{Err: errors.New("cache instance is nil")}
	}

	pk := m.PrimaryKey()
	if pk == "" {
		return &Row{Err: errors.New("primary is not defined")}
	}

	key := m.cacheKeyPrefix(id)

	c, err := cache.Get(key)
	if err != nil {
		return &Row{Err: err}
	}
	if c != "" {
		var result map[string]interface{}
		err = json.Unmarshal([]byte(c), &result)
		if err != nil {
			return &Row{Err: err}
		}
		_ = logger.Log(LevelDebug, "FindBy id:", id, "form cache", c)
		return &Row{Data: result}
	}

	row := m.SelectOne(WhereEq(pk, id))
	if row.Err == nil && row.Data != nil {
		c, err := json.Marshal(row.Data)
		if err != nil {
			return &Row{Err: err}
		}
		err = cache.Set(key, string(c))
		if err != nil {
			return &Row{Err: err}
		}
		_ = logger.Log(LevelDebug, "FindBy id:", id, "set cache", string(c))
	}

	return row
}

func (m *model) UpdateBy(id int64, record interface{}) (bool, error) {
	if cache == nil {
		return false, errors.New("cache instance is nil")
	}
	_, err := m.Update(record, WhereEq("id", id))
	if err != nil {
		return false, err
	}
	key := m.cacheKeyPrefix(id)
	err = cache.Del(key)
	if err != nil {
		return false, err
	}
	_ = logger.Log(LevelDebug, "del key after UpdateBy id:", id)

	return true, nil
}
