package dto

type CreateAddressRequest struct {
	AddressLine1 string `json:"addressLine1" binding:"required"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	Country      string `json:"country" binding:"required"`
	ZipCode      string `json:"zipCode" binding:"required"`
}

type UpdateAddressRequest struct {
	ID           int    `json:"id" binding:"required"`
	AddressLine1 string `json:"addressLine1" binding:"required"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	Country      string `json:"country" binding:"required"`
	ZipCode      string `json:"zipCode" binding:"required"`
}

type PlaceOrderRequest struct {
	SelectedAddressID int `json:"selectedAddressId" binding:"required"`
}

type CreateOfferRequest struct {
	BookName    string  `json:"bookName" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	ISBN        string  `json:"isbn"`
	BookTypeID  int     `json:"bookTypeId" binding:"required"`
	ConditionID int     `json:"conditionId" binding:"required"`
	GenreID     int     `json:"genreId" binding:"required"`
	PublisherID int     `json:"publisherId" binding:"required"`
	BookPrice   float64 `json:"bookPrice" binding:"required"`
	Summary     string  `json:"summary"`
}

type CreateBookRequest struct {
	Name        string  `json:"name" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	BookTypeID  int     `json:"bookTypeId" binding:"required"`
	ConditionID int     `json:"conditionId" binding:"required"`
	GenreID     int     `json:"genreId" binding:"required"`
	PublisherID int     `json:"publisherId" binding:"required"`
	Year        *int    `json:"year"`
	ISBN        string  `json:"isbn"`
	Summary     string  `json:"summary"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity"`
}

type UpdateBookRequest struct {
	ID          int     `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	BookTypeID  int     `json:"bookTypeId" binding:"required"`
	ConditionID int     `json:"conditionId" binding:"required"`
	GenreID     int     `json:"genreId" binding:"required"`
	PublisherID int     `json:"publisherId" binding:"required"`
	Year        *int    `json:"year"`
	ISBN        string  `json:"isbn"`
	Summary     string  `json:"summary"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity"`
}

type UpdateOrderStatusRequest struct {
	OrderID     int `json:"orderId" binding:"required"`
	OrderStatus int `json:"orderStatus" binding:"required"`
}

type CreateReferenceDataRequest struct {
	DataType int    `json:"dataType" binding:"required"`
	Text     string `json:"text" binding:"required"`
}

type UpdateReferenceDataRequest struct {
	ID       int    `json:"id" binding:"required"`
	DataType int    `json:"dataType" binding:"required"`
	Text     string `json:"text" binding:"required"`
}

type DashboardResponse struct {
	TotalOrders    int64 `json:"totalOrders"`
	PendingOrders  int64 `json:"pendingOrders"`
	PendingOffers  int64 `json:"pendingOffers"`
	TotalBooks     int64 `json:"totalBooks"`
	LowStockBooks  int64 `json:"lowStockBooks"`
}
