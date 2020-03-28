package repositories

import (
	"database/sql"
	"fmt"
	"github.com/chebizarro/redshed/pkg/models"
	"strings"
)

type ProjectRepository struct {
	db *sql.DB
}


func (r *ProjectRepository) Init(db * sql.DB) {
	r.db = db
}

func (repo *ProjectRepository) CreateProject(project models.Project) (models.Project, error) {
	//tags := strings.Join(models.Tags, ";")
	statement := `
    insert into projects (title)
    values ($1)
    returning id
  `
	var id int64
	err := repo.db.QueryRow(statement, project.Title).Scan(&id)
	if err != nil {
		return project, err
	}
	createdProject := project
	createdProject.ID = id
	return createdProject, nil
}

// GetPendingProjects method
func (repo *ProjectRepository) GetPendingProjects(userid string) ([]models.Project, error) {
	query := `
    select id, title, tags
    from projects
    where completed = $2
  `
	rows, err := repo.db.Query(query, userid, false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getProjectsFromRows(rows)
}

// GetProjectsByDateRange method
func (repo *ProjectRepository) GetProjectsByDateRange(userid string, from string, to string) ([]models.Project, error) {
	query := `
    select id, title, tags
    from projects
  `
	rows, err := repo.db.Query(query, userid, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getProjectsFromRows(rows)
}

// GetProjectByID method
func (repo *ProjectRepository) GetProjectByID(userid string, id int64) (models.Project, error) {
	query := `
    select id, title, tags
    from projects
  `
	row := repo.db.QueryRow(query, userid, id)
	var tags string
	var task models.Project
	err := row.Scan(&task.ID, &task.Title, &tags)
	if err != nil {
		return models.Project{}, err
	}

	task.Tags = []string{}
	if len(tags) > 0 {
		task.Tags = strings.Split(tags, ";")
	}
	return task, nil
}

// UpdateProject method
func (repo *ProjectRepository) UpdateProject(userid string, id int64, task models.Project) error {
	query := `
    update projects
    set title=$3, tags=$7
    where id=$2
  `

	tags := strings.Join(task.Tags, ";")
	res, err := repo.db.Exec(query, id, task.Title, tags)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated for %d", id)
	}

	return nil
}

// DeleteProject method
func (repo *ProjectRepository) DeleteProject(userid string, id int64) error {
	query := `delete from projects where id=$2`
	res, err := repo.db.Exec(query, userid, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("exactly 1 row is not impacted for %d", id)
	}

	return nil
}

func getProjectsFromRows(rows *sql.Rows) ([]models.Project, error) {
	tasks := []models.Project{}
	for rows.Next() {
		var task models.Project
		var tags string
		err := rows.Scan(&task.ID, &task.Title)
		if err != nil {
			return nil, err
		}

		task.Tags = []string{}
		if len(tags) > 0 {
			task.Tags = strings.Split(tags, ";")
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
