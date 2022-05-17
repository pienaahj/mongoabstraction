package mongoabstraction

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName          string = "testdb"
	collectionName        string = "movies"
	connectionStringAdmin string = "mongodb://admin:myadminpassword@192.168.0.148:27017"
	connectionStringUser  string = "mongodb://user2:user2password@192.168.0.148:27017/user2?authSource=testdb"
)

var (
	DBError   error
	adminUser options.Credential = options.Credential{
		Username:   "admin",
		Password:   "myadminpassword",
		AuthSource: "admin",
	}
	user2User options.Credential = options.Credential{
		Username:   "user2",
		Password:   "user2password",
		AuthSource: "testdb",
	}
	adminURI string = "192.168.0.148:27017"
	user2URI string = "192.168.0.148:27017/user2"
)

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect(context.Context) error
	StartSession() (mongo.Session, error)
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoSession struct {
	mongo.Session
}

func NewClient() (ClientHelper, error) {
	c, err := mongo.NewClient(options.Client().SetAuth(
		options.Credential{
			Username:   adminUser.Username,
			Password:   adminUser.Password,
			AuthSource: adminUser.AuthSource,
		}).ApplyURI(adminURI))

	return &mongoClient{cl: c}, err

}

// func NewDatabase(cnf *config.Config, client ClientHelper) DatabaseHelper {
// 	return client.Database(cnf.DatabaseName)
// }

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.cl.Connect(ctx)
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}
