using Amazon.Rekognition;
using Amazon.S3;
using Autofac;
using BobsBookstoreClassic.Data;
using Bookstore.Data;
using Bookstore.Data.FileServices;
using Bookstore.Data.ImageResizeService;
using Bookstore.Data.ImageValidationServices;
using Bookstore.Data.Repositories;
using Bookstore.Domain;
using Bookstore.Domain.Addresses;
using Bookstore.Domain.Books;
using Bookstore.Domain.Carts;
using Bookstore.Domain.Customers;
using Bookstore.Domain.Offers;
using Bookstore.Domain.Orders;
using Bookstore.Domain.ReferenceData;
using Bookstore.Web.Helpers;
using Microsoft.AspNetCore.Hosting;

namespace Bookstore.Web
{
    /// <summary>
    /// Autofac container composition, ported from App_Start/DependencyInjectionSetup.cs.
    /// Registered through <c>AutofacServiceProviderFactory</c>, so controllers and MVC filters are
    /// resolved by the ASP.NET Core infrastructure (the MVC5-era RegisterControllers /
    /// RegisterFilterProvider / AutofacDependencyResolver / UseAutofacMiddleware calls are no
    /// longer required).
    /// </summary>
    public static class DependencyInjectionSetup
    {
        public static void ConfigureDependencyInjection(ContainerBuilder builder, IWebHostEnvironment environment)
        {
            builder.RegisterType<BookService>().As<IBookService>();
            builder.RegisterType<OrderService>().As<IOrderService>();
            builder.RegisterType<ReferenceDataService>().As<IReferenceDataService>();
            builder.RegisterType<OfferService>().As<IOfferService>();
            builder.RegisterType<CustomerService>().As<ICustomerService>();
            builder.RegisterType<AddressService>().As<IAddressService>();
            builder.RegisterType<ShoppingCartService>().As<IShoppingCartService>();
            builder.RegisterType<ImageResizeService>().As<IImageResizeService>();

            var connectionString = BookstoreConfiguration.GetConnectionString("BookstoreDatabaseConnection");
            builder.RegisterType<ApplicationDbContext>().WithParameter("connectionString", connectionString).InstancePerLifetimeScope();

            builder.RegisterType<CustomerRepository>().As<ICustomerRepository>();
            builder.RegisterType<AddressRepository>().As<IAddressRepository>();
            builder.RegisterType<BookRepository>().As<IBookRepository>();
            builder.RegisterType<OfferRepository>().As<IOfferRepository>();
            builder.RegisterType<ShoppingCartRepository>().As<IShoppingCartRepository>();
            builder.RegisterType<OrderRepository>().As<IOrderRepository>();
            builder.RegisterType<ReferenceDataRepository>().As<IReferenceDataRepository>();

            builder.RegisterGeneric(typeof(PaginatedList<>)).As(typeof(IPaginatedList<>)).InstancePerLifetimeScope();

            if (BookstoreConfiguration.GetSetting("Services/FileService") == "aws")
            {
                builder.RegisterType<AmazonS3Client>().As<IAmazonS3>();
                builder.RegisterType<S3FileService>().As<IFileService>();
            }
            else
            {
                // Replaces HttpRuntime.AppDomainAppPath + "Content": locally uploaded book images
                // are written to, and served from, the static-file web root (wwwroot).
                var webRootPath = environment.WebRootPath ?? environment.ContentRootPath;

                builder.RegisterInstance(new LocalFileService(webRootPath)).As<IFileService>();
            }

            if (BookstoreConfiguration.GetSetting("Services/ImageValidationService") == "aws")
            {
                builder.RegisterType<AmazonRekognitionClient>().As<IAmazonRekognition>();
                builder.RegisterType<RekognitionImageValidationService>().As<IImageValidationService>();
            }
            else
            {
                builder.RegisterType<LocalImageValidationService>().As<IImageValidationService>();
            }

            if (BookstoreConfiguration.GetSetting("Services/Authentication") != "aws")
            {
                builder.RegisterType<LocalAuthenticationMiddleware>();
            }
        }
    }
}
