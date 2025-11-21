package data

import (
	"context"
	"errors"
	"task_manager/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TaskService contains business logic and MongoDB-based data storage.
type TaskService struct {
	tasksColl    *mongo.Collection
	countersColl *mongo.Collection
}

// NewTaskService initializes the service with Mongo collections.
func NewTaskService(tasksColl, countersColl *mongo.Collection) *TaskService {
	return &TaskService{
		tasksColl:    tasksColl,
		countersColl: countersColl,
	}
}

// getContext returns a context with timeout for DB operations.
func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// getNextID uses a counters collection to atomically increment and return a sequence.
func (s *TaskService) getNextID() (int, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"_id": "task_id"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}

	err := s.countersColl.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Seq, nil
}

// GetAllTasks returns all tasks.
func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	ctx, cancel := getContext()
	defer cancel()

	cur, err := s.tasksColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasks []models.Task
	for cur.Next(ctx) {
		var t models.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskByID returns a task by ID.
func (s *TaskService) GetTaskByID(id int) (models.Task, bool, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"id": id}
	var task models.Task
	err := s.tasksColl.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Task{}, false, nil
		}
		return models.Task{}, false, err
	}
	return task, true, nil
}

// CreateTask creates and stores a new task.
func (s *TaskService) CreateTask(input models.TaskInput) (models.Task, error) {
	ctx, cancel := getContext()
	defer cancel()

	newID, err := s.getNextID()
	if err != nil {
		return models.Task{}, err
	}

	task := models.Task{
		ID:          newID,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Status:      input.Status,
	}

	_, err = s.tasksColl.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// UpdateTask updates an existing task.
func (s *TaskService) UpdateTask(id int, input models.TaskInput) (models.Task, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"title":       input.Title,
			"description": input.Description,
			"due_date":    input.DueDate,
			"status":      input.Status,
		},
	}

	res := s.tasksColl.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	var updated models.Task
	err := res.Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}

	return updated, nil
}

// DeleteTask deletes a task by ID.
func (s *TaskService) DeleteTask(id int) error {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"id": id}
	result, err := s.tasksColl.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
