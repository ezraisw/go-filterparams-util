# go-filterparams-util

Utilities for parsing and cleaning go-filterparams query.

Currently in an experimental state. Do expect bugs!

## Cleaning

```go
func GetCleanSpecs(v interface{}) []CleanSpec
func CleanQuery(queryData *filterparams.QueryData, specs ...CleanSpec) *filterparams.QueryData
func CleanFilter(filter interface{}, specs ...CleanSpec) interface{}
func CleanOrders(orders []*definition.Order, specs ...CleanSpec) []*definition.Order
```

### Example

Defining your own cleaning specification.
```go
package main

import (
    "github.com/pwnedgod/go-filterparams-util"
)

func main() {
    //... obtain your query data
    var queryData *filterparams.QueryData

    queryData = fputil.CleanQuery(
        queryData,
        CleanSpec{Name: "username", DataType: fputil.StringDataType()},
        CleanSpec{Name: "email", DataType: fputil.StringDataType()},
        CleanSpec{Name: "age", DataType: fputil.IntDataType(0)},
    )

    //... use your query data
}
```

Or, you could use a struct for the specification. Do note that this **does not** transform the names from query parameter and is only used for validation.
```go
package main

import (
    "github.com/pwnedgod/go-filterparams-util"
)

type User struct {
    Name        string    `filterparams:"userName"`
    Email       string    `filterparams:"email"`
    Age         int       `filterparams:"age"`
    Appointment time.Time `filterparams:"appointmentTime"`
    InternalVal bool      `filterparams:"-"` // Ignore
}

var specs = GetCleanSpecs(new(User))

func main() {
    //... obtain your query data
    var queryData *filterparams.QueryData

    queryData = fputil.CleanQuery(queryData, specs...)
    //... use your query data
}
```
