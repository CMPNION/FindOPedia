package entity

import "time"

type RarityTier string

const (
	RarityCommon    RarityTier = "common"
	RarityUncommon  RarityTier = "uncommon"
	RarityRare      RarityTier = "rare"
	RarityEpic      RarityTier = "epic"
	RarityLegendary RarityTier = "legendary"
)

func ComputeRarity(contentLength int) RarityTier {
	switch {
	case contentLength < 5000:
		return RarityCommon
	case contentLength < 20000:
		return RarityUncommon
	case contentLength < 50000:
		return RarityRare
	case contentLength < 100000:
		return RarityEpic
	default:
		return RarityLegendary
	}
}

type Article struct {
	ID            int64
	WikipediaID   int
	Title         string
	Slug          string
	Content       string
	ContentLength int
	RarityTier    RarityTier
	Summary       string
	CreatedAt     time.Time
}
