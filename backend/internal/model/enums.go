package model

// OrderStatus enum
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusOrdered   OrderStatus = 1
	OrderStatusShipped   OrderStatus = 2
	OrderStatusDelivered OrderStatus = 3
	OrderStatusCancelled OrderStatus = 4
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "Pending"
	case OrderStatusOrdered:
		return "Ordered"
	case OrderStatusShipped:
		return "Shipped"
	case OrderStatusDelivered:
		return "Delivered"
	case OrderStatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// OfferStatus enum
type OfferStatus int

const (
	OfferStatusPendingApproval OfferStatus = 0
	OfferStatusApproved        OfferStatus = 1
	OfferStatusReceived        OfferStatus = 2
	OfferStatusPaid            OfferStatus = 3
	OfferStatusRejected        OfferStatus = 4
)

func (s OfferStatus) String() string {
	switch s {
	case OfferStatusPendingApproval:
		return "Pending Approval"
	case OfferStatusApproved:
		return "Approved"
	case OfferStatusReceived:
		return "Received"
	case OfferStatusPaid:
		return "Paid"
	case OfferStatusRejected:
		return "Rejected"
	default:
		return "Unknown"
	}
}

// ReferenceDataType enum
type ReferenceDataType int

const (
	ReferenceDataTypePublisher  ReferenceDataType = 0
	ReferenceDataTypeCondition  ReferenceDataType = 1
	ReferenceDataTypeBookType   ReferenceDataType = 2
	ReferenceDataTypeGenre      ReferenceDataType = 3
)

func (t ReferenceDataType) String() string {
	switch t {
	case ReferenceDataTypePublisher:
		return "Publisher"
	case ReferenceDataTypeCondition:
		return "Condition"
	case ReferenceDataTypeBookType:
		return "BookType"
	case ReferenceDataTypeGenre:
		return "Genre"
	default:
		return "Unknown"
	}
}
