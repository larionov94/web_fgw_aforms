package model

type RoleList struct {
	Roles []*Role `json:"roles"`
}

// Role роль.
type Role struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	AuditRec Audit  `json:"auditRec"`
}
