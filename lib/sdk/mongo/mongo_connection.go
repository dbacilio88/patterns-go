package sdkmongo

import (
	"context"
	"crypto/tls"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
	"sync"
	"time"
)

/**
 * mongo_connection
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
 * @since 7/10/2025
 */

type DatabaseMongoDb struct {
	Client     *mongo.Client
	Collection string
	Database   string
}

type ParamsMongoDb struct {
	Context context.Context
	Uri     string
	Tls     bool
	Timeout int
}

var (
	connectionMongo *mongo.Client
	mongoOnce       sync.Once
	err             error
)

func NewConnectionClient(params *ParamsMongoDb) (*mongo.Client, error) {

	if params == nil {
		return nil, errors.New("params is nil")
	}

	ctx := params.Context
	if ctx == nil {
		ctx = context.Background()
	}

	uri := params.Uri

	if uri == "" {
		uri = os.Getenv("MONGODB_URI")
	}

	if uri == "" {
		return nil, errors.New("MONGODB_URI is nil")
	}

	mongoOnce.Do(func() {

		opt := options.Client().
			ApplyURI(uri).
			SetServerSelectionTimeout(time.Duration(params.Timeout) * time.Second).
			SetConnectTimeout(time.Duration(params.Timeout) * time.Second)

		if params.Tls {
			opt.SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			})
		}

		connectionMongo, err = mongo.Connect(opt)

		if err != nil {
			connectionMongo = nil
		}

		err = connectionMongo.Ping(context.TODO(), nil)
	})

	if err == nil || connectionMongo == nil {
		return nil, errors.New("failed to connect to MongoDB")
	}

	return connectionMongo, nil
}

func GetCollection(param *DatabaseMongoDb) (*mongo.Collection, error) {

	filter := bson.M{
		"name": param.Collection,
	}

	collections, err := param.Client.Database(param.Database).ListCollectionNames(context.TODO(), filter)

	if err != nil {
		return nil, errors.New("failed to list collections")
	}

	for _, collection := range collections {
		if collection == param.Collection {
			return param.Client.Database(param.Database).Collection(param.Collection), nil
		}
	}
	return nil, errors.New("failed to find collection")
}

func CloseConnection(client *mongo.Client) {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			return
		}
	}
}
