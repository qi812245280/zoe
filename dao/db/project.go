package db

import (
	"database/sql"
	"errors"
	"github.com/cihub/seelog"
	"zoe/model"
)

func queryProject(conn *sql.Tx, sql string, args ...interface{}) (*[]model.Project, error) {
	var projects []model.Project
	rows, err := conn.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	var project model.Project
	for rows.Next() {
		err = rows.Scan(&project.Id, &project.Name, &project.ParentId, &project.Visibility, &project.CurrentVersionId,
			&project.IsDeleted, &project.UpdatedAt, &project.CreateAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return &projects, nil
}

func GetProjectById(conn *sql.Tx, id int) (*model.Project, error) {
	sql := "select * from project where id = ? and is_deleted = 0"
	projects, err := queryProject(conn, sql, id)
	if err != nil {
		return nil, err
	}
	if len(*projects) > 0 {
		return &(*projects)[0], nil
	}
	return nil, nil
}

func GetProjectByParentIdAndName(conn *sql.Tx, parentId int, name string) (*model.Project, error) {
	sql := "select * from project where name = ? and parent_id = ? and is_deleted = 0"
	projects, err := queryProject(conn, sql, name, parentId)
	if err != nil {
		return nil, err
	}
	if len(*projects) > 0 {
		return &(*projects)[0], nil
	}
	return nil, nil
}

func ListProject(conn *sql.Tx, orgId int) (*[]model.Project, error) {
	sql := "select * from project where parent_id = ? and is_deleted = 0"
	projects, err := queryProject(conn, sql, orgId)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func ListProjectByVisibility(conn *sql.Tx, orgId int) (*[]model.Project, *[]model.Project, error) {
	projects, err := ListProject(conn, orgId)
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
			_ = seelog.Critical("unKown project")
		}
	}
	return &publicProjects, &privateProjects, nil
}

func CreateProject(conn *sql.Tx, name string, visibility int, parentId int) (int, error) {
	project, err := GetProjectByParentIdAndName(conn, parentId, name)
	if err != nil {
		return 0, err
	}
	if project != nil {
		return 0, errors.New("该project已经存在")
	}
	sql := "insert into project (name, parent_id, visibility) values(?, ?, ?)"
	r, err := conn.Exec(sql, name, parentId, visibility)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdateProject(conn *sql.Tx, projectId int, private string) error {
	var visibility int
	if private == "true" {
		visibility = 0
	} else if private == "false" {
		visibility = 1
	}
	sql := "update project set visibility = ? where id = ? and is_deleted = 0"
	_, err := conn.Exec(sql, visibility, projectId)
	if err != nil {
		return err
	}
	return nil
}

func ListProjectByParentId(conn *sql.Tx, orgId int) (*[]model.Project, error) {
	sql := "select * from project where parent_id = ? and is_deleted = 0"
	projects, err := queryProject(conn, sql, orgId)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
