package pagination

import (
	"github.com/getevo/evo/v2"
	"github.com/getevo/evo/v2/lib/outcome"
	"github.com/getevo/evo/v2/lib/ptr"
	"gorm.io/gorm"
)

type Options struct {
	Size    int  `json:"size"`
	Page    int  `json:"page"`
	MaxSize int  `json:"maxSize"`
	Debug   bool `json:"debug"`
}

// Pagination represents a utility type for handling pagination in Go.
//
// Fields:
// - Records: Total number of rows.
// - CurrentPage: Current page loaded.
// - Pages: Total number of pages.
// - Limit: Number of rows per page.
// - First: First page.
// - Last: Last page.
// - PageRange: Range of visible pages.
//
// Methods:
// - SetCurrentPage: Sets the current page based on the provided value. If the value is 0, the current page is set to 1.
// - SetLimit: Sets the limit of rows per page. If the value is 0, the limit is set to the minimum limit of 10. If the limit is less than the minimum limit, it is set to the minimum
type Pagination struct {
	Model       *gorm.DB    `json:"-"`
	Executed    bool        `json:"-"`
	Debug       bool        `json:"-"`
	Success     bool        `json:"success"`
	Error       *string     `json:"error,omitempty"`
	Records     int         `json:"records"`      // Total rows
	CurrentPage int         `json:"current_page"` // Current Page loaded
	Pages       int         `json:"pages"`        // total number of pages
	Size        int         `json:"size"`         // number of rows per page
	MaxSize     int         `json:"max_size"`
	First       int         `json:"first"` // First Page
	Last        int         `json:"last"`  // Last Page
	Data        interface{} `json:"data"`
}

func (p *Pagination) SetMaxSize(limit int) {
	p.MaxSize = limit
}

// setPages sets the total number of pages in the pagination struct based on the number of records and the limit per page.
// If the number of records is 0, it sets the number of pages to 1.
// If there is no remainder when dividing the number of records by the limit, it sets the number of pages to the integer division.
// Otherwise, it sets the number of pages to the integer division plus 1.
// If the number of pages is 0, it sets it to 1.
// After setting the number of pages, it calls the SetLast and SetPageRange methods to update the last page indicator and the range of visible pages respectively.
func (p *Pagination) setPages() {

	if p.Records == 0 {
		p.Pages = 1
		return
	}

	res := p.Records % p.Size
	if res == 0 {
		p.Pages = p.Records / p.Size
	} else {
		p.Pages = (p.Records / p.Size) + 1

	}

	if p.Pages == 0 {
		p.Pages = 1
	}

	p.Last = p.GetOffset() + p.Size
	if p.Last > p.Records {
		p.Last = p.Records
	}

	if p.Size < 10 {
		p.Size = 10
	}
	if p.Size > p.MaxSize {
		if p.MaxSize == 0 {
			p.MaxSize = 50
		}
		p.Size = p.MaxSize
	}
	if p.CurrentPage > 0 {
		if p.CurrentPage > p.Pages {
			p.CurrentPage = p.Pages
		}
	} else {
		p.CurrentPage = 1
	}
}

// GetOffset calculates the offset for paginating the data based on the current page and limit
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.Size
}

// GetPage returns the current page of the pagination struct
func (p *Pagination) GetPage() int {
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}

	return p.CurrentPage
}

func New(model *gorm.DB, request *evo.Request, out interface{}, options ...Options) (*Pagination, error) {
	var err error
	var p = Pagination{
		Model: model,
		Size:  10,
	}
	p.Size = request.Query("size").Int()
	p.CurrentPage = request.Query("page").Int()

	for _, opt := range options {
		if opt.MaxSize > 0 {
			p.MaxSize = opt.MaxSize
		}
		if opt.Size > 0 {
			p.Size = opt.Size
		}
		if opt.Page > 0 {
			p.CurrentPage = opt.Page
		}
		if opt.Debug {
			p.Debug = opt.Debug
		}
	}

	n, err := p.LoadData(out)
	return n, err
}

func (p *Pagination) LoadData(out interface{}) (*Pagination, error) {
	if p.Debug {
		p.Model = p.Model.Debug()
	}
	var total int64
	if err := p.Model.Count(&total).Error; err != nil {
		return p, err
	}
	p.Records = int(total)
	p.setPages()

	p.Model = p.Model.Limit(p.Size)
	p.Model = p.Model.Offset(p.GetOffset())
	if err := p.Model.Find(out).Error; err != nil {
		if err != nil {
			p.Error = ptr.String("unable to load data from the database")
		}
		return p, err
	}
	p.Success = true
	p.Data = out
	return p, nil
}

func (p *Pagination) GetResponse() outcome.Response {
	var response = outcome.Response{
		Data: p,
	}
	if p.Success {
		response.Status(200)

	} else {
		response.Status(500)
	}
	response.ContentType = "application/json"
	return response
}
