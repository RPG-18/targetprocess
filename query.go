package targetprocess

import (
	"net/url"
	"strconv"
	"strings"
)

type orderEnum int

const (
	format           = `json`
	asc    orderEnum = iota
	desc
)

var (
	defaultQuery = NewQuery()
)

type Query struct {
	where     string
	include   string
	exclude   string
	append    string
	skip      int32
	take      int32
	innerTake int32
	order     Order
	format    string
}

func (q *Query) Where(str string) *Query {
	q.where = str
	return q
}

func (q *Query) Inclde(fields ...string) *Query {
	q.include = "[" + strings.Join(fields, ",") + "]"
	return q
}

func Where(str string) *Query {
	return NewQuery().Where(str)
}

func (q *Query) Order(order Order) *Query {
	q.order = order
	return q
}

func (q *Query) Skip(cnt int32) *Query {
	q.skip = cnt
	return q
}

func (q *Query) Take(cnt int32) *Query {
	q.take = cnt
	return q
}

func (q *Query) InnerTake(cnt int32) *Query {
	q.innerTake = cnt
	return q
}

func NewQuery() *Query {
	return &Query{
		format: format,
	}
}

func (q *Query) values() url.Values {
	v := url.Values{}
	v.Set("format", q.format)

	if len(q.include) != 0 {
		v.Add("include", q.include)
	}

	if q.take != 0 {
		v.Add("take", strconv.Itoa(int(q.take)))
	}

	if q.skip != 0 {
		v.Add("skip", strconv.Itoa(int(q.skip)))
	}

	if q.innerTake != 0 {
		v.Add("innertake", strconv.Itoa(int(q.innerTake)))
	}

	if len(q.order.field) != 0 {
		switch q.order.order {
		case asc:
			v.Add("orderby", q.order.field)
		case desc:
			v.Add("orderbydesc", q.order.field)
		}
	}

	return v
}

type Order struct {
	field string
	order orderEnum
}

func OrderBy(field string) Order {
	return Order{
		field: field,
		order: asc,
	}
}

func OrderByDesc(field string) Order {
	return Order{
		field: field,
		order: desc,
	}
}
