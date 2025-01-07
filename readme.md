# Pagination Library for EVO

This library provides a robust and flexible solution for handling pagination in Go applications. Built using GORM for database interaction and EVO for request handling, it is designed to simplify pagination logic, making it easy to retrieve and display paginated data in a structured and efficient way.

---

## Why This Library?

Pagination is essential for applications that deal with large datasets. Fetching all data at once is not practical for performance and usability reasons. This library offers:

- **Ease of Use**: Simplifies the implementation of pagination in APIs and applications.
- **Flexibility**: Provides customization options for page size, maximum size, and debug mode.
- **Efficiency**: Uses GORM’s database capabilities to fetch only the required data.
- **Structured Output**: Returns data with metadata like total records, current page, total pages, and more.

---

## Features

- Automatically calculates offsets and limits for database queries.
- Handles edge cases like invalid page numbers or page sizes.
- Supports debugging for SQL queries.
- Offers metadata such as total records, total pages, and page range.
- Easily integrates with GORM and EVO.

---

## Installation

Add this library to your Go project:

```sh
go get github.com/getevo/pagination
```

---

## Usage

### Basic Example

This example demonstrates how to paginate user orders:

```go
func (c Controller) getUserOrdersHandler(request *evo.Request) any {
if request.User().Anonymous() {
return errors.New(403, "access denied")
}
var orders []models.Order
var model = db.Where("uuid = ?", request.User().UUID())
var p, err = pagination.New(model, request, &orders, pagination.Options{MaxSize: 50})
if err != nil {
log.Error(err)
}
return p
}
```

### Customizing Page and Size

You can customize the page and size directly in the options:

```go
var p, err = pagination.New(model, request, &data, pagination.Options{
Page:    2,
Size:    20,
MaxSize: 100,
})
```

### Debug Mode

Enable debug mode to see detailed SQL logs:

```go
var p, err = pagination.New(model, request, &data, pagination.Options{
Debug: true,
})
```

### Handling Errors

The library gracefully handles errors. If an error occurs during data fetching, the response will include the error message:

```go
if err != nil {
return errors.New(500, "Failed to fetch paginated data")
}
```

### Accessing Metadata

The `Pagination` struct contains metadata you can use in your response:

```go
response := map[string]interface{}{
"total_records": p.Records,
"current_page": p.CurrentPage,
"total_pages":  p.Pages,
"data":         p.Data,
}
return response
```

### Custom Query Examples

#### Paginate Products with a Filter

```go
var products []models.Product
var model = db.Where("category = ?", "electronics")
var p, err = pagination.New(model, request, &products, pagination.Options{MaxSize: 20})
if err != nil {
log.Error(err)
}
return p
```

#### Paginate Logs with Date Range

```go
var logs []models.Log
var model = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
var p, err = pagination.New(model, request, &logs, pagination.Options{Size: 15, MaxSize: 50})
if err != nil {
log.Error(err)
}
return p
```

#### Paginate with Sorting

```go
var users []models.User
var model = db.Order("created_at DESC")
var p, err = pagination.New(model, request, &users)
if err != nil {
log.Error(err)
}
return p
```

---

## API Reference

### `Options`

The `Options` struct allows you to customize pagination behavior:

| Field     | Type   | Description                         |
| --------- | ------ | ----------------------------------- |
| `Size`    | `int`  | Number of records per page.         |
| `Page`    | `int`  | Current page number.                |
| `MaxSize` | `int`  | Maximum allowed page size.          |
| `Debug`   | `bool` | Enables debug mode for SQL queries. |

### `Pagination`

| Field         | Type          | Description                       |
| ------------- | ------------- | --------------------------------- |
| `Records`     | `int`         | Total number of records.          |
| `CurrentPage` | `int`         | Current page number.              |
| `Pages`       | `int`         | Total number of pages.            |
| `Size`        | `int`         | Number of records per page.       |
| `First`       | `int`         | First record on the current page. |
| `Last`        | `int`         | Last record on the current page.  |
| `Data`        | `interface{}` | The paginated data.               |

### Methods

#### `New`

Creates a new instance of the pagination object.

```go
func New(model *gorm.DB, request *evo.Request, out interface{}, options ...Options) (*Pagination, error)
```

#### `LoadData`

Loads paginated data into the output variable.

```go
func (p *Pagination) LoadData(out interface{}) (*Pagination, error)
```

#### `GetOffset`

Calculates the offset for paginated data.

```go
func (p *Pagination) GetOffset() int
```

#### `GetResponse`

Returns a standardized response object.

```go
func (p *Pagination) GetResponse() outcome.Response
```

---

## Example Using cURL

You can test the pagination API by adjusting the `size` and `page` parameters in the query string. For example:

```sh
curl --location 'https://api.example.com/api/user/orders?size=10&page=2'
```



