package tk

import "strings"

const (
	DescOrderPrefix = "-"
	AscOrderPrefix  = "+"
	DefaultOrderBy  = "id desc"
)

func ArrayToQueryOrder(s []string) string {
	v := make([]string, 0, len(s))
	for _, b := range s {
		asc := !strings.HasPrefix(b, DescOrderPrefix)
		o := strings.TrimPrefix(strings.TrimPrefix(b, AscOrderPrefix), DescOrderPrefix)
		if !asc {
			o += " desc"
		}
		v = append(v, o)
	}
	return strings.Join(v, ",")
}
