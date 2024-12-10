package gormx

import (
	"github.com/hanyougame/glib/tracing"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type GlobalHooks struct{}

// AfterCreate 全局钩子：创建后
func AfterCreate(db *gorm.DB) {
	var sql = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	tracing.Inject(db.Statement.Context, "gorm.create", func(span oteltrace.Span) oteltrace.Span {
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "gorm.sql",
				Value: attribute.StringValue(sql),
			},
		)
		return span
	})
}

// AfterUpdate 全局钩子：更新后
func AfterUpdate(db *gorm.DB) {
	var sql = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	tracing.Inject(db.Statement.Context, "gorm.update", func(span oteltrace.Span) oteltrace.Span {
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "gorm.sql",
				Value: attribute.StringValue(sql),
			},
		)
		return span
	})
}

// AfterSave 全局钩子：保存后
func AfterSave(db *gorm.DB) {
	var sql = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	tracing.Inject(db.Statement.Context, "gorm.save", func(span oteltrace.Span) oteltrace.Span {
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "gorm.sql",
				Value: attribute.StringValue(sql),
			},
		)
		return span
	})
}

// AfterDelete 全局钩子：删除后
func AfterDelete(db *gorm.DB) {
	var sql = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	tracing.Inject(db.Statement.Context, "gorm.delete", func(span oteltrace.Span) oteltrace.Span {
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "gorm.sql",
				Value: attribute.StringValue(sql),
			},
		)
		return span
	})
}

func AfterQuery(db *gorm.DB) {
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	tracing.Inject(db.Statement.Context, "gorm.query", func(span oteltrace.Span) oteltrace.Span {
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "gorm.sql",
				Value: attribute.StringValue(sql),
			},
		)
		return span
	})
}
