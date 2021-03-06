package wiki

import (
	"log"
	"strconv"
	"strings"
	"time"
)

var monthsGenitive = [...]string{
	"января",
	"февраля",
	"марта",
	"апреля",
	"мая",
	"июня",
	"июля",
	"августа",
	"сентября",
	"октября",
	"ноября",
	"декабря",
}

var weekDays = [...]string{
	"воскресенье",
	"понедельник",
	"вторник",
	"среда",
	"четверг",
	"пятница",
	"суббота",
}

const holidaysHeader = "Праздники и памятные дни"
const intHolidaysSubheader = "Международные"
const locHolidaysSubheader = "Национальные"
const rlgHolidaysSubheader = "Религиозные"
const profHolidaysSubheader = "Профессиональные"
const nameDaysSubheader = "Именины"

const moscowLocation = "Europe/Moscow"

type Report struct {
	Stats        string
	Common       []string
	HolidaysInt  []string
	HolidaysLoc  []string
	HolidaysProf []string
	HolidaysRlg  ReligiousHolidays
	NameDays     []string
	Omens        []string
}

type ReligiousHolidayDescr struct {
	Descriptions []string
	GroupAbbr    string
}

type ReligiousHolidays struct {
	Holidays []*ReligiousHolidayDescr
}

func (holidays *ReligiousHolidays) Empty() bool {
	empty := true
	for _, item := range holidays.Holidays {
		if len(item.Descriptions) > 0 {
			empty = false
		}
	}
	return empty
}

func (holidays *ReligiousHolidays) AppendString(formatted *string) {
	if len(holidays.Holidays) > 0 {
		for _, items := range holidays.Holidays {
			for _, line := range items.Descriptions {
				*formatted += "- " + line
				if items.GroupAbbr != "" {
					*formatted += " (" + items.GroupAbbr + ")"
				}
				*formatted += "\n"
			}
		}
	}
}

func (report *Report) String() string {
	formattedStr := ""
	if report.Stats != "" {
		formattedStr += report.Stats + "\n"
	}

	if len(report.HolidaysInt) > 0 || len(report.HolidaysLoc) > 0 || len(report.HolidaysProf) > 0 || !report.HolidaysRlg.Empty() {
		formattedStr += "*" + holidaysHeader + "*\n"
		if len(report.HolidaysInt) > 0 {
			formattedStr += "\n_" + intHolidaysSubheader + "_\n"
			for _, line := range report.HolidaysInt {
				formattedStr += "- " + line + "\n"
			}
		}
		if len(report.HolidaysLoc) > 0 {
			formattedStr += "\n_" + locHolidaysSubheader + "_\n"
			for _, line := range report.HolidaysLoc {
				formattedStr += "- " + line + "\n"
			}
		}
		if len(report.HolidaysProf) > 0 {
			formattedStr += "\n_" + profHolidaysSubheader + "_\n"
			for _, line := range report.HolidaysProf {
				formattedStr += "- " + line + "\n"
			}
		}
		if !report.HolidaysRlg.Empty() {
			formattedStr += "\n_" + rlgHolidaysSubheader + "_\n"
			report.HolidaysRlg.AppendString(&formattedStr)
		}
	}

	if len(report.NameDays) > 0 {
		formattedStr += "\n_" + nameDaysSubheader + "_"
		needAppend := false
		for _, line := range report.NameDays {
			if strings.Contains(line, ":") {
				formattedStr += "\n- " + line
				needAppend = false
			} else {
				if needAppend {
					formattedStr += ", " + line
				} else {
					formattedStr += "\n- " + line
					needAppend = true
				}
			}
		}
		formattedStr += "\n"
	}

	if l := len(report.Omens); l > 0 {
		formattedStr += "\n*" + "Приметы" + "*\n\n"
		for i, line := range report.Omens {
			if i > 0 && i < 5 {
				formattedStr += line + "\n"
			} else if i == 0 {
				formattedStr += "_" + line + "_\n"
			} else {
				break
			}
		}
	}
	return formattedStr
}

func (report *Report) setCalendarInfo(day *time.Time) {
	report.Stats = GenerateCalendarStats(day)
}

type Response struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         Query  `json:"query"`
}

type Query struct {
	Pages map[string]Pages `json:"pages"`
}

type Pages struct {
	Title   string `json:"title"`
	Extract string `json:"extract"`
	PageId  uint64 `json:"pageid"`
	NS      uint64 `json:"ns"`
}

func GetTodaysReport(holidays *Holidays) string {
	location, _ := time.LoadLocation(moscowLocation)
	log.Print(location)
	today := time.Now().In(location)

	return GetReport(holidays, &today)
}

func GetReport(holidays *Holidays, report_date *time.Time) string {
	var report = ExtractReport(holidays, report_date.Month(), report_date.Day())
	report.setCalendarInfo(report_date)
	return report.String()
}

func ExtractReport(holidays *Holidays, month time.Month, day int) Report {
	log.Println("Extract info", month, day)
	var m = *(*holidays)[month]
	var d = *m[day]
	return d.Report
}

func getDateString(day *time.Time) string {
	_, month, dayNum := day.Date()
	return strconv.Itoa(dayNum) + " " + monthsGenitive[month-1]
}

func getFullDateString(day *time.Time) string {
	year, month, dayNum := day.Date()
	weekDay := strings.Title(getWeekDateString(day))
	return "*" + weekDay + ", " + strconv.Itoa(dayNum) + " " + monthsGenitive[month-1] + " " + strconv.Itoa(year) + " года" + "*"
}

func getWeekDateString(day *time.Time) string {
	weekday := int(day.Weekday())
	return weekDays[weekday]
}

func GetDayNoun(day int) string {
	rest := day % 10
	if (day > 10) && (day < 20) {
		// для второго десятка - всегда третья форма
		return "дней"
	} else if rest == 1 {
		return "день"
	} else if rest > 1 && rest < 5 {
		return "дня"
	} else {
		return "дней"
	}
}

func GenerateCalendarStats(reportDay *time.Time) string {
	firstLine := getFullDateString(reportDay)

	year := time.Date(reportDay.Year(), time.December, 31, 0, 0, 0, 0, time.UTC)
	infoDay := reportDay.YearDay()
	full_days := year.YearDay()

	rest := full_days - infoDay
	secondLine := ""
	if rest > 0 {
		secondLine = strconv.Itoa(infoDay) + "-й день года. До конца года " + strconv.Itoa(rest) + " " + GetDayNoun(rest)
	} else {
		secondLine = "Завтра уже Новый Год!"
	}

	return firstLine + "\n" + secondLine + "\n"
}
