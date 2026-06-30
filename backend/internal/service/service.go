package service

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/repository"
)

// Services holds all service instances
type Services struct {
	Book          *BookService
	Order         *OrderService
	Customer      *CustomerService
	Address       *AddressService
	Offer         *OfferService
	ShoppingCart  *ShoppingCartService
	ReferenceData *ReferenceDataService
	FileService   FileService
	ImageResize   *ImageResizeService
	ImageValidation ImageValidationService
}

// NewServices creates all services
func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	var fileService FileService
	if cfg.IsAWS("fileservice") {
		fileService = NewS3FileService(cfg)
	} else {
		fileService = NewLocalFileService(cfg.LocalImagePath)
	}

	var imageValidation ImageValidationService
	if cfg.IsAWS("imagevalidation") {
		imageValidation = NewRekognitionImageValidationService()
	} else {
		imageValidation = NewLocalImageValidationService()
	}

	imageResize := NewImageResizeService()

	return &Services{
		Book:            NewBookService(repos.Book, repos.Order, imageResize, imageValidation, fileService),
		Order:           NewOrderService(repos.Order, repos.ShoppingCart, repos.Customer, repos.Book),
		Customer:        NewCustomerService(repos.Customer),
		Address:         NewAddressService(repos.Address, repos.Customer),
		Offer:           NewOfferService(repos.Offer, repos.Customer),
		ShoppingCart:    NewShoppingCartService(repos.ShoppingCart),
		ReferenceData:   NewReferenceDataService(repos.ReferenceData),
		FileService:     fileService,
		ImageResize:     imageResize,
		ImageValidation: imageValidation,
	}
}
