package repo

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"shortlink/common/consts"
	"shortlink/common/types"
	"shortlink/internal/infra/persistence/converter"
	"shortlink/internal/infra/persistence/po"
)

type RecycleBinRepository struct {
	db  *gorm.DB
	rdb *redis.Client
	cvt converter.LinkRecycleBinConverter
}

// SaveToRecycleBin 保存到回收站
func (r RecycleBinRepository) SaveToRecycleBin(
	ctx context.Context,
	id types.LinkID,
	updateFn func(ctx context.Context, link *types.Link) (*types.Link, error),
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
		var updatedLink *types.Link
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
	r.rdb.Del(ctx, fmt.Sprintf(consts.GotoShortLinkKey, id.FullShortUrl))

	return
}

// RemoveFromRecycleBin 从回收站移除
func (r RecycleBinRepository) RemoveFromRecycleBin(
	ctx context.Context,
	id types.LinkID,
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
	id types.LinkID,
	enableStatus int,
	updateFn func(ctx context.Context, link *types.Link) (*types.Link, error),
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
		var updatedLink *types.Link
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

	err = r.rdb.Del(ctx, fmt.Sprintf(consts.GotoIsNullShortLinkKey, id.FullShortUrl)).Err()
	return
}
