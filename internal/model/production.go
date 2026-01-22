package model

type Production struct {
	IdProduction   int     `json:"idProduction"`   // IdProduction - ид продукции.
	PrName         string  `json:"prName"`         // PrName - наименование варианта упаковки продукции для упаковщика.
	PrShortName    string  `json:"prShortName"`    // PrShortName - короткое наименование продукции для этикетки.
	PrPackName     string  `json:"prPackName"`     // PrPackName - вариант упаковки.
	PrType         string  `json:"prType"`         // PrType - декларированная или нет.
	PrArticle      string  `json:"prArticle"`      // PrArticle - артикул варианта упаковки.
	PrColor        string  `json:"prColor"`        // PrColor - цвет продукции.
	PrBarCode      string  `json:"prBarCode"`      // PrBarCode - бар-код.
	PrCount        int     `json:"prCount"`        // PrCount - количество продукции в ряду.
	PrRows         int     `json:"prRows"`         // PrRows - количество рядов.
	PrWeight       float64 `json:"prWeight"`       // PrWeight - вес п\п (кг).
	PrHWD          string  `json:"prHWD"`          // PrHWD - габариты (мм) 1000(высота)х1200(ширина)х1000(глубина).
	PrInfo         string  `json:"prInfo"`         // PrInfo - информация о продукции\комментарий.
	PrStatus       bool    `json:"prStatus"`       // PrStatus - статус продукции.
	PrEditDate     string  `json:"prEditDate"`     // PrEditDate - дата и время изменения записи.
	PrEditUser     int     `json:"prEditUser"`     // PrEditUser - роль сотрудника. По умолчанию 1 - администратор, 5 - оператор.
	PrPart         int     `json:"prPart"`         // PrPart - номер текущей партии, номер партии и дата указываются вручную и не будут изменяться автоматически с течением времени.
	PrPartLastDate string  `json:"prPartLastDate"` // PrPartLastDate - дата выпуска партии. Тут в ручную указывается дата для отображения на этикетке.
	PrPartAutoInc  int     `json:"prPartAutoInc"`  // PrPartAutoInc - нумерация партии и даты! Ручная(0), Автоматическая(1), С указанной даты(2).
	PrPartRealDate string  `json:"prPartRealDate"` // PrPartRealDate - дата продукции пока неизвестное поле.
	PrArchive      bool    `json:"prArchive"`      // PrArchive - архивная запись или нет.
	PrPerGodn      int     `json:"prPerGodn"`      // PrPerGodn - срок годности в месяцах.
	PrSAP          string  `json:"prSAP"`          // PrSAP - сап-код.
	PrProdType     bool    `json:"prProdType"`     // PrProdType - тип продукции пищевая\не пищевая.
	PrUmbrella     bool    `json:"prUmbrella"`     // PrUmbrella - беречь от влаги.
	PrPerfumery    bool    `json:"prPerfumery"`    // prPerfumery - парфбмерия.
	PrSun          bool    `json:"prSun"`          // PrSun - беречь от солнца.
	PrDecl         bool    `json:"prDecl"`         // PrDecl - декларирования или нет.
	PrParty        bool    `json:"prParty"`        // PrParty - партионная или нет.
	PrGL           int     `json:"prGL"`           // PrGL - петля Мёбиуса.
	PrVP           int     `json:"prVP"`           // PrVP - ванная печь.
	PrML           int     `json:"prML"`           // PrML - машинная линия на печи.
	AuditRec       Audit   `json:"auditRec"`       // AuditRec - аудит.
}
