package readrepo

import (
	"context"
	"gorm.io/gorm"
	"math"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/types"
	po2 "shortlink/internal/link_stats/adapter/po"
	"shortlink/internal/link_stats/adapters/readrepo/dao"
	"shortlink/internal/link_stats/app/query"
	"shortlink/internal/user/adapter/po"
)

type LinkStatsQuery struct {
	linkAccessStatsDao  dao.LinkAccessStatsDao
	linkAccessLogsDao   dao.LinkAccessLogsDao
	linkLocaleStatsDao  dao.LinkLocaleStatsDao
	linkBrowserStatsDao dao.LinkBrowserStatsDao
	linkOsStatsDao      dao.LinkOsStatsDao
	linkDeviceStatsDao  dao.LinkDeviceStatsDao
	linkNetworkStatsDao dao.LinkNetworkStatsDao
}

func NewLinkStatsQuery(db *gorm.DB) LinkStatsQuery {
	return LinkStatsQuery{
		linkAccessStatsDao:  dao.NewLinkAccessStatsDao(db),
		linkAccessLogsDao:   dao.NewLinkAccessLogsDao(db),
		linkLocaleStatsDao:  dao.NewLinkLocaleStatsDao(db),
		linkBrowserStatsDao: dao.NewLinkBrowserStatsDao(db),
		linkOsStatsDao:      dao.NewLinkOsStatsDao(db),
		linkDeviceStatsDao:  dao.NewLinkDeviceStatsDao(db),
		linkNetworkStatsDao: dao.NewLinkNetworkStatsDao(db),
	}
}

