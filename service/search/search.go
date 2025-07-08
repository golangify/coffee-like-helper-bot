package search

import (
	"coffee-like-helper-bot/config"
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/service/search/levenshtein"

	"gorm.io/gorm"
)

const sens float64 = 78

type SearchEngine struct {
	config   *config.Config
	database *gorm.DB
}

func NewEngine(config *config.Config, database *gorm.DB) *SearchEngine {
	return &SearchEngine{
		config:   config,
		database: database,
	}
}

func (e *SearchEngine) SearchProducts(query string) ([]models.Product, error) {
	query = prepareQueryString(query)

	var products []models.Product
	err := e.database.Find(&products).Error
	if err != nil {
		return nil, err
	}

	var results []models.Product
	for _, product := range products {
		if levenshtein.Fuzzy(query, prepareQueryString(product.Name)) > e.config.SearchSensitive {
			results = append(results, product)
		}
	}

	return results, nil
}
