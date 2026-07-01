package service

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/repository"
)

type Services struct {
	Book          *BookService
	Order         *OrderService
	Cart          *CartService
	Customer      *CustomerService
	Address       *AddressService
	Offer         *OfferService
	ReferenceData *ReferenceDataService
	File          FileService
	Image         ImageValidationService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	var fileService FileService
	if cfg.FileServiceMode == "aws" {
		fileService = NewS3FileService(cfg)
	} else {
		fileService = NewLocalFileService(cfg)
	}

	var imageService ImageValidationService
	if cfg.ImageValidationMode == "aws" {
		imageService = &RekognitionImageValidationService{}
	} else {
		imageService = &LocalImageValidationService{}
	}

	return &Services{
		Book:          NewBookService(repos.Book, repos.Order, fileService, imageService),
		Order:         NewOrderService(repos.Order, repos.Cart, repos.Customer, repos.Book),
		Cart:          NewCartService(repos.Cart),
		Customer:      NewCustomerService(repos.Customer),
		Address:       NewAddressService(repos.Address, repos.Customer),
		Offer:         NewOfferService(repos.Offer, repos.Customer),
		ReferenceData: NewReferenceDataService(repos.ReferenceData),
		File:          fileService,
		Image:         imageService,
	}
}
