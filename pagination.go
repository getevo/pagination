package pagination

import (
	"github.com/getevo/evo/v2"
	"github.com/getevo/evo/v2/lib/outcome"
	"github.com/getevo/evo/v2/lib/ptr"
	"gorm.io/gorm"
)

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

// SetCurrentPage sets the value of CurrentPage in the Pagination struct.
// If the input page is not equal to zero, p.CurrentPage will be set to the input page.
// Otherwise, p.CurrentPage will be set to 1.
func (p *Pagination) setCurrentPage(page int) {
	if page > 0 {
		p.CurrentPage = page
	} else {
		p.CurrentPage = 1
	}
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

	p.setLast()

}

// setLast sets the value of the Last page in the pagination struct.
// It calculates the value by adding the current offset to the limit.
// If the calculated value is greater than the total number of records,
// it sets the Last page to the total number of records.
func (p *Pagination) setLast() {
	p.Last = p.GetOffset() + p.Size
	if p.Last > p.Records {
		p.Last = p.Records
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

func New(model *gorm.DB, request *evo.Request, out ...interface{}) (*Pagination, error) {
	var err error
	var p = Pagination{}
	var limit = request.Query("limit").Int()
	var page = request.Query("page").Int()
	if limit < 10 {
		limit = 10
	}
	if limit > p.MaxSize {
		if p.MaxSize == 0 {
			p.MaxSize = 50
		}
		limit = p.MaxSize
	}
	if page < 1 {
		page = 1
	}
	p.setCurrentPage(page)

	if len(out) > 0 {
		n, err := p.LoadData(out[0])
		return n, err
	}
	return &p, err
}

func (p *Pagination) LoadData(out interface{}) (*Pagination, error) {
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
