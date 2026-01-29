package model

type Plan struct {
	IdPlan        int    `json:"idPlan"`        // IdPlan - ид плана
	PlanShift     int    `json:"planShift"`     // PlanShift - номер смены
	ExtProduction int    `json:"extProduction"` // ExtProduction - внешний ключ с таблицей svTB_Production.idProduction
	ExtSector     int    `json:"extSector"`     // ExtSector - внешний ключ с таблицей svTB_Sector.idSector
	PlanCount     int    `json:"planCount"`     // PlanCount - плановое кол-во продукции
	PlanDate      string `json:"planDate"`      // PlanDate - дата плана создания
	PlanInfo      string `json:"planInfo"`      // PlanInfo - комментарий к плану
	PlEditDate    string `json:"plEditDate"`    // PlEditDate - дата редактирования плана
	AuditRec      Audit  `json:"auditRec"`      // AuditRec - аудит
}
