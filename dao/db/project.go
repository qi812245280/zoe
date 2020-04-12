package db

import (
	"errors"
	"github.com/cihub/seelog"
	"zoe/model"
)

func GetProjectById(id int) (*model.Project, error) {
	var projects []model.Project
	sql := "select * from project where id = ? and is_deleted = 0"
	err := DB.Select(&projects, sql, id)
	if err != nil {
		return nil, err
	}
	if len(projects) > 0 {
		return &projects[0], nil
	}
	return nil, nil
}

func GetProjectByParentIdAndName(parentId int, name string) (*model.Project, error) {
	var projects []model.Project
	sql := "select * from project where name = ? and parent_id = ? and is_deleted = 0"
	err := DB.Select(&projects, sql, name, parentId)
	if err != nil {
		return nil, err
	}
	if len(projects) > 0 {
		return &projects[0], nil
	}
	return nil, nil
}

func ListProject(orgId int) (*[]model.Project, error) {
	var projects []model.Project
	sql := "select * from project where parent_id = ? and is_deleted = 0"
	err := DB.Select(&projects, sql, orgId)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}

func ListProjectByVisibility(orgId int) (*[]model.Project, *[]model.Project, error) {
	projects, err := ListProject(orgId)
	if err != nil {
		return nil, nil, err
	}
	var publicProjects []model.Project
	var privateProjects []model.Project
	for _, item := range *projects {
		if item.Visibility == 1 {
			publicProjects = append(publicProjects, item)
		} else if item.Visibility == 0 {
			privateProjects = append(privateProjects, item)
		} else {
			seelog.Critical("unKown project")
		}
	}
	return &publicProjects, &privateProjects, nil
}

func CreateProject(name, private string, parentId int) (*model.Project, error) {
	project, err := GetProjectByParentIdAndName(parentId, name)
	if err != nil {
		return nil, err
	}
	if project != nil {
		return nil, errors.New("该project已经存在")
	}
	var visibility int
	if private == "true" {
		visibility = 0
	} else if private == "false" {
		visibility = 1
	}
	sql := "insert into project (name, parent_id, visibility) values(?, ?, ?)"
	r, err := DB.Exec(sql, name, parentId, visibility)
	if err != nil {
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	project, err = GetProjectById(int(id))
	if err != nil {
		return nil, err
	}
	return project, nil
}

func UpdateProject(projectId int, private string) error {
	var visibility int
	if private == "true" {
		visibility = 0
	} else if private == "false" {
		visibility = 1
	}
	sql := "update project set visibility = ? where id = ? and is_deleted = 0"
	_, err := DB.Exec(sql, visibility, projectId)
	if err != nil {
		return err
	}
	return nil
}
