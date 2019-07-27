package models

import (
	"fmt"
	"time"

	"github.com/pmoorani/booksAPI/database"
	u "github.com/pmoorani/booksAPI/utils"
	uuid "github.com/satori/go.uuid"
)

type (
	Task struct {
		BaseModel
		Title       string    `json:"title"`
		Description string    `json:"description"`
		PriorityID  uint      `json:"priority_id" gorm:"association_autocreate:false;default:1"`
		StatusID    uint      `json:"status_id" gorm:"association_autocreate:false;default:1"`
		Start       time.Time `json:"start"`
		End         time.Time `json:"end"`
		UserID      uuid.UUID `json:"user_id"`
	}

	TransformedTask struct {
		BaseModel
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Start       time.Time       `json:"start"`
		End         time.Time       `json:"end"`
		User        TransformedUser `json:"user"`
		Status      Status          `json:"status"`
		Priority    Priority        `json:"priority"`
	}

	Status struct {
		ID        uint       `json:"id" gorm:"primary_key"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
		Name      string     `json:"name"`
		NameDE    string     `json:"name_de"`
	}

	Priority struct {
		ID        uint       `json:"id" gorm:"primary_key"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
		Name      string     `json:"name"`
		NameDE    string     `json:"name_de"`
	}
)

func (p *Priority) Scan(v interface{}) error {
	var priority Priority
	err := database.DB.Where("id = ?", v).Find(&priority).Error

	if err != nil {
		return err
	}

	*p = priority
	return nil
}

func (s *Status) Scan(v interface{}) error {
	var status Status
	err := database.DB.Where("id = ?", v).Find(&status).Error

	if err != nil {
		return err
	}

	*s = status
	return nil
}

func (task *Task) Validate() (map[string]interface{}, bool) {
	return u.Message(true, "Requirement passed!"), true
}


func AllTasks() ([]TransformedTask, error) {
	var tasks []Task
	var _tasks []TransformedTask
	var status Status
	var priority Priority
	var user TransformedUser


	err := database.DB.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	for _, item := range tasks {
		if item.StatusID > 0 {
			database.DB.Where("id = ?", item.StatusID).Find(&status)
		} else {
			status = Status{}
		}

		if item.PriorityID > 0 {
			database.DB.Where("id = ?", item.PriorityID).Find(&priority)
		} else {
			priority = Priority{}
		}

		user, err = FindUserByID(item.UserID)
		if err != nil {
			fmt.Println(err)
			user = TransformedUser{}
		}

		_tasks = append(_tasks, TransformedTask{
			BaseModel:   item.BaseModel,
			Title:       item.Title,
			Description: item.Description,
			Start:       item.Start,
			End:         item.End,
			Status:      status,
			Priority:    priority,
			User:        user,
		})
	}

	return _tasks, nil
}

func FindTaskByID(id interface{}) (TransformedTask, error) {
	var task Task
	var status Status
	var priority Priority
	var user TransformedUser

	err := database.DB.Where("id = ?", id).Find(&task).Error

	if err != nil {
		return TransformedTask{}, err
	}

	if task.StatusID > 0 {
		database.DB.Where("id = ?", task.StatusID).Find(&status)
	} else {
		status = Status{}
	}

	if task.PriorityID > 0 {
		database.DB.Where("id = ?", task.PriorityID).Find(&priority)
	} else {
		priority = Priority{}
	}

	user, err = FindUserByID(task.UserID)
	if err != nil {
		fmt.Println(err)
		user = TransformedUser{}
	}
	_task := TransformedTask{
		BaseModel:   task.BaseModel,
		Title:       task.Title,
		Description: task.Description,
		Start:       task.Start,
		End:         task.End,
		User:        user,
		Status:      status,
		Priority:    priority,
	}
	return _task, nil
}

func DeleteTaskByID(id interface{}) error {
	var task Task
	err := database.DB.Where("id = ?", id).Find(&task).Error

	if err != nil {
		return err
	}

	if err = database.DB.Delete(&task).Error; err != nil {
		return err
	}
	return nil
}

func AllStatuses() ([]Status, error) {
	var statuses []Status
	err := database.DB.Find(&statuses).Error

	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func AllPriorities() ([]Priority, error) {
	var priorities []Priority
	err := database.DB.Find(&priorities).Error

	if err != nil {
		return nil, err
	}
	return priorities, nil
}
