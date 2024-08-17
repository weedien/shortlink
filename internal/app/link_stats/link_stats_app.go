package link_stats

import "shortlink/internal/app/link_stats/query"

type Application struct {
	Queries Queries
}

type Queries struct {
	GetLinkStats               query.GetLinkStatsHandler
	GroupLinkStats             query.GroupLinkStatsHandler
	GetLinkStatsAccessRecord   query.GetLinkStatsAccessRecordHandler
	GroupLinkStatsAccessRecord query.GroupLinkStatsAccessRecord
}
