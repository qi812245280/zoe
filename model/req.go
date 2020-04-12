package model

type OrgCreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Private bool   `json:"private"`
}

type OrgUpdateRequest struct {
	Private bool `json:"private"`
}

type AuthorizeOrgRequest struct {
	Type   string `json:"type"`
	UserId int    `json:"user_id"`
}

type CreateProjectRequest struct {
	ParentId int    `json:"parent_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Private  string `json:"private"`
}

type UpdateProjectRequest struct {
	Private string `json:"private"`
}