// GetLinkStats 获取单个短链接监控数据
func (q LinkStatsQuery) GetLinkStats(ctx context.Context, param query.GetLinkStats) (res *query.LinkStats, err error) {

	queryParam := dao.LinkQueryParam{
		FullShortUrl: param.FullShortUrl,
		Gid:          param.Gid,
		EnableStatus: param.EnableStatus,
		StartDate:    param.StartDate,
		EndDate:      param.EndDate,
	}

	var stats []po.LinkAccessStats
	stats, err = q.linkAccessStatsDao.ListStatsByLink(ctx, queryParam)
	if err != nil {
		return
	}
	if len(stats) == 0 {
		return
	}
	// 基础访问数据
	pvUvUidStats, err := q.linkAccessLogsDao.FindPvUvUidStatsByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	// 基础访问详情
	daily := make([]query.LinkStatsAccessDaily, 0)
	var rangeDates []string
	for d := param.StartDate; !d.After(param.EndDate); d = d.AddDate(0, 0, 1) {
		rangeDates = append(rangeDates, d.Format("2006-01-02"))
	}
	statsMap := make(map[string]po.LinkAccessStats)
	for _, item := range stats {
		statsMap[item.Date.Format("2006-01-02")] = item
	}
	for _, date := range rangeDates {
		if item, found := statsMap[date]; found {
			accessDailyRespDTO := query.LinkStatsAccessDaily{
				Date: date,
				Pv:   item.Pv,
				Uv:   item.Uv,
				Uip:  item.Uip,
			}
			daily = append(daily, accessDailyRespDTO)
		} else {
			accessDailyRespDTO := query.LinkStatsAccessDaily{
				Date: date,
				Pv:   0,
				Uv:   0,
				Uip:  0,
			}
			daily = append(daily, accessDailyRespDTO)
		}
	}
	// 地区访问详情（仅国内）
	locales := make([]query.LinkStatsLocale, 0)
	localeStats, err := q.linkLocaleStatsDao.ListLocaleByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var localeCnTotal int
	for _, item := range localeStats {
		localeCnTotal += item.Cnt
	}
	for _, item := range localeStats {
		ratio := float64(item.Cnt) / float64(localeCnTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		locale := query.LinkStatsLocale{
			Cnt:    item.Cnt,
			Locale: item.Province,
			Ratio:  actualRatio,
		}
		locales = append(locales, locale)
	}
	// 小时访问详情
	hours := make([]int, 24)
	hourStats, err := q.linkAccessStatsDao.ListHourStatsByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	for _, item := range hourStats {
		hours[item.Hour] = item.Pv
	}
	// 高频访问IP详情
	var topIps []query.LinkStatsTopIp
	topIpStats, err := q.linkAccessLogsDao.ListTopIpByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	for _, item := range topIpStats {
		topIp := query.LinkStatsTopIp{
			Ip:  item.Ip,
			Cnt: item.Cnt,
		}
		topIps = append(topIps, topIp)
	}
	// 一周访问详情
	weekdays := make([]int, 7)
	weekdayStats, err := q.linkAccessStatsDao.ListWeekdayStatsByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	for _, item := range weekdayStats {
		weekdays[item.Weekday] = item.Pv
	}
	// 浏览器访问情况
	browsers := make([]query.LinkStatsBrowser, 0)
	browserStats, err := q.linkBrowserStatsDao.ListBrowserStatsByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var browserTotal int
	for _, item := range browserStats {
		browserTotal += item.Cnt
	}
	for _, item := range browserStats {
		ratio := float64(item.Cnt) / float64(browserTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		browser := query.LinkStatsBrowser{
			Browser: item.Browser,
			Cnt:     item.Cnt,
			Ratio:   actualRatio,
		}
		browsers = append(browsers, browser)
	}
	// 操作系统访问详情
	oss := make([]query.LinkStatsOs, 0)
	osStats, err := q.linkOsStatsDao.ListOsStatsByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var osTotal int
	for _, item := range osStats {
		osTotal += item.Cnt
	}
	for _, item := range osStats {
		ratio := float64(item.Cnt) / float64(osTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		os := query.LinkStatsOs{
			Os:    item.Os,
			Cnt:   item.Cnt,
			Ratio: actualRatio,
		}
		oss = append(oss, os)
	}
	// 访客访问类型详情
	uvTypes := make([]query.LinkStatsUv, 2)
	uvTypeStats, err := q.linkAccessLogsDao.FindUvTypeCntByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	oldUserCnt := uvTypeStats.OldUserCnt
	newUserCnt := uvTypeStats.NewUserCnt
	uvTotal := oldUserCnt + newUserCnt
	oldUserRatio := float64(oldUserCnt) / float64(uvTotal)
	newUserRatio := float64(newUserCnt) / float64(uvTotal)
	oldUserRatio = math.Round(oldUserRatio*100.0) / 100.0
	newUserRatio = math.Round(newUserRatio*100.0) / 100.0
	oldUser := query.LinkStatsUv{
		VisitorType: "老访客",
		Cnt:         oldUserCnt,
		Ratio:       oldUserRatio,
	}
	newUser := query.LinkStatsUv{
		VisitorType: "新访客",
		Cnt:         newUserCnt,
		Ratio:       newUserRatio,
	}
	uvTypes = append(uvTypes, oldUser, newUser)
	// 访问设备类型详情
	devices := make([]query.LinkStatsDevice, 0)
	deviceStats, err := q.linkDeviceStatsDao.ListDeviceStatsByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var deviceTotal int
	for _, item := range deviceStats {
		deviceTotal += item.Cnt
	}
	for _, item := range deviceStats {
		ratio := float64(item.Cnt) / float64(deviceTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		device := query.LinkStatsDevice{
			Device: item.Device,
			Cnt:    item.Cnt,
			Ratio:  actualRatio,
		}
		devices = append(devices, device)
	}
	// 访问网络类型详情
	networks := make([]query.LinkStatsNetwork, 0)
	networkStats, err := q.linkNetworkStatsDao.ListNetworkStatsByLink(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var networkTotal int
	for _, item := range networkStats {
		networkTotal += item.Cnt
	}
	for _, item := range networkStats {
		ratio := float64(item.Cnt) / float64(networkTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		network := query.LinkStatsNetwork{
			Network: item.Network,
			Cnt:     item.Cnt,
			Ratio:   actualRatio,
		}
		networks = append(networks, network)
	}
	// 组装返回数据
	res = &query.LinkStats{
		Pv:               pvUvUidStats.Pv,
		Uv:               pvUvUidStats.Uv,
		Uip:              pvUvUidStats.Uip,
		Hourly:           hours,
		Daily:            daily,
		Weekly:           weekdays,
		LocationCnStats:  locales,
		TopIpStats:       topIps,
		BrowserStats:     browsers,
		OsStats:          oss,
		VisitorTypeStats: uvTypes,
		DeviceStats:      devices,
		NetworkStats:     networks,
	}
	return
}

// GroupLinkStats 获取分组短链接监控数据
func (q LinkStatsQuery) GroupLinkStats(ctx context.Context, param query.GroupLinkStats) (res *query.LinkStats, err error) {

	queryParam := dao.LinkGroupQueryParam{
		Gid:          param.Gid,
		EnableStatus: constant.StatusEnable,
		StartDate:    param.StartDate,
		EndDate:      param.EndDate,
	}
	var stats []po.LinkAccessStats
	stats, err = q.linkAccessStatsDao.ListStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	if len(stats) == 0 {
		return
	}
	// 基础访问数据
	pvUvUidStats, err := q.linkAccessLogsDao.FindPvUvUidStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	// 基础访问详情
	daily := make([]query.LinkStatsAccessDaily, 0)
	var rangeDates []string
	for d := param.StartDate; !d.After(param.EndDate); d = d.AddDate(0, 0, 1) {
		rangeDates = append(rangeDates, d.Format("2006-01-02"))
	}
	statsMap := make(map[string]po.LinkAccessStats)
	for _, item := range stats {
		statsMap[item.Date.Format("2006-01-02")] = item
	}
	for _, date := range rangeDates {
		if item, found := statsMap[date]; found {
			accessDailyRespDTO := query.LinkStatsAccessDaily{
				Date: date,
				Pv:   item.Pv,
				Uv:   item.Uv,
				Uip:  item.Uip,
			}
			daily = append(daily, accessDailyRespDTO)
		} else {
			accessDailyRespDTO := query.LinkStatsAccessDaily{
				Date: date,
				Pv:   0,
				Uv:   0,
				Uip:  0,
			}
			daily = append(daily, accessDailyRespDTO)
		}
	}
	// 地区访问详情（仅国内）
	locales := make([]query.LinkStatsLocale, 0)
	localeStats, err := q.linkLocaleStatsDao.ListLocaleByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var localeCnTotal int
	for _, item := range localeStats {
		localeCnTotal += item.Cnt
	}
	for _, item := range localeStats {
		ratio := float64(item.Cnt) / float64(localeCnTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		locale := query.LinkStatsLocale{
			Cnt:    item.Cnt,
			Locale: item.Province,
			Ratio:  actualRatio,
		}
		locales = append(locales, locale)
	}
	// 小时访问详情
	hours := make([]int, 24)
	hourStats, err := q.linkAccessStatsDao.ListHourStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	for _, item := range hourStats {
		hours[item.Hour] = item.Pv
	}
	// 高频访问IP详情
	var topIps []query.LinkStatsTopIp
	topIpStats, err := q.linkAccessLogsDao.ListTopIpByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	for _, item := range topIpStats {
		topIp := query.LinkStatsTopIp{
			Ip:  item.Ip,
			Cnt: item.Cnt,
		}
		topIps = append(topIps, topIp)
	}
	// 一周访问详情
	weekdays := make([]int, 7)
	weekdayStats, err := q.linkAccessStatsDao.ListWeekdayStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	for _, item := range weekdayStats {
		weekdays[item.Weekday] = item.Pv
	}
	// 浏览器访问情况
	browsers := make([]query.LinkStatsBrowser, 0)
	browserStats, err := q.linkBrowserStatsDao.ListBrowserStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var browserTotal int
	for _, item := range browserStats {
		browserTotal += item.Cnt
	}
	for _, item := range browserStats {
		ratio := float64(item.Cnt) / float64(browserTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		browser := query.LinkStatsBrowser{
			Browser: item.Browser,
			Cnt:     item.Cnt,
			Ratio:   actualRatio,
		}
		browsers = append(browsers, browser)
	}
	// 操作系统访问详情
	oss := make([]query.LinkStatsOs, 0)
	osStats, err := q.linkOsStatsDao.ListOsStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var osTotal int
	for _, item := range osStats {
		osTotal += item.Cnt
	}
	for _, item := range osStats {
		ratio := float64(item.Cnt) / float64(osTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		os := query.LinkStatsOs{
			Os:    item.Os,
			Cnt:   item.Cnt,
			Ratio: actualRatio,
		}
		oss = append(oss, os)
	}
	// 访问设备类型详情
	devices := make([]query.LinkStatsDevice, 0)
	deviceStats, err := q.linkDeviceStatsDao.ListDeviceStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var deviceTotal int
	for _, item := range deviceStats {
		deviceTotal += item.Cnt
	}
	for _, item := range deviceStats {
		ratio := float64(item.Cnt) / float64(deviceTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		device := query.LinkStatsDevice{
			Device: item.Device,
			Cnt:    item.Cnt,
			Ratio:  actualRatio,
		}
		devices = append(devices, device)
	}
	// 访问网络类型详情
	networks := make([]query.LinkStatsNetwork, 0)
	networkStats, err := q.linkNetworkStatsDao.ListNetworkStatsByGroup(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	var networkTotal int
	for _, item := range networkStats {
		networkTotal += item.Cnt
	}
	for _, item := range networkStats {
		ratio := float64(item.Cnt) / float64(networkTotal)
		actualRatio := math.Round(ratio*100.0) / 100.0
		network := query.LinkStatsNetwork{
			Network: item.Network,
			Cnt:     item.Cnt,
			Ratio:   actualRatio,
		}
		networks = append(networks, network)
	}
	// 组装返回数据
	res = &query.LinkStats{
		Pv:              pvUvUidStats.Pv,
		Uv:              pvUvUidStats.Uv,
		Uip:             pvUvUidStats.Uip,
		Hourly:          hours,
		Daily:           daily,
		Weekly:          weekdays,
		LocationCnStats: locales,
		TopIpStats:      topIps,
		BrowserStats:    browsers,
		OsStats:         oss,
		DeviceStats:     devices,
		NetworkStats:    networks,
	}
	return
}

// GetLinkStatsAccessRecord 访问单个短链接指定时间内访问记录监控数据
func (q LinkStatsQuery) GetLinkStatsAccessRecord(
	ctx context.Context,
	param query.GetLinkStatsAccessRecord,
) (res *types.PageResp[query.LinkStatsAccessRecord], err error) {

	queryParam := dao.LinkQueryParam{
		FullShortUrl: param.FullShortUrl,
		Gid:          param.Gid,
		EnableStatus: constant.StatusEnable,
		StartDate:    param.StartDate,
		EndDate:      param.EndDate,
	}

	logPoPage, err := q.linkAccessLogsDao.Page(ctx, queryParam, param.Current, param.Size)

	if err != nil {
		return nil, err
	}

	if logPoPage.Total == 0 {
		return types.NewEmptyPageResp[query.LinkStatsAccessRecord](), nil
	}

	return q.buildStatsAccessRecordResult(logPoPage, func(users []string) (userTypes []dao.UserType, err error) {
		return q.linkAccessLogsDao.SelectUvTypeByUsers(ctx, queryParam, users)
	})
}

// GroupLinkStatsAccessRecord 访问分组短链接指定时间内访问记录监控数据
func (q LinkStatsQuery) GroupLinkStatsAccessRecord(
	ctx context.Context,
	param query.GroupLinkStatsAccessRecord,
) (res *types.PageResp[query.LinkStatsAccessRecord], err error) {

	queryParam := dao.LinkGroupQueryParam{
		Gid:          param.Gid,
		EnableStatus: constant.StatusEnable,
		StartDate:    param.StartDate,
		EndDate:      param.EndDate,
	}

	logPoPage, err := q.linkAccessLogsDao.PageGroup(ctx, queryParam, param.Current, param.Size)
	if err != nil {
		return
	}

	if logPoPage.Total == 0 {
		return types.NewEmptyPageResp[query.LinkStatsAccessRecord](), nil
	}

	return q.buildStatsAccessRecordResult(logPoPage, func(users []string) (userTypes []dao.UserType, err error) {
		return q.linkAccessLogsDao.SelectGroupUvTypeByUsers(ctx, queryParam, users)
	})
}

func (q LinkStatsQuery) buildStatsAccessRecordResult(
	logPoPage *types.PageResp[po2.LinkAccessLog],
	getUserTypeFn func(users []string) (userTypes []dao.UserType, err error),
) (res *types.PageResp[query.LinkStatsAccessRecord], err error) {

	// 构建用户信息列表
	logPos := logPoPage.Records
	users := make([]string, logPoPage.Total)
	for idx, logPo := range logPos {
		users[idx] = logPo.User
	}
	// 获取用户类型
	var userTypes []dao.UserType
	userTypes, err = getUserTypeFn(users)
	if err != nil {
		return
	}
	// 构建map用于查找
	userTypeMap := make(map[string]dao.UserType)
	for _, userType := range userTypes {
		userTypeMap[userType.User] = userType
	}

	// 分页结果类型转换
	res = types.ConvertRecords(logPoPage, func(logPo po2.LinkAccessLog) query.LinkStatsAccessRecord {
		record := query.LinkStatsAccessRecord{
			Browser:    logPo.Browser,
			Os:         logPo.Os,
			Ip:         logPo.IP,
			Network:    logPo.Network,
			Device:     logPo.Device,
			Locale:     logPo.Locale,
			User:       logPo.User,
			AccessTime: logPo.CreateTime,
		}
		// 加上用户类型信息
		if userType, found := userTypeMap[logPo.User]; found {
			record.UvType = userType.UvType
		}
		return record
	})
	return
}

//func (q LinkStatsQuery) buildStatsAccessRecordResultV1(
//	logPos []po.LinkAccessLog,
//	getUserTypeFn func(users []string) (userTypes []dao.UserType, err error),
//) (records []query.LinkStatsAccessRecord, err error) {
//
//	users := make([]string, len(logPos))
//	for idx, logPo := range logPos {
//		users[idx] = logPo.User
//	}
//
//	var userTypes []dao.UserType
//	userTypes, err = getUserTypeFn(users)
//	if err != nil {
//		return
//	}
//
//	for _, logPo := range logPos {
//		record := query.LinkStatsAccessRecord{
//			Browser:    logPo.Browser,
//			Os:         logPo.Os,
//			Ip:         logPo.IP,
//			Network:    logPo.Network,
//			Device:     logPo.Device,
//			Locale:     logPo.Locale,
//			User:       logPo.User,
//			AccessTime: logPo.CreateTime,
//		}
//		for _, userType := range userTypes {
//			if record.User == userType.User {
//				record.UvType = userType.UvType
//				break
//			}
//		}
//		records = append(records, record)
//	}
//	return
//}
