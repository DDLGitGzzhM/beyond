package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LikeCountModel = (*customLikeCountModel)(nil)

type (
	// LikeCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeCountModel.
	LikeCountModel interface {
		likeCountModel
		withSession(session sqlx.Session) LikeCountModel
	}

	customLikeCountModel struct {
		*defaultLikeCountModel
	}
)

// NewLikeCountModel returns a model for the database table.
func NewLikeCountModel(conn sqlx.SqlConn) LikeCountModel {
	return &customLikeCountModel{
		defaultLikeCountModel: newLikeCountModel(conn),
	}
}

func (m *customLikeCountModel) withSession(session sqlx.Session) LikeCountModel {
	return NewLikeCountModel(sqlx.NewSqlConnFromSession(session))
}
