package controllers

import "encoding/xml"

type CategoryData struct {
	CategoryName string
	Items        []UnifiedImageData
	Categories   []string
	Resources    []interface{}
}

type CategoryInfo struct {
	Name     string
	ImageURL string
}

type MarsRoverResponse struct {
	Photos []struct {
		ImgSrc    string `json:"img_src"`
		EarthDate string `json:"earth_date"`
		Camera    struct {
			FullName string `json:"full_name"`
		} `json:"camera"`
		Rover struct {
			Name string `json:"name"`
		} `json:"rover"`
	} `json:"photos"`
}

type UnifiedImageData struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Date        string   `json:"date"`
	Explanation string   `json:"explanation"`
	Source      string   `json:"source"`
	Type        string   `json:"media_type"`
	Keywords    []string `json:"keywords"`
	ID          string   `json:"id"`
}

type CollectionData struct {
	Items []UnifiedImageData
	Error string
}

type RoverConfig struct {
	Name string
	Sols []int
}

// Structure pour les données RSS de la NASA
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Enclosure   struct {
		URL    string `xml:"url,attr"`
		Type   string `xml:"type,attr"`
		Length string `xml:"length,attr"`
	} `xml:"enclosure"`
}

// Structure pour l'API NASA Images
type ImageLibraryResponse struct {
	Collection struct {
		Items []struct {
			Data []struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				DateCreated string `json:"date_created"`
			} `json:"data"`
			Links []struct {
				Href string `json:"href"`
			} `json:"links"`
		} `json:"items"`
	} `json:"collection"`
}

type TopicResponse struct {
	Collection struct {
		Items []struct {
			Data []struct {
				Title       string   `json:"title"`
				Description string   `json:"description"`
				Keywords    []string `json:"keywords"`
				MediaType   string   `json:"media_type"`
			} `json:"data"`
			Links []struct {
				Href string `json:"href"`
			} `json:"links"`
		} `json:"items"`
	} `json:"collection"`
}

type Topic struct {
	Title       string
	Description string
	ImageURL    string
	Keywords    []string
}

type HomePageData struct {
	APOD    APODResponse
	News    []NewsItem
	Gallery []GalleryItem
	Topics  []Topic
}

type APODResponse struct {
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	HDURL       string `json:"hdurl"`
	MediaType   string `json:"media_type"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

type NewsItem struct {
	Title       string
	Description string
	ImageURL    string
	PublishDate string
	ReadMoreURL string
}

type GalleryItem struct {
	Title       string
	Description string
	ImageURL    string
	DateCreated string
	ID          string
}

// Structure for Earth Imagery API response
type EarthImageryResponse struct {
	Date       string  `json:"date"`
	ID         string  `json:"id"`
	URL        string  `json:"url"`
	CloudScore float64 `json:"cloud_score,omitempty"`
	Resource   struct {
		Dataset string `json:"dataset"`
		Planet  string `json:"planet"`
	} `json:"resource"`
}

type EarthAsset struct {
	Href string `json:"href"`
}

// Unified data structure for Earth imagery
type EarthData struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Date        string   `json:"date"`
	Coordinates string   `json:"coordinates"`
	Type        string   `json:"type"`
	Dataset     string   `json:"dataset"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Source      string   `json:"source"`
}

// And this struct for the Earth page data
type EarthPageData struct {
	Imagery         []EarthData `json:"imagery"`
	CurrentPage     int         `json:"currentPage"`
	TotalPages      int         `json:"totalPages"`
	PaginationRange []int       `json:"paginationRange"`
	PreviousPage    int         `json:"previousPage"`
	NextPage        int         `json:"nextPage"`
}

type EarthAssetsResponse struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Date     string `json:"date"`
	Resource struct {
		Dataset string `json:"dataset"`
	} `json:"resource"`
	Assets map[string]struct {
		Href string `json:"href"`
	} `json:"assets"`
}
