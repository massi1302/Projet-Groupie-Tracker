package controllers

import (
	"fmt"
	"time"
)

func getReliableEarthData(lat, lon, date string) []EarthData {
	if date == "" {
		date = time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	}

	// Try Sentinel Hub or another satellite imagery provider's API first
	sentinelData := tryAlternativeSatelliteSource(lat, lon, date)
	if len(sentinelData) > 0 {
		return sentinelData
	}

	// If that fails, try NASA's GIBS (Global Imagery Browse Services) as a fallback
	gibsData := tryNASAGIBS(lat, lon, date)
	if len(gibsData) > 0 {
		return gibsData
	}

	// As a last resort, try a static map service using the coordinates
	return generateStaticMapImage(lat, lon)
}

// Try an alternative satellite imagery source like Sentinel Hub
func tryAlternativeSatelliteSource(lat, lon, date string) []EarthData {
	// This would normally use Sentinel Hub or similar API
	// For demo purposes, we'll simulate a successful response with Mapbox satellite imagery

	// Format the coordinates
	coords := fmt.Sprintf("Lat: %s, Lon: %s", lat, lon)

	// Generate a Mapbox satellite imagery URL
	zoom := "12"
	mapboxToken := "pk.eyJ1IjoibWFzc2ktMTMiLCJhIjoiY203cnZvb3l4MTBzcTJqcXk4cXkzYW1ldSJ9.17SowS5wR2AGGR8Al9U5dA" // Replace with your actual token
	imageURL := fmt.Sprintf(
		"https://api.mapbox.com/styles/v1/mapbox/satellite-v9/static/%s,%s,%s,0/600x400?access_token=%s",
		lon, lat, zoom, mapboxToken,
	)

	// Create Earth data
	title := fmt.Sprintf("Satellite View")

	return []EarthData{
		{
			ID:          fmt.Sprintf("satellite-%s-%s", lat, lon),
			Title:       title,
			URL:         imageURL,
			Date:        date,
			Coordinates: coords,
			Type:        "image",
			Dataset:     "Satellite Imagery",
			Description: fmt.Sprintf("Satellite image at coordinates %s taken on or before %s", coords, date),
			Keywords:    []string{"Earth", "Satellite", "Geography"},
			Source:      "Satellite Imagery",
		},
	}
}

// Try NASA's Global Imagery Browse Services (GIBS)
func tryNASAGIBS(lat, lon, date string) []EarthData {
	// NASA GIBS provides various layers of Earth imagery
	// This is a simplified example using a static MODIS layer

	// Format the coordinates
	coords := fmt.Sprintf("Lat: %s, Lon: %s", lat, lon)

	// Parse date for GIBS URL format
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		parsedDate = time.Now()
	}

	// Format date for GIBS URL
	gibsDate := parsedDate.Format("2006-01-02")

	// Create a GIBS URL for MODIS Terra imagery
	// Note: A real implementation would use the GIBS API properly
	imageURL := fmt.Sprintf(
		"https://gibs.earthdata.nasa.gov/wmts/epsg4326/best/MODIS_Terra_CorrectedReflectance_TrueColor/default/%sT00:00:00Z/250m/6/%s/%s.jpg",
		gibsDate, lat, lon,
	)

	return []EarthData{
		{
			ID:          fmt.Sprintf("gibs-%s-%s", lat, lon),
			Title:       fmt.Sprintf("NASA GIBS View"),
			URL:         imageURL,
			Date:        date,
			Coordinates: coords,
			Type:        "image",
			Dataset:     "NASA GIBS",
			Description: fmt.Sprintf("NASA GIBS satellite image at coordinates %s from %s", coords, date),
			Keywords:    []string{"Earth", "NASA", "MODIS", "Geography"},
			Source:      "NASA GIBS",
		},
	}
}

// Generate a static map image as a last resort
func generateStaticMapImage(lat, lon string) []EarthData {
	// Format the coordinates
	coords := fmt.Sprintf("Lat: %s, Lon: %s", lat, lon)

	// Use OpenStreetMap or Google Maps static image as fallback
	// This is a reliable fallback that should always work
	zoom := "10"
	imageURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/staticmap?center=%s,%s&zoom=%s&size=600x400&maptype=satellite&key=AIzaSyA_DlCuRTJ_hm7xayDTXy3-l9Qed6jQS70",
		lat, lon, zoom,
	)

	date := time.Now().Format("2006-01-02")

	return []EarthData{
		{
			ID:          fmt.Sprintf("map-%s-%s", lat, lon),
			Title:       fmt.Sprintf("Map Satellite View"),
			URL:         imageURL,
			Date:        date,
			Coordinates: coords,
			Type:        "image",
			Dataset:     "Map Satellite",
			Description: fmt.Sprintf("Satellite map image at coordinates %s", coords),
			Keywords:    []string{"Earth", "Map", "Satellite", "Geography"},
			Source:      "Maps API",
		},
	}
}
