package admin_field_now

import (
	"time"

	"github.com/ecletus/core"

	"github.com/ecletus/admin"
	"github.com/moisespsena-go/i18n-modular/i18nmod"
	path_helpers "github.com/moisespsena-go/path-helpers"
)

var (
	group = i18nmod.PkgToGroup(path_helpers.GetCalledDir())
)

const (
	FieldName  = "Now"
	TimeFormat = time.Stamp + " GMT -07:00"
)

type Field struct {
	Name  string
	Label string

	// Layout time format layout
	Layout       string
	FormatFunc   func(recorde interface{}, context *core.Context, now time.Time) string
	Location     *time.Location
	LocationFunc func(recorde interface{}, context *core.Context) *time.Location
	NowFunc      func(recorde interface{}, context *core.Context) time.Time
}

func (f Field) Now(recorde interface{}, context *core.Context) time.Time {
	return f.NowFunc(recorde, context)
}

func (f Field) Format(recorde interface{}, context *core.Context) string {
	return f.FormatFunc(recorde, context, f.NowFunc(recorde, context))
}

func (f Field) Setup() Field {
	if f.Name == "" {
		f.Name = FieldName
	}

	if f.FormatFunc == nil {
		if f.Layout == "" {
			f.Layout = TimeFormat
		}
		f.FormatFunc = func(recorde interface{}, context *core.Context, now time.Time) string {
			return now.Format(f.Layout)
		}
	}

	if f.NowFunc == nil {
		if f.LocationFunc != nil {
			f.NowFunc = func(recorde interface{}, context *core.Context) time.Time {
				return time.Now().In(f.LocationFunc(recorde, context))
			}
		} else if f.Location != nil {
			f.NowFunc = func(recorde interface{}, context *core.Context) time.Time {
				return time.Now().In(f.Location)
			}
		} else {
			f.NowFunc = func(recorde interface{}, context *core.Context) time.Time {
				return time.Now()
			}
		}
	}

	if f.Label == "" {
		f.Label = group + "." + FieldName
	}
	return f
}

func (f Field) Apply(res *admin.Resource) *Field {
	res.Meta(&admin.Meta{
		Name:  f.Name,
		Label: f.Label,
		Type:  "string",
		Valuer: func(recorde interface{}, context *core.Context) interface{} {
			return f.Now(recorde, context)
		},
		FormattedValuer: func(recorde interface{}, context *core.Context) interface{} {
			return f.FormatFunc(recorde, context, f.NowFunc(recorde, context))
		},
	})
	return &f
}

func Setup(f *Field, res *admin.Resource) *Field {
	return f.Setup().Apply(res)
}
