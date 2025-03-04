package filter

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Time window based filter
// Supported date formats :
// "May 8, 2009 5:57:51 PM",
// "oct 7, 1970",
// "oct 7, '70",
// "oct. 7, 1970",
// "oct. 7, 70",
// "Mon Jan  2 15:04:05 2006",
// "Mon Jan  2 15:04:05 MST 2006",
// "Mon Jan 02 15:04:05 -0700 2006",
// "Monday, 02-Jan-06 15:04:05 MST",
// "Mon, 02 Jan 2006 15:04:05 MST",
// "Tue, 11 Jul 2017 16:28:13 +0200 (CEST)",
// "Mon, 02 Jan 2006 15:04:05 -0700",
// "Thu, 4 Jan 2018 17:53:36 +0000",
// "Mon Aug 10 15:44:11 UTC+0100 2015",
// "Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)",
// "September 17, 2012 10:09am",
// "September 17, 2012 at 10:09am PST-08",
// "September 17, 2012, 10:10:09",
// "October 7, 1970",
// "October 7th, 1970",
// "12 Feb 2006, 19:17",
// "12 Feb 2006 19:17",
// "7 oct 70",
// "7 oct 1970",
// "03 February 2013",
// "1 July 2013",
// "2013-Feb-03",
// //   mm/dd/yy
// "3/31/2014",
// "03/31/2014",
// "08/21/71",
// "8/1/71",
// "4/8/2014 22:05",
// "04/08/2014 22:05",
// "4/8/14 22:05",
// "04/2/2014 03:00:51",
// "8/8/1965 12:00:00 AM",
// "8/8/1965 01:00:01 PM",
// "8/8/1965 01:00 PM",
// "8/8/1965 1:00 PM",
// "8/8/1965 12:00 AM",
// "4/02/2014 03:00:51",
// "03/19/2012 10:11:59",
// "03/19/2012 10:11:59.3186369",
// // yyyy/mm/dd
// "2014/3/31",
// "2014/03/31",
// "2014/4/8 22:05",
// "2014/04/08 22:05",
// "2014/04/2 03:00:51",
// "2014/4/02 03:00:51",
// "2012/03/19 10:11:59",
// "2012/03/19 10:11:59.3186369",
// // Chinese
// "2014年04月08日",
// //   yyyy-mm-ddThh
// "2006-01-02T15:04:05+0000",
// "2009-08-12T22:15:09-07:00",
// "2009-08-12T22:15:09",
// "2009-08-12T22:15:09Z",
// //   yyyy-mm-dd hh:mm:ss
// "2014-04-26 17:24:37.3186369",
// "2012-08-03 18:31:59.257000000",
// "2014-04-26 17:24:37.123",
// "2013-04-01 22:43",
// "2013-04-01 22:43:22",
// "2014-12-16 06:20:00 UTC",
// "2014-12-16 06:20:00 GMT",
// "2014-04-26 05:24:37 PM",
// "2014-04-26 13:13:43 +0800",
// "2014-04-26 13:13:43 +0800 +08",
// "2014-04-26 13:13:44 +09:00",
// "2012-08-03 18:31:59.257000000 +0000 UTC",
// "2015-09-30 18:48:56.35272715 +0000 UTC",
// "2015-02-18 00:12:00 +0000 GMT",
// "2015-02-18 00:12:00 +0000 UTC",
// "2015-02-08 03:02:00 +0300 MSK m=+0.000000001",
// "2015-02-08 03:02:00.001 +0300 MSK m=+0.000000001",
// "2017-07-19 03:21:51+00:00",
// "2014-04-26",
// "2014-04",
// "2014",
// "2014-05-11 08:20:13,787",
// // mm.dd.yy
// "3.31.2014",
// "03.31.2014",
// "08.21.71",
// "2014.03",
// "2014.03.30",
// //  yyyymmdd and similar
// "20140601",
// "20140722105203",
// // unix seconds, ms, micro, nano
// "1332151919",
// "1384216367189",
// "1384216367111222",
// "1384216367111222333",
type DatesFilter struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (dt DatesFilter) Validate() error {
	zero := time.Time{}
	if dt.From.Equal(zero) {
		return errors.New("'from' field is missing or not parsed correctly")
	}

	if dt.To.Equal(zero) {
		return errors.New("'to' field is missing or not parsed correctly")
	}

	if dt.To.Before(dt.From) {
		return errors.New("'to' should be after 'from'")
	}

	return nil
}

func (filter DatesFilter) Where(column string) qm.QueryMod {
	return qm.Where(column+" BETWEEN ? AND ?", filter.From, filter.To)
}

func (filter DatesFilter) AndIn(column string) qm.QueryMod {
	return qm.And(column+" BETWEEN ? AND ?", filter.From, filter.To)
}

func (filter DatesFilter) Or(column string) qm.QueryMod {
	return qm.Or(column+" BETWEEN ? AND ?", filter.From, filter.To)
}

func (df *DatesFilter) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	df.From = time.Now()
	df.To = time.Now()

	for k, v := range rawStrings {
		if strings.ToLower(k) == "from" {
			df.From, err = dateparse.ParseAny(v)
			if err != nil {
				return err
			}
		}

		if strings.ToLower(k) == "to" {
			df.To, err = dateparse.ParseAny(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (df1 *DatesFilter) Equal(df2 *DatesFilter) bool {
	if df1 == nil {
		return df2 == nil
	}

	return df1.From.Equal(df2.From) && df1.To.Equal(df2.To)
}
