package mongo

import (
	"context"
	"log"

	"github.com/Krylphi/rockspoon-cart/internal/domain"
	"github.com/Krylphi/rockspoon-cart/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// Config is used for repository configuration parameters
	Config struct {
		ConnStr         string
		Database        string
		UsersCollection string
	}

	mongoRepository struct {
		collection *mongo.Collection
		config     *Config
	}

	cartModelWrapper struct {
		ObjID primitive.ObjectID `bson:"_id" json:"_id"`
		*domain.Cart
	}
)

// NewMongoRepository is a constructor for Repository instance with Mongo implementation
func NewMongoRepository(c *Config) (repository.CartRepository, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(c.ConnStr)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print(err)

		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print(err)

		return nil, err
	}

	log.Print("Connected to MongoDB!")

	return &mongoRepository{
		collection: client.Database(c.Database).Collection(c.UsersCollection),
		config:     c,
	}, nil
}

func (repo *mongoRepository) CreateCart(ctx context.Context) (*domain.Cart, error) {
	id := primitive.NewObjectID()
	sid := id.Hex()
	cart := &domain.Cart{ID: sid, Items: make([]*domain.CartItem, 0)}

	cartModel := cartModelWrapper{
		ObjID: id,
		Cart:  cart,
	}

	_, err := repo.collection.InsertOne(ctx, cartModel)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Println("Inserted a single Cart document: ", sid)

	return cart, nil
}

func (repo *mongoRepository) AddItem(ctx context.Context, cartID string, product string, quantity int) (*domain.CartItem, error) {
	objID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	item := &domain.CartItem{
		ID:       primitive.NewObjectID().Hex(),
		CartID:   cartID,
		Product:  product,
		Quantity: quantity,
	}

	if err := item.Validate(); err != nil {
		log.Print(err)
		return nil, err
	}

	update := bson.M{
		"$push": bson.M{"cart.items": item},
	}

	updateRes, err := repo.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if updateRes.MatchedCount == 0 {
		return nil, domain.ErrCartNotFound
	}

	if updateRes.ModifiedCount == 0 {
		return nil, domain.ErrCartNotFound
	}

	log.Println("Updated a single Cart data document: ", cartID)

	return item, nil
}

func (repo *mongoRepository) RemoveItem(ctx context.Context, cartID string, itemID string) error {
	objID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		log.Print(err)
		return err
	}

	update := bson.M{
		"$pull": bson.M{"cart.items": bson.M{"id": itemID}},
	}

	updateRes, err := repo.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Print(err)
		return err
	}

	if updateRes.MatchedCount == 0 {
		return domain.ErrCartNotFound
	}

	if updateRes.ModifiedCount == 0 {
		return domain.ErrNoSuchCartItem
	}

	log.Println("Updated a single User data document: ", cartID)

	return nil
}

func (repo *mongoRepository) DeleteCart(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Print(err)
		return err
	}

	singleResult := repo.collection.FindOneAndDelete(ctx, bson.M{"_id": objID})
	if singleResult != nil {
		if singleResult.Err() != nil {
			log.Println(singleResult.Err())
			return singleResult.Err()
		}
	}

	return nil
}

func (repo *mongoRepository) Cart(ctx context.Context, id string) (*domain.Cart, error) {
	var result cartModelWrapper

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = repo.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result.Cart, nil
}
