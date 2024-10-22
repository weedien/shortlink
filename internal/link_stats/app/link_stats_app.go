package app

import "shortlink/internal/link_stats/app/query"

type Application struct {
	Queries Queries
}

type Queries struct {
	GetLinkStat               query.GetLinkStatHandler
	GroupLinkStat             query.GroupLinkStatHandler
	GetLinkStatAccessRecord   query.GetLinkStatAccessRecordHandler
	GroupLinkStatAccessRecord query.GroupLinkStatAccessRecordHandler
}
