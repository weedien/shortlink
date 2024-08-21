package adapters

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/lock"
	"shortlink/internal/common/persistence/po"
	"shortlink/internal/common/toolkit"
	"shortlink/internal/link/domain/valobj"
)

type LinkStatsRepository struct {
	db     *gorm.DB
	rdb    *redis.Client
	locker lock.DistributedLock
}

func NewLinkStatsRepository(db *gorm.DB, rdb *redis.Client) LinkStatsRepository {
	return LinkStatsRepository{db: db, rdb: rdb}
}

func (r LinkStatsRepository) SaveLinkStats(ctx context.Context, statsInfo valobj.ShortLinkStatsRecordVo) error {
	lockKey := constant.LockGidUpdateKey + statsInfo.FullShortUrl
	if _, err := r.locker.Acquire(ctx, lockKey, -1); err != nil {
		return err
	}
	defer func() {
		if err := r.locker.Release(ctx, lockKey); err != nil {
			slog.Error("release lock failed", "lockKey", lockKey, "err", err)
		}
	}()

	fullShortUrl := statsInfo.FullShortUrl
	currentDate := statsInfo.CurrentDate
	hour := currentDate.Hour() + 1
	weekDay := int(currentDate.Weekday())

	// 访问统计
	// 确定两个值的信息，uvFirstFlag 和 uipFirstFlag
	uv, uip := 0, 0
	uvAdded, err := r.rdb.SAdd(ctx, constant.ShortLinkStatsUvKey+fullShortUrl, statsInfo.UV).Result()
	if err != nil {
		return err
	}
	if uvAdded > 0 {
		uv = 1
	}
	uipAdded, err := r.rdb.SAdd(ctx, constant.ShortLinkStatsUipKey+fullShortUrl, statsInfo.RemoteAddr).Result()
	if err != nil {
		return err
	}
	if uipAdded > 0 {
		uip = 1
	}
	linkAccessStatsPo := po.LinkAccessStats{
		Pv:           1,
		Uv:           uv,
		Uip:          uip,
		Hour:         hour,
		Weekday:      weekDay,
		FullShortURL: fullShortUrl,
		Date:         currentDate,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "full_short_url"}, {Name: "date"}, {Name: "hour"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"pv":  gorm.Expr("pv + ?", 1),
			"uv":  gorm.Expr("uv + ?", uv),
			"uip": gorm.Expr("uip + ?", uip),
		}),
	}).Create(&linkAccessStatsPo).Error; err != nil {
		return err
	}
	// 今日统计
	linkStatsToday := po.LinkStatsToday{
		TodayPv:      1,
		TodayUv:      uv,
		TodayUip:     uip,
		Date:         currentDate,
		FullShortURL: fullShortUrl,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "full_short_url"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"today_pv":  gorm.Expr("today_pv + ?", 1),
			"today_uv":  gorm.Expr("today_uv + ?", uv),
			"today_uip": gorm.Expr("today_uip + ?", uip),
		}),
	}).Create(&linkStatsToday).Error; err != nil {
		return err
	}
	// 地区信息
	location := toolkit.GetLocationByIP(statsInfo.RemoteAddr)
	linkLocaleStatsPo := po.LinkLocaleStats{
		FullShortURL: fullShortUrl,
		Date:         currentDate,
		Cnt:          1,
		Province:     location.RegionName,
		City:         location.City,
		Country:      location.Country,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "full_short_url"}, {Name: "date"}, {Name: "province"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"cnt": gorm.Expr("cnt + ?", 1)}),
	}).Create(&linkLocaleStatsPo).Error; err != nil {
		return err
	}
	// 操作系统信息
	linkOsStatsPo := po.LinkOsStats{
		Os:           statsInfo.OS,
		Cnt:          1,
		FullShortURL: fullShortUrl,
		Date:         currentDate,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "os"}, {Name: "full_short_url"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"cnt": gorm.Expr("cnt + ?", 1)}),
	}).Create(&linkOsStatsPo).Error; err != nil {
		return err
	}
	// 浏览器信息
	linkBrowserStatsPo := po.LinkBrowserStats{
		Browser:      statsInfo.Browser,
		Cnt:          1,
		FullShortURL: fullShortUrl,
		Date:         currentDate,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "browser"}, {Name: "full_short_url"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"cnt": gorm.Expr("cnt + ?", 1)}),
	}).Create(&linkBrowserStatsPo).Error; err != nil {
		return err
	}
	// 设备信息
	linkDeviceStatsPo := po.LinkDeviceStats{
		Device:       statsInfo.Device,
		Cnt:          1,
		FullShortURL: fullShortUrl,
		Date:         currentDate,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device"}, {Name: "full_short_url"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"cnt": gorm.Expr("cnt + ?", 1)}),
	}).Create(&linkDeviceStatsPo).Error; err != nil {
		return err
	}
	// 网络信息
	linkNetworkStatsPo := po.LinkNetworkStats{
		Network:      statsInfo.Network,
		Cnt:          1,
		FullShortURL: fullShortUrl,
		Date:         currentDate,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "network"}, {Name: "full_short_url"}, {Name: "date"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"cnt": gorm.Expr("cnt + ?", 1)}),
	}).Create(&linkNetworkStatsPo).Error; err != nil {
		return err
	}
	// 访问日志
	linkAccessLogPo := po.LinkAccessLog{
		FullShortURL: fullShortUrl,
		User:         statsInfo.UV,
		IP:           statsInfo.RemoteAddr,
		Browser:      statsInfo.Browser,
		Os:           statsInfo.OS,
		Network:      statsInfo.Network,
		Device:       statsInfo.Device,
		Locale:       location.Country + "-" + location.RegionName + "-" + location.City,
	}
	if err := r.db.Create(&linkAccessLogPo).Error; err != nil {
		return err
	}
	// 更新shortLink表中的状态pv, uv, uip
	linkGotoPo := po.LinkGoto{FullShortURL: fullShortUrl}
	if err := r.db.First(&linkGotoPo).Error; err != nil {
		return err
	}
	r.db.Model(&po.Link{}).
		Where("gid = ? and full_short_url = ?", linkGotoPo.Gid, fullShortUrl).
		Updates(map[string]interface{}{
			"total_pv":  gorm.Expr("total_pv + ?", 1),
			"total_uv":  gorm.Expr("total_uv + ?", uv),
			"total_uip": gorm.Expr("total_uip + ?", uip),
		})

	return nil
}
