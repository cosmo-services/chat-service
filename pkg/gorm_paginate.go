package pkg

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type Page struct {
	Items   interface{}
	HasNext bool
	HasPrev bool
}

func Paginate[T any](db *gorm.DB, dest *[]*T, cursor string, column string, direction string, limit int) (*Page, error) {
	if limit <= 0 {
		limit = 20
	}

	base := db.Session(&gorm.Session{})
	data := base.Session(&gorm.Session{})

	if cursor != "" {
		switch direction {
		case "next":
			data = data.Where(column+"> ?", cursor).Order(column + " ASC")
		case "prev":
			data = data.Where(column+"< ?", cursor).Order(column + " DESC")
		default:
			data = data.Order(column + " DESC")
		}
	} else {
		data = data.Order(column + " DESC")
	}

	if err := data.Limit(limit).Find(dest).Error; err != nil {
		return nil, err
	}

	if len(*dest) == 0 {
		return &Page{Items: dest, HasNext: false, HasPrev: false}, nil
	}

	var model T

	var globalMin, globalMax sql.NullString
	err := base.Model(&model).
		Select(fmt.Sprintf("MIN(%s), MAX(%s)", column, column)).
		Row().Scan(&globalMin, &globalMax)
	if err != nil {
		return nil, err
	}

	fieldName, err := getFieldNameByColumn(db, &model, column)
	if err != nil {
		return nil, err
	}
	hasPrev := true
	hasNext := true

	for _, item := range *dest {
		val := reflect.ValueOf(item).Elem().FieldByName(fieldName).Interface()
		strVal := toString(val)

		if globalMin.Valid && strVal == globalMin.String {
			hasPrev = false
		}
		if globalMax.Valid && strVal == globalMax.String {
			hasNext = false
		}
	}

	return &Page{
		Items:   dest,
		HasNext: hasNext,
		HasPrev: hasPrev,
	}, nil
}

func getFieldNameByColumn(db *gorm.DB, model interface{}, columnName string) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}

	for _, field := range stmt.Schema.Fields {
		if field.DBName == columnName {
			return field.Name, nil
		}
	}

	return "", fmt.Errorf("column %s not found in model", columnName)
}

func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case time.Time:
		return v.UTC().Format("2006-01-02T15:04:05.999999Z")
	case *time.Time:
		if v != nil {
			return v.UTC().Format("2006-01-02T15:04:05.999999Z")
		}
		return ""
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}
