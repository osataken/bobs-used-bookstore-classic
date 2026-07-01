package model

type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusOrdered   OrderStatus = 1
	OrderStatusShipped   OrderStatus = 2
	OrderStatusDelivered OrderStatus = 3
	OrderStatusCancelled OrderStatus = 4
)

type OfferStatus int

const (
	OfferStatusPendingApproval OfferStatus = 0
	OfferStatusApproved        OfferStatus = 1
	OfferStatusReceived        OfferStatus = 2
	OfferStatusPaid            OfferStatus = 3
	OfferStatusRejected        OfferStatus = 4
)

type ReferenceDataType int

const (
	ReferenceDataTypePublisher  ReferenceDataType = 0
	ReferenceDataTypeCondition  ReferenceDataType = 1
	ReferenceDataTypeBookType   ReferenceDataType = 2
	ReferenceDataTypeGenre      ReferenceDataType = 3
)
