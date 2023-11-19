package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error)
	DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error)
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}

func (G GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return G.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt": gorm.Expr("read_cnt + 1"),
			"utime":    time.Now().UnixMilli(),
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

func (G GORMInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status": 1,
				"utime":  now,
			}),
		},
		).Create(&UserLikeBiz{
			Biz:    biz,
			BizId:  bizId,
			Uid:    uid,
			Ctime:  now,
			Utime:  now,
			Status: 1,
		}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			// MySQL 不写
			//Columns:
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("like_cnt + 1"),
				"utime":    time.Now().UnixMilli(),
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   bizId,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (G GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := G.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ? AND status=?", biz, bizId, uid, 1).First(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil
}

func (G GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).Where("biz = ? AND biz_id = ? AND uid = ?", biz, bizId, uid).Updates(map[string]any{
			"status": 0,
			"utime":  now,
		}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, bizId).Updates(map[string]any{
			"like_cnt": gorm.Expr("like_cnt - 1"),
			"utime":    now,
		}).Error
	})
}

func (G GORMInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var res Interactive
	err := G.db.WithContext(ctx).Where("biz=? AND biz_id=?", biz, bizId).First(&res).Error
	if err != nil {
		return Interactive{}, err
	}
	return res, nil
}

func (G GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	//TODO implement me
	panic("implement me")
}

func (G GORMInteractiveDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}

type Interactive struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 业务标识符
	// 同一个资源，在这里应该只有一行
	// 也就是说我要在 bizId 和 biz 上创建联合唯一索引
	// 1. bizId, biz。优先选择这个，因为 bizId 的区分度更高
	// 2. biz, bizId。如果有 WHERE biz = xx 这种查询条件（不带 bizId）的，就只能这种
	//
	// 联合索引的列的顺序：查询条件，区分度
	// 这个名字无所谓
	BizId int64 `gorm:"uniqueIndex:biz_id_type"`
	// 我这里biz 用的是 string，有些公司枚举使用的是 int 类型
	// 0-article
	// 1- xxx
	// 默认是 BLOB/TEXT 类型
	Biz string `gorm:"uniqueIndex:biz_id_type;type:varchar(128)"`
	// 这个是阅读计数
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}

type UserLikeBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Biz   string `gorm:"uniqueIndex:uid_biz_id_type;type:varchar(128)"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_id_type"`

	// 谁的操作
	Uid int64 `gorm:"uniqueIndex:uid_biz_id_type"`

	Ctime  int64
	Utime  int64
	Status uint8
}

// Collection 收藏夹
type Collection struct {
	Id   int64  `gorm:"primaryKey,autoIncrement"`
	Name string `gorm:"type=varchar(1024)"`
	Uid  int64  `gorm:""`

	Ctime int64
	Utime int64
}

// UserCollectionBiz 收藏的东西
type UserCollectionBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 收藏夹 ID
	// 作为关联关系中的外键，我们这里需要索引
	Cid   int64  `gorm:"index"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	// 这算是一个冗余，因为正常来说，
	// 只需要在 Collection 中维持住 Uid 就可以
	Uid   int64 `gorm:"uniqueIndex:biz_type_id_uid"`
	Ctime int64
	Utime int64
}
