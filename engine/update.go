package engine

import (
	"github.com/BurntSushi/toml"
	"github.com/zhenghaoz/gorse/base"
	"github.com/zhenghaoz/gorse/core"
	"log"
)

func UpdateItemPop(cacheSize int, dataSet core.DataSetInterface, db *DB) error {
	log.Printf("update popular items")
	items, scores := core.Popular(dataSet, cacheSize)
	recommends := make([]RecommendedItem, len(items))
	for i := range recommends {
		recommends[i].ItemId = items[i]
		recommends[i].Score = scores[i]
	}
	if err := db.SetPopular(recommends); err != nil {
		return err
	}
	return nil
}

func UpdateNeighbors(name string, cacheSize int, dataSet core.DataSetInterface, db *DB) error {
	log.Printf("update neighbors by %v", name)
	similarity := LoadSimilarity(name)
	for denseItemId := 0; denseItemId < dataSet.ItemCount(); denseItemId++ {
		itemId := dataSet.ItemIndexer().ToID(denseItemId)
		neighbors, similarities := core.Neighbors(dataSet, itemId, cacheSize, similarity)
		recommends := make([]RecommendedItem, len(neighbors))
		for i := range recommends {
			recommends[i].ItemId = neighbors[i]
			recommends[i].Score = similarities[i]
		}
		if err := db.SetNeighbors(itemId, recommends); err != nil {
			return err
		}
	}
	return nil
}

func UpdateRecommends(name string, params base.Params, cacheSize int, dataSet core.DataSetInterface, db *DB) error {
	// Create model
	log.Printf("create model %v with params = %v\n", name, params)
	model := LoadModel(name, params)
	// Training model
	log.Println("training model")
	model.Fit(dataSet, nil)
	// Generate recommends
	log.Println("update recommends")
	items := core.Items(dataSet)
	for userIndex := 0; userIndex < dataSet.UserCount(); userIndex++ {
		userId := dataSet.UserIndexer().ToID(userIndex)
		exclude := dataSet.UserByIndex(userIndex)
		recommendItems, ratings := core.Top(items, userId, cacheSize, exclude, model)
		recommends := make([]RecommendedItem, len(recommendItems))
		for i := range recommends {
			recommends[i].ItemId = recommendItems[i]
			recommends[i].Score = ratings[i]
		}
		if err := db.SetRecommends(userId, recommends); err != nil {
			return err
		}
	}
	return nil
}

func Update(config TomlConfig, metaData toml.MetaData, db *DB) error {
	// Load data
	log.Println("load data from database")
	dataSet, err := db.ToDataSet()
	if err != nil {
		return err
	}
	// Generate recommends
	params := config.Params.ToParams(metaData)
	if err = UpdateRecommends(config.Recommend.Model, params, config.Recommend.CacheSize, dataSet, db); err != nil {
		return err
	}
	// Generate neighbors
	if err = UpdateNeighbors(config.Recommend.Similarity, config.Recommend.CacheSize, dataSet, db); err != nil {
		return err
	}
	// Generate popular items
	if err = UpdateItemPop(config.Recommend.CacheSize, dataSet, db); err != nil {
		return err
	}
	return nil
}
