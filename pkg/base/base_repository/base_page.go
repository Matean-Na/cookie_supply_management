package base_repository

import "gorm.io/gorm"

type Pager interface {
	GetPage() int
	GetPageSize() int
	GetOffset() int

	Paginate() func(db *gorm.DB) *gorm.DB
}

type Page struct {
	Page, PageSize, Offset int
}

func NewPager(page, pageSize, offset int) Pager {
	return &Page{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

func (p *Page) GetPage() int {
	return p.Page
}

func (p *Page) GetPageSize() int {
	return p.PageSize
}

func (p *Page) GetOffset() int {
	return p.Offset
}

func (p *Page) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetPageSize())
	}
}
