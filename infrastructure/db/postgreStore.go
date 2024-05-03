package db

import (
	"kurs-server/domain/entities"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreStore struct {
	Logger *zap.Logger
	db     *gorm.DB
}

func NewPostgreStore(logger *zap.Logger, dsn string) *PostgreStore {
	store := &PostgreStore{
		Logger: logger,
	}

	store.Init(dsn)

	return store
}

func (s *PostgreStore) DB() *gorm.DB {
	return s.db
}

func (s *PostgreStore) Init(dsn string) {
	s.Logger.Debug("Connecting to the database...", zap.String("dsn: ", dsn))
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		s.Logger.Error("Failed to connect to the database", zap.String("Error: ", err.Error()))
		panic(err)
	}

	s.db = db

	err = s.db.AutoMigrate(
		&entities.Category{},
		&entities.Detail{},
		&entities.DetailCategory{},
		&entities.Product{},
		&entities.ProductCategory{},
		&entities.DetailValue{},
		&entities.User{},
		&entities.OrderGroup{},
		&entities.Order{},
		&entities.Review{},
		&entities.Star{},
	)
	if err != nil {
		s.Logger.Error("Failed to migrate the database", zap.String("Error: ", err.Error()))
		return
	}

	s.AutomateRating()
}

// AutomateRating TODO: Probably will NOT work because of wrong table names
func (s *PostgreStore) AutomateRating() {
	// Creating a view
	s.db.Exec(`CREATE OR REPLACE VIEW products_with_rating AS
        SELECT p.id, 
               COALESCE(SUM(s.amount * 1.0) / COUNT(s.id), 0) AS rating
        FROM products p
        LEFT JOIN stars s ON s.product_id = p.id
        GROUP BY p.id;`)

	// Creating a function to count ratings
	s.db.Exec(`CREATE OR REPLACE FUNCTION update_product_rating()
        RETURNS TRIGGER AS $$
        BEGIN
            UPDATE products p
            SET rating = (
                SELECT COALESCE(SUM(s.amount * 1.0) / COUNT(s.id), 0)
                FROM stars s
                WHERE s.product_id = p.id
            )
            WHERE p.id = NEW.product_id;
            RETURN NEW;
        END;
        $$ LANGUAGE plpgsql;`)

	// Creating triggers
	s.db.Exec(`CREATE TRIGGER update_product_rating_insert
        AFTER INSERT ON stars
        FOR EACH ROW
        EXECUTE PROCEDURE update_product_rating();`)

	s.db.Exec(`CREATE TRIGGER update_product_rating_update
        AFTER UPDATE ON stars
        FOR EACH ROW
        EXECUTE PROCEDURE update_product_rating();`)

	s.db.Exec(`CREATE TRIGGER update_product_rating_delete
        AFTER DELETE ON stars
        FOR EACH ROW
        EXECUTE PROCEDURE update_product_rating();`)
}
