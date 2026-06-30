package dto

// OfferResponse represents an offer in API responses
type OfferResponse struct {
	ID          int     `json:"id"`
	BookName    string  `json:"bookName"`
	Author      string  `json:"author"`
	ISBN        string  `json:"isbn"`
	GenreID     int     `json:"genreId"`
	Genre       string  `json:"genre"`
	ConditionID int     `json:"conditionId"`
	Condition   string  `json:"condition"`
	PublisherID int     `json:"publisherId"`
	Publisher   string  `json:"publisher"`
	BookTypeID  int     `json:"bookTypeId"`
	BookType    string  `json:"bookType"`
	Summary     string  `json:"summary"`
	OfferStatus int     `json:"offerStatus"`
	StatusText  string  `json:"statusText"`
	Comment     string  `json:"comment"`
	CustomerID  int     `json:"customerId"`
	BookPrice   float64 `json:"bookPrice"`
}

// OfferCreateRequest for creating an offer
type OfferCreateRequest struct {
	BookName    string  `json:"bookName" binding:"required"`
	Author      string  `json:"author"`
	ISBN        string  `json:"isbn"`
	BookTypeID  int     `json:"bookTypeId" binding:"required"`
	ConditionID int     `json:"conditionId" binding:"required"`
	GenreID     int     `json:"genreId" binding:"required"`
	PublisherID int     `json:"publisherId" binding:"required"`
	BookPrice   float64 `json:"bookPrice" binding:"required"`
	Summary     string  `json:"summary"`
}

// OfferFiltersRequest for filtering offers
type OfferFiltersRequest struct {
	BookName    string `form:"bookName"`
	Author      string `form:"author"`
	GenreID     *int   `form:"genreId"`
	ConditionID *int   `form:"conditionId"`
	OfferStatus *int   `form:"offerStatus"`
}

// OfferStatisticsResponse for dashboard stats
type OfferStatisticsResponse struct {
	PendingOffers   int `json:"pendingOffers"`
	OffersThisMonth int `json:"offersThisMonth"`
	OffersTotal     int `json:"offersTotal"`
}
