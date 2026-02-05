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
	FGWsvAFormsProductionAllQuery      = "exec dbo.svAFormsProductionAll ?, ?;"                                                                      //ХП выводит список продукции.
	FGWsvAFormsProductionFilterQuery   = "exec dbo.svAFormsProductionFilterById ?, ?, ?;"                                                            // ХП ищет продукцию по артиклю и наименованию и коду продукции.
	FGWsvAFormsProductionAddQuery      = "exec dbo.svAFormsProductionAdd ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?;"    // ХП добавляет продукцию.
	FGWsvAFormsProductionUpdQuery      = "exec dbo.svAFormsProductionUpd ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?;" // ХП добавляет продукцию.
	FGWsvAFormsProductionFindByIdQuery = "exec dbo.svAFormsProductionFindById ?"
)

const (
	FGWsvAFormsDesignNameAllQuery = "exec dbo.svAFormsDesignNameAll;" // ХП возвращает список конструкторских наименование.
	FGWsvAFormsColorAllQuery      = "exec dbo.svAFormsColorAll;"      // ХП возвращает цвета.
)

const (
	FGWsvAFormsPlansWithSortingAndFilteringQuery = "exec dbo.svAFormsPlansWithSortingAndFiltering ?, ?, ?, ?, ?, ?;" //ХП возвращает список планов с сортировкой и с параметрами фильтрации.
)

const (
	FGWsvAFormsSectorAllQuery = "exec dbo.svAFormsSectorAll;" // ХП возвращает список печек
)
