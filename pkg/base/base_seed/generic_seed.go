package base_seed

import (
	"bytes"
	"cookie_supply_management/core/config"
	"cookie_supply_management/core/connect"
	"cookie_supply_management/pkg/base/base_model"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

const (
	BASE_ERROR_MESSAGE string = "FAILED_SEEDING"
)

type Summary struct {
	Created int64
	Exist   int64
	Updated int64
	Errors  int64
}

type GenericSeed struct {
	Model        interface{}
	FixtureBased bool
	Summary      Summary
	Query        *gorm.DB
}

func (s *GenericSeed) Seed() (Summary, error) {
	defer s.Summarize()
	return Summary{}, fmt.Errorf("a custom implementation is needed")
}

func NewGenericSeed(model interface{}) GenericSeed {
	tableName := base_model.GetTableName(model, connect.DB)
	query := connect.DB.Table(tableName)
	return GenericSeed{Model: model, Query: query}
}

func (s *GenericSeed) Log(message string, values ...interface{}) error {
	tableName := base_model.GetTableName(s.Model, connect.DB)
	if len(values) != 0 {
		message = fmt.Sprintf(message, values)
	}
	log.Printf("[%s] - %s", tableName, message)
	return nil
}

func (s *GenericSeed) LogFail(message string, values ...interface{}) {
	if len(values) != 0 {
		message = fmt.Sprintf(message, values)
	}
	s.Log(fmt.Sprintf("%s - %s", BASE_ERROR_MESSAGE, message))
}

func (s *GenericSeed) GetFileName() string {
	tableName := base_model.GetTableName(s.Model, connect.DB)
	return fmt.Sprintf("%s.json", tableName)
}

func (s *GenericSeed) LoadFixture() ([]byte, error) {
	baseDir := config.Get().Dir.Seeder
	fileName := s.GetFileName()
	path := filepath.Join(baseDir, fileName)
	data, err := os.ReadFile(path)
	return data, err
}

func (s *GenericSeed) Error(err error) (Summary, error) {
	s.LogFail(err.Error())
	s.Summary.Errors++
	return s.Summary, err
}

type Kwargs struct {
	TableName string
	Name      string
	Value     interface{}
}

type KwargsToField struct {
	TableName string
	Name1     string
	Name2     string
	Value1    interface{}
	Value2    interface{}
}

func (s *GenericSeed) Exists(kwargs interface{}) bool {
	var exists bool
	var data Kwargs

	switch v := kwargs.(type) {
	case map[string]interface{}:
		for name, value := range v {
			data = Kwargs{
				Name:  name,
				Value: value,
			}
		}
	case Kwargs:
		data = v
	}

	var buf *bytes.Buffer = &bytes.Buffer{}
	var existsTemplate string = `
		SELECT EXISTS(
			SELECT 1
			FROM "{{.TableName}}"
			WHERE "{{.Name}}" = '{{.Value}}'
		)
	`

	data.TableName = base_model.GetTableName(s.Model, connect.DB)

	var tmpl *template.Template = template.New("exists")
	tmpl.Parse(existsTemplate)

	if err := tmpl.Execute(buf, data); err != nil {
		s.LogFail(err.Error())
	}

	r := connect.DB.Raw(buf.String()).Scan(&exists)

	if err := r.Error; err != nil {
		s.LogFail(err.Error())
		return false
	}

	return exists
}

func (s *GenericSeed) ExistsToField(data KwargsToField) bool {
	var exists bool

	//switch v := kwargs.(type) {
	//case map[string]interface{}:
	//	for name, value := range v {
	//		data = KwargsToField{
	//			Name1:  name,
	//			Value1: value,
	//		}
	//	}
	//case KwargsToField:
	//	data = v
	//}

	var buf *bytes.Buffer = &bytes.Buffer{}
	var existsTemplate string = `
		SELECT EXISTS(
			SELECT 1
			FROM "{{.TableName}}"
			WHERE "{{.Name1}}" = '{{.Value1}}' and "{{.Name2}}" = '{{.Value2}}'
		)
	`

	data.TableName = base_model.GetTableName(s.Model, connect.DB)

	var tmpl *template.Template = template.New("exists")
	tmpl.Parse(existsTemplate)

	if err := tmpl.Execute(buf, data); err != nil {
		s.LogFail(err.Error())
	}

	r := connect.DB.Raw(buf.String()).Scan(&exists)

	if err := r.Error; err != nil {
		s.LogFail(err.Error())
		return false
	}

	return exists
}

type LanguageValues struct {
	En string
	Ru string
	Ky string
}

func (s *GenericSeed) MakeNameJSON(values LanguageValues) []byte {
	data, err := json.Marshal(values)
	if err != nil {
		s.LogFail("failed to marshal data to json: %s", err.Error())
		s.Summary.Errors++
	}
	return data
}

func (s *GenericSeed) Summarize() {
	colorGreen := "\033[32m"
	colorRed := "\033[31m"
	colorYellow := "\033[33m"
	colorReset := "\033[0m"

	created := fmt.Sprintf("%s %d created %s", colorGreen, s.Summary.Created, colorReset)
	exist := fmt.Sprintf("%s %d exist %s", colorYellow, s.Summary.Exist, colorReset)
	errors := fmt.Sprintf("%s %d errors %s", colorRed, s.Summary.Errors, colorReset)
	update := fmt.Sprintf("%s %d updated %s", colorGreen, s.Summary.Updated, colorReset)
	tableName := base_model.GetTableName(s.Model, connect.DB)
	log.Printf("[%s] summary: %s - %s - %s - %s", tableName, created, update, exist, errors)
}
