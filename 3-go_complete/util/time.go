package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	Loc *time.Location
)

var (
	dateFormatReg1  = regexp.MustCompile(`(?P<year>\d{4}).(?P<month>\d{1,2}).(?P<day>\d{1,2})`)
	dateFormatReg2  = regexp.MustCompile(`(?P<year>\d{4}).(?P<month>\d{1,2})`)
	dateFormatReg3  = regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{1,2})-(?P<day>\d{1,2})`)
	dateFormatReg4  = regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{1,2})`)
	dateFormatReg5  = regexp.MustCompile(`(?P<year>\d{4})年(?P<month>\d{1,2})月(?P<day>\d{1,2})日`)
	dateFormatReg6  = regexp.MustCompile(`(?P<year>\d{4})年(?P<month>\d{1,2})月`)
	dateFormatReg7  = regexp.MustCompile(`(?P<year>\d{4})年`)
	dateFormatReg8  = regexp.MustCompile(`(?P<year>\d{4})/(?P<month>\d{1,2})/(?P<day>\d{1,2})`)
	dateFormatReg9  = regexp.MustCompile(`(?P<year>\d{4})/(?P<month>\d{1,2})`)
	dateFormatReg10 = regexp.MustCompile(`(?P<year>\d{4})`)
)

func init() {
	Loc, _ = time.LoadLocation("Asia/Shanghai")
}

const (
	DATE_LAYOUT1     = "2006-01-02"
	DATE_LAYOUT2     = "20060102"
	DATETIME_LAYOUT1 = "2006-01-02 15:04:05"
	DATETIME_LAYOUT2 = "20060102150405"
	DATETIME_LAYOUT3 = "2006-01-02T15:04:05-07:00"
)

func IsValidTime(t time.Time) bool {
	return t.Year() > 1900
}

func ParseDate(text string) (time.Time, error) {
	year := 0
	month := 1
	day := 1
	matches := []string{}
	var fitReg *regexp.Regexp
	for _, reg := range []*regexp.Regexp{dateFormatReg1, dateFormatReg2, dateFormatReg3, dateFormatReg4, dateFormatReg5, dateFormatReg6, dateFormatReg7, dateFormatReg8, dateFormatReg9, dateFormatReg10} {
		matches = reg.FindStringSubmatch(text)
		if len(matches) > 0 {
			fitReg = reg
			break
		}
	}
	if len(matches) > 0 {
		year, _ = strconv.Atoi(matches[fitReg.SubexpIndex("year")])
		if fitReg.SubexpIndex("month") >= 0 {
			month, _ = strconv.Atoi(matches[fitReg.SubexpIndex("month")])
		}
		if fitReg.SubexpIndex("day") >= 0 {
			day, _ = strconv.Atoi(matches[fitReg.SubexpIndex("day")])
		}
	} else {
		return time.Time{}, fmt.Errorf("invalid date format %s", text)
	}

	rect, err := time.ParseInLocation(DATE_LAYOUT1, fmt.Sprintf("%4d-%02d-%02d", year, month, day), Loc)
	if err != nil {
		return time.Time{}, err
	} else {
		return rect, nil
	}
}

func GetMonthFirstDay(t time.Time) time.Time {
	currentYear, currentMonth, _ := t.Date()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, Loc)
}

// 自定义json里的时间格式，核心是自定义一个type，实现MarshalJSON和UnmarshalJSON这两个方法
type MyJsonDateTime struct {
	time.Time
	Format string //希望在json里以什么格式展示
}

func (d MyJsonDateTime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%s\"", time.Time(d.Time).Format(d.Format))
	return []byte(s), nil
}

// 要改变自己，必须传指针
func (d *MyJsonDateTime) UnmarshalJSON(bs []byte) (err error) {
	now, err := time.ParseInLocation(`"`+d.Format+`"`, string(bs), time.Local) //注意Format前后还得加引号
	*d = MyJsonDateTime{now, d.Format}                                         // 要改变自己
	return
}

// 在print(MyJsonDateTime)时会调用String()方法
func (d MyJsonDateTime) String() string {
	return time.Time(d.Time).Format(d.Format)
}
