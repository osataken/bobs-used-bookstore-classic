package dto

// BookResponse represents a book in API responses
type BookResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Author        string  `json:"author"`
	Year          *int    `json:"year"`
	ISBN          string  `json:"isbn"`
	PublisherID   int     `json:"publisherId"`
	Publisher     string  `json:"publisher"`
	BookTypeID    int     `json:"bookTypeId"`
	BookType      string  `json:"bookType"`
	GenreID       int     `json:"genreId"`
	Genre         string  `json:"genre"`
	ConditionID   int     `json:"conditionId"`
	Condition     string  `json:"condition"`
	CoverImageURL string  `json:"coverImageUrl"`
	Summary       string  `json:"summary"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	IsInStock     bool    `json:"isInStock"`
	IsLowInStock  bool    `json:"isLowInStock"`
}

// BookCreateRequest for creating a book (form data)
type BookCreateRequest struct {
	Name        string  `form:"name" binding:"required"`
	Author      string  `form:"author" binding:"required"`
	BookTypeID  int     `form:"bookTypeId" binding:"required"`
	ConditionID int     `form:"conditionId" binding:"required"`
	GenreID     int     `form:"genreId" binding:"required"`
	PublisherID int     `form:"publisherId" binding:"required"`
	Year        *int    `form:"year"`
	ISBN        string  `form:"isbn"`
	Summary     string  `form:"summary"`
	Price       float64 `form:"price" binding:"required"`
	Quantity    int     `form:"quantity" binding:"required"`
}

// BookUpdateRequest for updating a book (form data)
type BookUpdateRequest struct {
	ID          int     `form:"id" binding:"required"`
	Name        string  `form:"name" binding:"required"`
	Author      string  `form:"author" binding:"required"`
	BookTypeID  int     `form:"bookTypeId" binding:"required"`
	ConditionID int     `form:"conditionId" binding:"required"`
	GenreID     int     `form:"genreId" binding:"required"`
	PublisherID int     `form:"publisherId" binding:"required"`
	Year        *int    `form:"year"`
	ISBN        string  `form:"isbn"`
	Summary     string  `form:"summary"`
	Price       float64 `form:"price" binding:"required"`
	Quantity    int     `form:"quantity" binding:"required"`
}

// BookFiltersRequest for filtering books
type BookFiltersRequest struct {
	Name        string `form:"name"`
	Author      string `form:"author"`
	PublisherID *int   `form:"publisherId"`
	GenreID     *int   `form:"genreId"`
	BookTypeID  *int   `form:"bookTypeId"`
	ConditionID *int   `form:"conditionId"`
	LowStock    bool   `form:"lowStock"`
}

// BookStatisticsResponse for inventory statistics
type BookStatisticsResponse struct {
	OutOfStock int `json:"outOfStock"`
	LowStock   int `json:"lowStock"`
	StockTotal int `json:"stockTotal"`
}
