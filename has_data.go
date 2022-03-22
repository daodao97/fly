package ggm

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type hasOpts struct {
	Conn        string
	DB          string
	Table       string
	LocalKey    string
	ForeignKey  string
	OtherKeys   []string
	RefType     reflect.Type
	StructField string
}

// [pool.]db.table:[local_key->]foreign_key,other_key
func explodeHasStr(str string) (opt *hasOpts, err error) {
	var re = regexp.MustCompile(`([a-zA-Z_0-9]+\.)?([a-zA-Z_0-9]+\.)?([a-zA-Z_0-9]+):([a-zA-Z_0-9]+->)?([a-zA-Z_0-9 ]+)?([a-zA-Z_,0-9 ]+)`)
	if !re.MatchString(str) {
		return nil, fmt.Errorf("has string syntux is error, mast like [pool.]db.table:[local_key->]foreign_key,other_key")
	}
	matched := re.FindStringSubmatch(str)
	for i, v := range matched {
		if i == 0 {
			continue
		}
		matched[i] = strings.ReplaceAll(strings.ReplaceAll(v, "->", ""), ".", "")
	}

	opt = &hasOpts{
		Conn:       matched[1],
		DB:         matched[2],
		Table:      matched[3],
		LocalKey:   matched[4],
		ForeignKey: matched[5],
		OtherKeys:  filterEmptyStr(strings.Split(matched[6], ",")),
	}
	if opt.Conn == "" {
		opt.Conn = "default"
	}
	if opt.LocalKey == "" {
		opt.LocalKey = "id"
	}
	if opt.ForeignKey == "" {
		opt.ForeignKey = "id"
	}

	return opt, nil
}
