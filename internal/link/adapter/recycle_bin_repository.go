package adapter

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"shortlink/internal/common/constant"
	"shortlink/internal/link/adapter/po"
	"shortlink/internal/link/adapters/assembler"
	"shortlink/internal/link/domain/entity"
)

type RecycleBinRepository struct {
	db  *gorm.DB
	rdb *redis.Client
	cvt assembler.LinkRecycleBinConverter
}

func NewRecycleBinRepository(db *gorm.DB, rdb *redis.Client) RecycleBinRepository {
	return RecycleBinRepository{
		db:  db,
		rdb: rdb,
		cvt: assembler.LinkRecycleBinConverter{},
	}
}

// SaveToRecycleBin 保存到回收站
func (r RecycleBinRepository) SaveToRecycleBin(
	ctx context.Context,
	id entity.LinkID,
	updateFn func(ctx context.Context, link *entity.Link) (*entity.Link, error),
) (err error) {

	// 持久化操作
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询
		var linkPo po.Link
		if err = r.db.WithContext(ctx).Model(&linkPo).
			Where("gid = ? and fullShortUrl = ?", id.Gid, id.FullShortUrl).First(&linkPo).Error; err != nil {
			return err
		}

		// 执行更新操作
		link := r.cvt.LinkPoToEntity(linkPo)
		var updatedLink *entity.Link
		if updatedLink, err = updateFn(ctx, &link); err != nil {
			return err
		}
		updatedLinkPo := r.cvt.LinkEntityToPo(*updatedLink)

		// 更新
		if err = r.db.WithContext(ctx).Model(&updatedLinkPo).Save(&linkPo).Error; err != nil {
			return err
		}
		return nil
	})

	// 删除缓存
	r.rdb.Del(ctx, fmt.Sprintf(constant.GotoShortLinkKey, id.FullShortUrl))

	return
}

// RemoveFromRecycleBin 从回收站移除
func (r RecycleBinRepository) RemoveFromRecycleBin(
	ctx context.Context,
	id entity.LinkID,
	enableStatus int,
) (err error) {

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		var linkPo po.Link
		err = r.db.WithContext(ctx).Model(&linkPo).
			Where("gid = ? and fullShortUrl = ? and enableStatus = ?", id.Gid, id.FullShortUrl, enableStatus).
			Find(&linkPo).Error
		if err != nil {
			return
		}
		err = r.db.Delete(&linkPo).Error
		return
	})
	return
}

// RecoverFromRecycleBin 从回收站中恢复
func (r RecycleBinRepository) RecoverFromRecycleBin(
	ctx context.Context,
	id entity.LinkID,
	enableStatus int,
	updateFn func(ctx context.Context, link *entity.Link) (*entity.Link, error),
) (err error) {

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		var linkPo po.Link
		err = r.db.WithContext(ctx).Model(&linkPo).
			Where("gid = ? and fullShortUrl = ? and enableStatus = ?", id.Gid, id.FullShortUrl, enableStatus).
			Find(&linkPo).Error
		if err != nil {
			return
		}

		link := r.cvt.LinkPoToEntity(linkPo)
		var updatedLink *entity.Link
		if updatedLink, err = updateFn(ctx, &link); err != nil {
			return err
		}
		updatedLinkPo := r.cvt.LinkEntityToPo(*updatedLink)

		err = r.db.Save(&updatedLinkPo).Error
		return
	})
	if err != nil {
		return
	}

	err = r.rdb.Del(ctx, fmt.Sprintf(constant.GotoIsNullShortLinkKey, id.FullShortUrl)).Err()
	return
}
