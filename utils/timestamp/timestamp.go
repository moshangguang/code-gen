package timestamp

import (
	"sync"
	"sync/atomic"
	"time"
)

// year, month, day, hour, min, sec, msec
const (
	MilliSecond Duration = 1
	Second               = 1000 * MilliSecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
	Day                  = 24 * Hour
)

// 当前时间偏移量，单位为ms
var _offset int64 = 0

// 默认时区，北京时间
var _defLocation Timezone = 8
var _defLocationOnce sync.Once

// 设置当前时间的偏移量，单位为ms
func SetNowOffset(offset Duration) {
	atomic.StoreInt64(&_offset, int64(offset))

}

// 设置默认的Location
// 只可调用一次
func SetDefaultLocation(location Timezone) {
	_defLocationOnce.Do(func() {
		_defLocation = location
	})
}

// 获取当前时间
func Now() Time {
	return Time{
		num:      Timestamp(nanoToMilli(time.Now().UTC().UnixNano()) + _offset),
		timeZone: _defLocation,
	}
}

// 根据日期获取时间；t.num存储UTC时间
func Date(year, month, day, hour, min, sec, msec int, timezone Timezone) Time {
	zone := time.FixedZone("UTC", int(timezone)*60*60)
	goTime := time.Date(year, time.Month(month), day, hour, min,
		sec, msec*1e6, zone)
	return Time{
		num:      Timestamp(nanoToMilli(goTime.UnixNano())),
		timeZone: timezone,
	}
}

// 计算时间差，精确到ms
func Since(t Time) Duration {
	return Duration(Now().num - t.num)
}

// 转化为go原生的时间
func ToGoTime(t Time) time.Time {
	return time.Unix(0, t.timestampByNano())
}

// go原生的时间转为Time
func FromGoTime(t time.Time) Time {
	return Time{
		num:      Timestamp(nanoToMilli(t.UnixNano())),
		timeZone: _defLocation,
	}
}

// 转化为go原生的时间戳
func ToGoTimestamp(t Timestamp) int64 {
	return milliToNano(int64(t))
}

// go原生的时间戳转为毫秒时间戳
func FromGoTimestamp(t int64) Timestamp {
	return Timestamp(nanoToMilli(t))
}

// 解析时间戳，采用默认时区
func Parse(timestamp Timestamp) Time {
	return Time{
		timeZone: _defLocation,
		num:      timestamp,
	}
}

func ParseWithZone(timestamp Timestamp, timezone Timezone) Time {
	return Time{
		timeZone: timezone,
		num:      timestamp,
	}
}

// 精确到ms的时间戳
type Timestamp int64

// 时区
type Timezone int

// 时长
type Duration int64

func (d Duration) ToGoDuration() time.Duration {
	return time.Duration(d.Milliseconds() * 1e6)
}

func (d Duration) Milliseconds() int64 {
	return int64(d)
}

func (d Duration) Second() float64 {
	sec := d / Second
	nsec := d % Second
	return float64(sec) + float64(nsec)/1e3
}

func (d Duration) Minute() float64 {
	min := d / Minute
	nsec := d % Minute
	return float64(min) + float64(nsec)/(60*1e3)
}

func (d Duration) Hour() float64 {
	hour := d / Hour
	nsec := d % Hour
	return float64(hour) + float64(nsec)/(60*60*1e3)
}

type Time struct {
	num      Timestamp
	timeZone Timezone
}

// 纳秒转毫秒
func nanoToMilli(t int64) int64 {
	return t / 1e6
}

// 毫秒转纳秒
func milliToNano(t int64) int64 {
	return t * 1e6
}

// 返回 偏移量 单位：毫秒
func (t Time) timestamp() int64 {
	return int64(t.num)
}

// 返回 偏移量  单位：纳秒
func (t Time) timestampByNano() int64 {
	return t.timestamp() * 1e6
}

// 返回 时区单位：毫秒
func (t Time) timestampWithTimezone() int64 {
	return int64(t.num) + int64(t.timeZone)*int64(Hour)
}

// 返回 时区单位：纳秒
func (t Time) timestampWithTimezoneByNano() int64 {
	return t.timestampWithTimezone() * 1e6
}

// 返回算上时区和偏移量的时间戳
func (t Time) TimeStamp() Timestamp {
	return Timestamp(t.timestamp())
}

func (t Time) Zone() Timezone {
	return t.timeZone
}

// 设置时区
func (t Time) In(z Timezone) Time {
	t.timeZone = z
	return t
}

func (t Time) Year() int {
	return time.Unix(0, t.timestampByNano()).Year()
}

func (t Time) Month() int {
	return int(time.Unix(0, t.timestampByNano()).Month())
}

