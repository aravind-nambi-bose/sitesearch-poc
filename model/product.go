package model

import "time"

type Product struct {
	BaseProduct string `json:"baseProduct"`
	Name        string `json:"name"`
	ProductType string `json:"productType"`
	//This will be CE/Pro/etc
	Category   string    `json:"category"`
	LaunchDate time.Time `json:"launchDate"`
	Price      float64   `json:"price"`

	Description     string  `json:"description"`
	DiscountPrice   float64 `json:"discountPrice"`
	DiscountPercent float64 `json:"discountPercent"`
	//This needs to be more precise than what we get from Product Information API
	ProductCategories []string `json:"productCategories"`
	AverageRating     float64  `json:"averageRating"`
	ReviewsTotal      int      `json:"reviewsTotal"`
	Promote           bool     `json:"promote"`
}
