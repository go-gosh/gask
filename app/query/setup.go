package query

import "github.com/go-gosh/gask/app/global"

func Setup() error {
	db, err := global.GetDatabase()
	if err != nil {
		return err
	}
	SetDefault(db)
	return nil
}
