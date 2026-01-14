package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
	"time"

	"gorm.io/gorm"
)

type projectTaskService struct {
	taskRepo       domain.ProjectTaskRepository
	assignmentRepo domain.ProjectAssignmentRepository
	projectRepo    domain.ProjectRepository
}

// NewProjectTaskService 创建项目任务服务实例
func NewProjectTaskService(
	taskRepo domain.ProjectTaskRepository,
	assignmentRepo domain.ProjectAssignmentRepository,
	projectRepo domain.ProjectRepository,
) domain.ProjectTaskService {
	return &projectTaskService{
		taskRepo:       taskRepo,
		assignmentRepo: assignmentRepo,
		projectRepo:    projectRepo,
	}
}

// AcceptTask 用户接受项目任务
func (s *projectTaskService) AcceptTask(projectID, userID uint) (*domain.ProjectTask, error) {
	// 1. 验证项目是否存在
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 2. 验证项目状态是否为招募中
	if project.Status != "recruiting" && project.Status != "in_progress" {
		return nil, errors.New("该项目当前不接受新任务")
	}

	// 3. 验证用户是否被指派到该项目
	_, err = s.assignmentRepo.FindByProjectAndUser(projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("您未被指派到该项目")
		}
		return nil, err
	}

	// 4. 检查用户是否已经接受过该任务
	existingTask, err := s.taskRepo.FindByProjectAndUser(projectID, userID)
	if err == nil && existingTask != nil {
		return nil, errors.New("您已接受过该任务")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 5. 创建任务记录
	task := &domain.ProjectTask{
		ProjectID:  projectID,
		UserID:     userID,
		Status:     "accepted",
		AcceptedAt: time.Now(),
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	// 6. 如果项目状态为"招募中"，更新为"进行中"
	if project.Status == "recruiting" {
		project.Status = "in_progress"
		if err := s.projectRepo.Update(project); err != nil {
			// 记录错误但不影响任务接受结果
			// 可以考虑添加日志记录
		}
	}

	// 加载关联的项目信息（使用更新后的状态）
	task.Project = *project

	return task, nil
}

// GetUserTasks 获取用户的所有任务
func (s *projectTaskService) GetUserTasks(userID uint) ([]domain.ProjectTask, error) {
	return s.taskRepo.FindByUserID(userID)
}

// GetUserAcceptedProjectIDs 获取用户已接受任务的项目ID列表
func (s *projectTaskService) GetUserAcceptedProjectIDs(userID uint) ([]uint, error) {
	tasks, err := s.taskRepo.FindAcceptedByUserID(userID)
	if err != nil {
		return nil, err
	}

	projectIDs := make([]uint, 0, len(tasks))
	for _, task := range tasks {
		projectIDs = append(projectIDs, task.ProjectID)
	}

	return projectIDs, nil
}
