package filter

import (
	"errors"

	"github.com/huandu/go-sqlbuilder"
)

type Filter interface {
	Apply(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder
	Validate() error
}

type Filters []Filter

func (f Filters) Validate() error {
	var errs error

	for _, filter := range f {
		if err := filter.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (f Filters) Apply(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder {
	for _, filter := range f {
		sb = filter.Apply(sb)
	}
	return sb
}
