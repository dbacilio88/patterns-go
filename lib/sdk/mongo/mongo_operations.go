package sdkmongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

/**
 * mongo_operations
 * <p>
 * This file contains core data structures and logic used throughout the application.
 *
 * <p><strong>Copyright © 2025 – All rights reserved.</strong></p>
 *
 * <p>This source code is distributed under a collaborative license.</p>
 *
 * <p>
 * Contributions, suggestions, and improvements are welcome!
 * You are free to fork, modify, and submit pull requests under the terms of the repository's license.
 * Please ensure proper attribution to the original author(s) and preserve this notice in derivative works.
 * </p>
 *
 * @author Christian Bacilio De La Cruz
 * @email dbacilio88@outlook.es
 * @since 7/11/2025
 */

const DefaultTimeout = 5 * time.Second

func InsertOne(context context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := getContext(context)
	defer cancel()
	return collection.InsertOne(ctx, document)
}

func FindOne(context context.Context, collection *mongo.Collection, filter bson.M, result interface{}) error {
	ctx, cancel := getContext(context)
	defer cancel()
	return collection.FindOne(ctx, filter).Decode(result)
}

func UpdateOne(context context.Context, collection *mongo.Collection, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := getContext(context)
	defer cancel()
	return collection.UpdateOne(ctx, filter, update)
}

func DeleteOne(context context.Context, collection *mongo.Collection, filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := getContext(context)
	defer cancel()
	return collection.DeleteOne(ctx, filter)
}

func FindAll(context context.Context, collection *mongo.Collection, filter bson.M) (*mongo.Cursor, error) {
	ctx, cancel := getContext(context)
	defer cancel()
	return collection.Find(ctx, filter)
}

func SaveBatch(context context.Context, collection *mongo.Collection, documents []any) (*mongo.BulkWriteResult, error) {

	if len(documents) == 0 {
		return nil, errors.New("documents is empty")
	}

	ctx, cancel := getContext(context)
	defer cancel()

	var models []mongo.WriteModel

	for _, document := range documents {
		models = append(models, mongo.NewInsertOneModel().SetDocument(document))
	}

	const batchSize = 100

	bulkWriteResult := &mongo.BulkWriteResult{
		UpsertedIDs: make(map[int64]interface{}),
	}

	for i := 0; i < len(models); i += batchSize {
		end := i + batchSize
		if end > len(models) {
			end = len(models)
		}
		batch := models[i:end]

		writed, err := collection.BulkWrite(ctx, batch)

		if err != nil {
			return nil, err
		}

		if bulkWriteResult == nil {
			bulkWriteResult = writed
		} else {
			bulkWriteResult.InsertedCount += writed.InsertedCount
			bulkWriteResult.MatchedCount += writed.MatchedCount
			bulkWriteResult.ModifiedCount += writed.ModifiedCount
			bulkWriteResult.DeletedCount += writed.DeletedCount
			bulkWriteResult.UpsertedCount += writed.UpsertedCount
			for k, v := range writed.UpsertedIDs {
				bulkWriteResult.UpsertedIDs[k] = v
			}
		}
	}
	return bulkWriteResult, nil
}

func getContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return context.WithTimeout(context.Background(), DefaultTimeout)
	}
	return context.WithTimeout(ctx, DefaultTimeout)
}
