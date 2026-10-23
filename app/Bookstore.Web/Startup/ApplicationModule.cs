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

namespace Bookstore.Web
{
    /// <summary>
    /// Port of the .NET Framework App_Start/DependencyInjectionSetup.cs. Every registration of the
    /// original container survives here; the single Autofac container is created by
    /// <c>AutofacServiceProviderFactory</c> in Program.cs (no legacy MVC5 dependency resolver, no second
    /// container). Controller registration is no longer required — ASP.NET Core activates
    /// controllers itself through AddControllersWithViews().
    /// </summary>
    public class ApplicationModule : Module
    {
        private readonly string _webRootPath;

        public ApplicationModule(string webRootPath)
        {
            _webRootPath = webRootPath;
        }

        protected override void Load(ContainerBuilder builder)
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

            // The EF Core ApplicationDbContext takes DbContextOptions instead of a raw
            // connection string, so the connection string is applied through UseSqlServer.
            // InstancePerRequest (MVC5/OWIN) becomes InstancePerLifetimeScope (scoped).
            builder.Register(c =>
            {
                var options = new DbContextOptionsBuilder<ApplicationDbContext>();

                options.UseSqlServer(connectionString);

                return new ApplicationDbContext(options.Options);
            }).AsSelf().InstancePerLifetimeScope();

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
                // The legacy ASP.NET app-domain path APIs no longer exist; the web
                // root supplied by IWebHostEnvironment is the ASP.NET Core equivalent and the
                // cover images now live under wwwroot/images/coverimages.
                builder.RegisterInstance(new LocalFileService(_webRootPath)).As<IFileService>();
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
                builder.RegisterType<LocalAuthenticationMiddleware>().AsSelf().InstancePerLifetimeScope();
            }
        }
    }
}
