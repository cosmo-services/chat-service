package pkg

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type Page struct {
	Items   interface{}
	HasNext bool
	HasPrev bool
}

func Paginate(db *gorm.DB, dest interface{}, cursor string, column string, direction string, limit int) (*Page, error) {
	baseQuery := db.Session(&gorm.Session{})
	dataQuery := baseQuery.Session(&gorm.Session{})

	if cursor != "" {
		switch direction {
		case "next":
			dataQuery = dataQuery.Where(fmt.Sprintf("%s > ?", column), cursor).
				Order(fmt.Sprintf("%s ASC", column))
		case "prev":
			dataQuery = dataQuery.Where(fmt.Sprintf("%s < ?", column), cursor).
				Order(fmt.Sprintf("%s DESC", column))
		default:
			dataQuery = dataQuery.Order(fmt.Sprintf("%s DESC", column))
		}
	} else {
		dataQuery = dataQuery.Order(fmt.Sprintf("%s DESC", column))
	}

	if err := dataQuery.Limit(limit).Find(dest).Error; err != nil {
		return nil, err
	}

	items := reflect.ValueOf(dest).Elem()
	if items.Len() == 0 {
		return &Page{
			Items:   dest,
			HasNext: false,
			HasPrev: false,
		}, nil
	}

	var minVal, maxVal string
	err := baseQuery.Session(&gorm.Session{}).
		Select(fmt.Sprintf("MIN(%s), MAX(%s)", column, column)).
		Row().Scan(&minVal, &maxVal)
	if err != nil {
		return nil, err
	}

	var hasPrev bool
	if minVal != "" {
		prevCheck := baseQuery.Session(&gorm.Session{}).
			Where(fmt.Sprintf("%s < ?", column), minVal)
		err = prevCheck.Select("EXISTS(?)", prevCheck).Find(&hasPrev).Error
		if err != nil {
			return nil, err
		}
	}

	var hasNext bool
	if maxVal != "" {
		nextCheck := baseQuery.Session(&gorm.Session{}).
			Where(fmt.Sprintf("%s > ?", column), maxVal)
		err = nextCheck.Select("EXISTS(?)", nextCheck).Find(&hasNext).Error
		if err != nil {
			return nil, err
		}
	}

	return &Page{
		Items:   dest,
		HasNext: hasNext,
		HasPrev: hasPrev,
	}, nil
}
