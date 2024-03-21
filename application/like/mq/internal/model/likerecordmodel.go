package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LikeRecordModel = (*customLikeRecordModel)(nil)

type (
	// LikeRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeRecordModel.
	LikeRecordModel interface {
		likeRecordModel
		withSession(session sqlx.Session) LikeRecordModel
	}

	customLikeRecordModel struct {
		*defaultLikeRecordModel
	}
)

// NewLikeRecordModel returns a model for the database table.
func NewLikeRecordModel(conn sqlx.SqlConn) LikeRecordModel {
	return &customLikeRecordModel{
		defaultLikeRecordModel: newLikeRecordModel(conn),
	}
}

func (m *customLikeRecordModel) withSession(session sqlx.Session) LikeRecordModel {
	return NewLikeRecordModel(sqlx.NewSqlConnFromSession(session))
}
