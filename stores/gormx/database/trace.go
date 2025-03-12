package database

import (
	"github.com/hanyougame/glib/tracing"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func registerTraceHook(tx *gorm.DB) {
	tx.Callback().Create().Before("gorm:created").Register("trace:create", func(db *gorm.DB) {
		traceSql("gorm:create", db)
	})
	tx.Callback().Create().After("gorm:saved").Register("trace:save", func(db *gorm.DB) {
		traceSql("gorm:save", db)
	})
	tx.Callback().Query().After("gorm:queried").Register("trace:query", func(db *gorm.DB) {
		traceSql("gorm:query", db)
	})
	tx.Callback().Delete().After("gorm:deleted").Register("trace:delete", func(db *gorm.DB) {
		traceSql("gorm:delete", db)
	})
	tx.Callback().Update().After("gorm:updated").Register("trace:update", func(db *gorm.DB) {
		traceSql("gorm:update", db)
	})
	tx.Callback().Raw().After("*").Register("trace:raw", func(db *gorm.DB) {
		traceSql("gorm:raw", db)
	})
	tx.Callback().Row().After("*").Register("trace:row", func(db *gorm.DB) {
		traceSql("gorm:row", db)
	})
}

func traceSql(spanName string, db *gorm.DB) {
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	logx.WithContext(db.Statement.Context).Infof("spanName : %s, sql:%s", spanName, sql)
	tracing.Inject(db.Statement.Context, spanName, func(span oteltrace.Span) oteltrace.Span {
		span.SetAttributes(attribute.KeyValue{
			Key:   "gorm.sql",
			Value: attribute.StringValue(sql),
		})
		return span
	})
}
