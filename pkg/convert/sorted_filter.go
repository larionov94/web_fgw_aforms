package convert

import "net/url"

func BuildSortURL(sortField, sortOrder, startDate, endDate, newSortField string) string {
	params := url.Values{}

	// Устанавливаем поле сортировки
	params.Set("sort", newSortField)

	// Устанавливаем порядок сортировки
	if sortField == newSortField {
		if sortOrder == "ASC" {
			params.Set("order", "DESC")
		} else {
			params.Set("order", "ASC")
		}
	} else {
		params.Set("order", "ASC")
	}

	// Добавляем даты, если они есть
	if startDate != "" {
		params.Set("startDate", startDate)
	}

	if endDate != "" {
		params.Set("endDate", endDate)
	}

	return "/aforms/plans?" + params.Encode()
}
