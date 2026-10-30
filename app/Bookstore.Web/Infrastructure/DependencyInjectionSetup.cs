using System.IO;
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
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;

namespace Bookstore.Web.Infrastructure
{
    /// <summary>
    /// Ported from App_Start/DependencyInjectionSetup.cs. The OWIN/Autofac.Mvc5 integration
    /// (DependencyResolver.SetResolver + UseAutofacMiddleware + RegisterControllers) is replaced by
    /// the Autofac service-provider factory; every original registration is preserved.
    /// </summary>
    public static class DependencyInjectionSetup
    {
        /// <summary>
        /// Registrations that must live in the framework service collection: the EF Core DbContext
        /// and the AWS SDK clients (conditional, exactly as in the original container).
        /// </summary>
        public static void ConfigureServices(IServiceCollection services)
        {
            var connectionString = BookstoreConfiguration.GetConnectionString("BookstoreDatabaseConnection");

            // Was: RegisterType<ApplicationDbContext>().WithParameter("connectionString", ...).InstancePerRequest()
            services.AddDbContext<ApplicationDbContext>(options => options.UseSqlServer(connectionString));

            if (BookstoreConfiguration.GetSetting("Services/FileService") == "aws")
            {
                services.AddAWSService<IAmazonS3>();
            }

            if (BookstoreConfiguration.GetSetting("Services/ImageValidationService") == "aws")
            {
                services.AddAWSService<IAmazonRekognition>();
            }
        }

        /// <summary>
        /// The Autofac registrations, one-for-one with the original container.
        /// </summary>
        /// <param name="webRootPath">Web root used by <see cref="LocalFileService"/> (was the HttpRuntime probe).</param>
        public static void ConfigureContainer(ContainerBuilder builder, string webRootPath)
        {
            builder.RegisterType<BookService>().As<IBookService>();
            builder.RegisterType<OrderService>().As<IOrderService>();
            builder.RegisterType<ReferenceDataService>().As<IReferenceDataService>();
            builder.RegisterType<OfferService>().As<IOfferService>();
            builder.RegisterType<CustomerService>().As<ICustomerService>();
            builder.RegisterType<AddressService>().As<IAddressService>();
            builder.RegisterType<ShoppingCartService>().As<IShoppingCartService>();
            builder.RegisterType<ImageResizeService>().As<IImageResizeService>();

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
                // IAmazonS3 itself comes from AWSSDK.Extensions.NETCore.Setup (see ConfigureServices).
                builder.RegisterType<S3FileService>().As<IFileService>();
            }
            else
            {
                builder.RegisterInstance(new LocalFileService(webRootPath)).As<IFileService>();
            }

            if (BookstoreConfiguration.GetSetting("Services/ImageValidationService") == "aws")
            {
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