// 返回这是一周中的第几天
func (t Time) WeekDay() int {
	return int(time.Unix(0, t.timestampByNano()).Weekday())
}

// Day返回一个月中的第几天
func (t Time) Day() int {
	return time.Unix(0, t.timestampByNano()).Day()
}

// YearDay返回一年中的第几天
func (t Time) YearDay() int {
	return time.Unix(0, t.timestampByNano()).YearDay()
}

func (t Time) Hour() int {
	return int(Duration(t.timestampWithTimezone()) % Day / Hour)
}

func (t Time) Minute() int {
	return int(Duration(t.timestampWithTimezone()) % Hour / Minute)
}

func (t Time) Second() int {
	return int(Duration(t.timestampWithTimezone()) % Minute / Second)
}

func (t Time) MilliSecond() int {
	return int(Duration(t.timestampWithTimezone()) % Second)
}

func (t Time) Add(duration Duration) Time {
	t.num = t.num + Timestamp(duration)
	return t
}

func (t Time) AddYear(year int) Time {
	t.num = Timestamp(nanoToMilli(ToGoTime(t).AddDate(year, 0, 0).UnixNano()))
	return t
}

func (t Time) AddMonth(mon int) Time {
	t.num = Timestamp(nanoToMilli(ToGoTime(t).AddDate(0, mon, 0).UnixNano()))
	return t
}

func (t Time) AddWeek(week int) Time {
	t.num = Timestamp(nanoToMilli(ToGoTime(t).AddDate(0, 0, 7*week).UnixNano()))
	return t
}

func (t Time) AddDay(day int) Time {
	t.num = t.num + Timestamp(Duration(day)*Day)
	return t
}

func (t Time) AddHour(hour int) Time {
	t.num = t.num + Timestamp(Duration(hour)*Hour)
	return t
}

func (t Time) AddMinute(min int) Time {
	t.num = t.num + Timestamp(Duration(min)*Minute)
	return t
}

func (t Time) AddSecond(sec int) Time {
	t.num = t.num + Timestamp(Duration(sec)*Second)
	return t
}

func (t Time) AddMilliSecond(ms int) Time {
	t.num = t.num + Timestamp(ms)
	return t
}

// 当日的开始
func (t Time) BeginOfDay() Time {
	return Time{
		num:      t.num - Timestamp(t.timestampWithTimezone())%Timestamp(Day),
		timeZone: t.timeZone,
	}
}

// 当前小时时的开始
func (t Time) BeginOfHour() Time {
	return Time{
		num:      t.num - t.num%Timestamp(Hour),
		timeZone: t.timeZone,
	}
}

// 当前分钟的开始
func (t Time) BeginOfMinute() Time {
	minutePassedSec := t.num % Timestamp(Minute)
	return Time{
		num:      t.num - minutePassedSec,
		timeZone: t.timeZone,
	}
}

// 当日的结束
func (t Time) EndOfDay() Time {
	return Time{
		num:      t.BeginOfDay().num + Timestamp(Day),
		timeZone: t.timeZone,
	}
}

// 当前小时时的结束
func (t Time) EndOfHour() Time {
	return Time{
		num:      t.BeginOfHour().num + Timestamp(Hour),
		timeZone: t.timeZone,
	}
}

// 当前分钟的结束
func (t Time) EndOfMinute() Time {
	return Time{
		num:      t.BeginOfMinute().num + Timestamp(Minute),
		timeZone: t.timeZone,
	}
}

// 是否在u Time之前
func (t Time) Before(u Time) bool {
	return t.num < u.num
}

// 是否在u Time之后
func (t Time) After(u Time) bool {
	return t.num > u.num
}

// 减去u Time之后的时间，返回的时间单位为ms
func (t Time) Sub(u Time) Duration {
	return Duration(t.num - u.num)
}

func IsSameDay(now, last Time) bool {
	return now.YearDay() == last.YearDay() && now.Year() == last.Year()
}

func IsPassDay(now, last Time) bool {
	if now.Year() > last.Year() {
		return true
	}
	return now.YearDay() > last.YearDay() && now.Year() >= last.Year()
}

// 输出"Mon Jan _2 15:04:05 MST 2006"标准的时间
func (t Time) Format() string {
	return ToGoTime(t).Format(time.UnixDate)
}

//将时间字符串转成 Time
func StringToTime(t string) Time {
	timeTemplate := "2006-01-02 15:04:05"
	//解决与当前系统时间相差8小时
	loc := time.FixedZone("CST", 3600*int(_defLocation))
	theTime, err := time.ParseInLocation(timeTemplate, t, loc)
	Must(err)
	return FromGoTime(theTime)
}
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
