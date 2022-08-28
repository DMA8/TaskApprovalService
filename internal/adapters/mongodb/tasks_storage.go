package mongodb

import (
	"context"
	"errors"

	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrTaskIsAlreadyExistsInStorage = errors.New("such task is already in storage")
)

//CRUD документов
//db.inventory.find({ $and: [ { price: 1.99 }, { qty: { $lt: 20 } }, { sale: true } ] } )
func (d *Database) CreateTask(ctx context.Context, task *models.Task) (primitive.ObjectID, error) {
	var dullTask models.Task
	ctx, span := otel.Tracer("team31_tasks").Start(ctx, "database createTask")
	defer span.End()
	if err := d.TasksCollection.FindOne(ctx, bson.M{"creator": task.Creator, "name": task.Name}).Decode(&dullTask); err == nil {
		return primitive.ObjectID{}, ErrTaskIsAlreadyExistsInStorage
	}
	insertResult, err := d.TasksCollection.InsertOne(ctx, task)
	span.SetAttributes(attribute.KeyValue{Key: "mongo_uuid", Value: attribute.StringValue(insertResult.InsertedID.(primitive.ObjectID).String())})
	if err != nil {
		return insertResult.InsertedID.(primitive.ObjectID), err
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func (d *Database) ReadTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	var outTask models.Task
	if err := d.TasksCollection.FindOne(ctx, bson.M{"creator": task.Creator, "name": task.Name}).Decode(&outTask); err != nil {
		if err != nil {
			return nil, err
		}
	}
	return &outTask, nil
}

func (d *Database) ReadTaskById(ctx context.Context, task *models.Task) (*models.Task, error) {
	var outTask models.Task
	if err := d.TasksCollection.FindOne(ctx, bson.M{"_id": task.ID}).Decode(&outTask); err != nil {
		if err != nil {
			return nil, err
		}
	}
	return &outTask, nil
}

//апдейтим только имя и описание таски
func (d *Database) UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	filter := bson.M{
		"creator": task.Creator,
		"name":    task.Name,
	}
	update := bson.M{
		"$set": task,
	}
	_, err := d.TasksCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return d.ReadTask(ctx, task)
}

func (d *Database) DeleteTask(ctx context.Context, task *models.Task) error {
	filter := createFilter(task.Creator, task.Name)
	_, err := d.TasksCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) ListTask(ctx context.Context, task *models.Task) ([]*models.Task, error) {
	var tasks []*models.Task
	filter := bson.M{"creator": task.Creator}
	cur, err := d.TasksCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var t models.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func createFilter(userLogin, taskName string) *primitive.M {
	filter := bson.M{
		"creator": userLogin,
		"name":    taskName,
	}
	return &filter
}
