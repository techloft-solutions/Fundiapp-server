package sqlite

import (
	"context"
	"fmt"
	"math"
	"strconv"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
)

type SearchService struct {
	db *DB
}

func NewSearchService(db *DB) *SearchService {
	return &SearchService{db}
}

func (s *SearchService) SearchByQuery(ctx context.Context, search model.Search) ([]app.SearchResult, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	results, err := searchByQuery(ctx, tx, search)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func searchByQuery(ctx context.Context, tx *Tx, search model.Search) ([]app.SearchResult, error) {
	query := "%" + search.Query + "%"
	rows, err := tx.QueryContext(ctx, `
	SELECT
		providers.provider_id,
		CONCAT_WS(' ', first_name, last_name) AS name,
		users.photo_url,
		locations.latitude,
		locations.longitude
	FROM
		users
	LEFT JOIN providers ON providers.user_id = users.user_id
	LEFT JOIN categories ON categories.id = providers.category_id
	LEFT JOIN industries ON industries.id = providers.industry_id
	LEFT JOIN locations ON locations.location_id = users.location_id
	WHERE
		CONCAT_WS(
			'',
			users.first_name,
			users.last_name,
			categories.name,
			industries.name
		) LIKE(?)
		`,
		query,
	)
	if err != nil {
		return nil, err
	}

	var latitude *float64
	var longitude *float64

	defer rows.Close()

	results := make([]app.SearchResult, 0)
	for rows.Next() {
		var result app.SearchResult
		if err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.Photo,
			&latitude,
			&longitude,
		); err != nil {
			return nil, err
		}
		// get distance info
		var distance float64
		if search.Latitude != "" && search.Longitude != "" {
			searchLat, err := strconv.ParseFloat(search.Latitude, 64)
			if err != nil {
				return nil, err
			}
			searchLong, err := strconv.ParseFloat(search.Longitude, 64)
			if err != nil {
				return nil, err
			}
			if latitude != nil && longitude != nil {
				distance = calculateDistance(searchLat, searchLong, *latitude, *longitude)
				distanceStr := fmt.Sprintf("%.1f %s", distance, "Km")
				result.Distance = &distanceStr
			}
		}

		searchDistance, _ := strconv.ParseFloat(search.Distance, 64)
		// if search distance is less than result distance dont include result in results
		if searchDistance != 0 {
			if distance > searchDistance {
				continue
			}
		}

		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// calculateDistance calculates the distance between two points in km
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6371 // Earth radius in KM

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
