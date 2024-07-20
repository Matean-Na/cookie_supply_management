package base_model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

type HasId interface {
	GetId() uint
	SetId(id uint)
	GetReadOnlyFields() []string
}

type Entity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func GetTableName(m interface{}, db *gorm.DB) string {
	st := &gorm.Statement{DB: db}
	if err := st.Parse(m); err != nil {
		return ""
	}
	return st.Schema.Table
}

func (e *Entity) GetId() uint {
	return e.ID
}

func (e *Entity) SetId(id uint) {
	e.ID = id
}

func (e *Entity) GetReadOnlyFields() []string {
	return []string{}
}

func GetType(m interface{}) string {
	t := reflect.TypeOf(m)

	for t.Kind() != reflect.Struct {
		t = t.Elem()
	}

	ft := fmt.Sprintf("%s", t)
	a := strings.Split(ft, ".")
	return fmt.Sprintf("%s", a[1])
}

func (e Entity) MakeJson(En string, Ru string, Ky string) (result []byte, err error) {
	var data = map[string]interface{}{
		"en": En,
		"ru": Ru,
		"ky": Ky,
	}

	j, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	return j, nil
}
