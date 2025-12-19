package repository

// СОТРУДНИКИ
const (
	FGWsvPerformerAuthQuery     = "exec dbo.svPerformerAuth ?, ?;"  // ХП проверяет сотрудника по табельному номеру и паролю для авторизации.
	FGWsvPerformerFindByIdQuery = "exec dbo.svPerformerFindById ?;" // ХП ищет информацию о сотруднике по ИД.

)

// РОЛИ
const (
	FGWsvRoleFindByIdQuery = "exec dbo.svRoleFindById ?;" // ХП ищет роль.
)

// ПРОДУКЦИЯ
const (
	FGWsvAFormsProductionAllQuery = "exec dbo.svAFormsProductionAll;" //ХП выводит список продукции.
)
