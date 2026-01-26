package convert

import (
	"log"
	"strings"
	"time"
)

// GetCurrentDateTime получить текущую дату и время в формате "2006-01-02 15:04:05".
func GetCurrentDateTime() string {
	return time.Now().Format(time.DateTime)
}

// FormatDateTime - функция форматирования даты в формате ДД.ММ.ГГГГ ЧЧ:ММ
func FormatDateTime(dateTime string) string {
	t, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return dateTime
	}

	return t.Format("02.01.2006 15:04:05")
}

// FormatDateTimeLocal преобразует дату в формат для datetime-local
func FormatDateTimeLocal(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Пробуем разные форматы из БД
	layouts := []string{
		"2006-01-02T15:04:05Z", // 2024-01-15T14:30:00Z
		"2006-01-02 15:04:05",  // 2024-01-15 14:30:00
		"2006-01-02T15:04:05",  // 2024-01-15T14:30:00
		"2006-01-02",           // 2024-01-15
		"02.01.2006",           // 15.01.2024
		"02.01.2006 15:04:05",  // 15.01.2024 14:30:00
		time.RFC3339,           // 2024-01-15T14:30:00+03:00
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, dateStr)
		if err == nil {
			return t.Format("2006-01-02T15:04")
		}
	}

	// Если не удалось распарсить, возвращаем как есть
	return dateStr
}

func ParseToMSSQLDateTime(goDateTime string) (string, error) {
	if goDateTime == "" {
		return "", nil
	}

	goDateTime = strings.TrimSpace(goDateTime)

	goDateTime = strings.Replace(goDateTime, "T", " ", 1)

	var t time.Time
	var err error

	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
	}

	for _, layout := range layouts {
		t, err = time.Parse(layout, goDateTime)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s]", err.Error(), pathToStrCode(), goDateTime)
	}

	return t.Format("20060102 15:04:05.000"), nil
}

func FormatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "не указано"
	}
	t := time.Unix(timestamp, 0)
	return t.Format("02.01.2006 15:04:05")
}
