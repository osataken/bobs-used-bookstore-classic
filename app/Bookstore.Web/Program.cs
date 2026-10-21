using Amazon.Rekognition;
using Amazon.S3;
using Amazon.SimpleSystemsManagement;
using Amazon.SimpleSystemsManagement.Model;
using Autofac;
using Autofac.Extensions.DependencyInjection;
using BobsBookstoreClassic.Data;
using Bookstore.Common;
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
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authentication.OpenIdConnect;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.IdentityModel.Protocols.OpenIdConnect;
using Microsoft.IdentityModel.Tokens;
using NLog;
using NLog.AWS.Logger;
using NLog.Config;
using NLog.Web;
using System.Security.Claims;
using System.Threading.Tasks;
using NLogTarget = NLog.Targets.Target;

var builder = WebApplication.CreateBuilder(args);

// ── 1. Initialize BookstoreConfiguration from IConfiguration ──────────────
BookstoreConfiguration.Initialize(builder.Configuration);

// Load remote config from SSM Parameter Store if services are configured for AWS
LoadRemoteConfiguration(builder.Configuration);

// ── 2. Logging (NLog) ────────────────────────────────────────────────────
ConfigureNLog();
builder.Logging.ClearProviders();
builder.Host.UseNLog();

// ── 3. Autofac DI ────────────────────────────────────────────────────────
builder.Host.UseServiceProviderFactory(new AutofacServiceProviderFactory());
builder.Host.ConfigureContainer<ContainerBuilder>(ConfigureDependencyInjection);

// ── 4. MVC + Views ───────────────────────────────────────────────────────
builder.Services.AddControllersWithViews();
builder.Services.AddHttpContextAccessor();

// ── 5. EF Core ───────────────────────────────────────────────────────────
var connectionString = BookstoreConfiguration.GetConnectionString("BookstoreDatabaseConnection");
builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseSqlServer(connectionString));

// ── 6. Authentication ────────────────────────────────────────────────────
if (BookstoreConfiguration.GetSetting("Services/Authentication") == "aws")
{
    builder.Services
        .AddAuthentication(options =>
        {
            options.DefaultScheme = CookieAuthenticationDefaults.AuthenticationScheme;
            options.DefaultChallengeScheme = OpenIdConnectDefaults.AuthenticationScheme;
        })
        .AddCookie(CookieAuthenticationDefaults.AuthenticationScheme)
        .AddOpenIdConnect(OpenIdConnectDefaults.AuthenticationScheme, options =>
        {
            options.ClientId = BookstoreConfiguration.GetSetting("Authentication/Cognito/LocalClientId");
            options.MetadataAddress = BookstoreConfiguration.GetSetting("Authentication/Cognito/MetadataAddress");
            options.ResponseType = OpenIdConnectResponseType.Code;
            options.Scope.Add("openid");
            options.Scope.Add("profile");
            options.SaveTokens = true;
            options.UseTokenLifetime = false;
            options.TokenValidationParameters = new TokenValidationParameters
            {
                NameClaimType = "cognito:username",
                RoleClaimType = "cognito:groups"
            };
            options.Events = new OpenIdConnectEvents
            {
                OnRedirectToIdentityProvider = ctx =>
                {
                    ctx.ProtocolMessage.RedirectUri = $"{ctx.Request.Scheme}://{ctx.Request.Host}/signin-oidc";
                    return Task.CompletedTask;
                },
                OnAuthorizationCodeReceived = ctx =>
                {
                    ctx.TokenEndpointRequest.RedirectUri = $"{ctx.Request.Scheme}://{ctx.Request.Host}/signin-oidc";
                    return Task.CompletedTask;
                },
                OnTokenValidated = async ctx =>
                {
                    var service = ctx.HttpContext.RequestServices.GetRequiredService<ICustomerService>();
                    var identity = ctx.Principal?.Identity as ClaimsIdentity;
                    if (identity != null)
                    {
                        var dto = new CreateOrUpdateCustomerDto(
                            identity.GetSub() ?? string.Empty,
                            identity.Name ?? string.Empty,
                            identity.FindFirst(y => y.Type.Contains("givenname"))?.Value ?? string.Empty,
                            identity.FindFirst(y => y.Type.Contains("surname"))?.Value ?? string.Empty);
                        await service.CreateOrUpdateCustomerAsync(dto);
                    }
                }
            };
        });
}
else
{
    builder.Services
        .AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
        .AddCookie(CookieAuthenticationDefaults.AuthenticationScheme, options =>
        {
            options.LoginPath = "/Authentication/Login";
        });
}

// ── 7. Build ─────────────────────────────────────────────────────────────
var app = builder.Build();

// Initialize database
using (var scope = app.Services.CreateScope())
{
    var dbContext = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
    await BookstoreDbInitializer.InitializeAsync(dbContext);
}

// ── 8. Middleware Pipeline ────────────────────────────────────────────────
if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Home/Error");
}

app.UseStaticFiles();
app.UseRouting();

app.UseAuthentication();

if (BookstoreConfiguration.GetSetting("Services/Authentication") != "aws")
{
    app.UseMiddleware<LocalAuthenticationMiddleware>();
}

app.UseAuthorization();

// ── 9. Routing ────────────────────────────────────────────────────────────
app.MapAreaControllerRoute(
    name: "Admin",
    areaName: "Admin",
    pattern: "Admin/{controller=Dashboard}/{action=Index}/{id?}");

app.MapControllerRoute(
    name: "default",
    pattern: "{controller=Home}/{action=Index}/{id?}");

app.Run();

// ──────────────────────────────────────────────────────────────────────────
// Local helper methods
// ──────────────────────────────────────────────────────────────────────────

static void ConfigureNLog()
{
    var config = new LoggingConfiguration();
    NLogTarget loggingTarget;

    if (BookstoreConfiguration.GetSetting("Services/LoggingService") == "aws")
    {
        loggingTarget = new AWSTarget { LogGroup = Constants.AppName };
    }
    else
    {
        loggingTarget = new NLog.Targets.DebuggerTarget();
    }

    config.AddTarget("target", loggingTarget);
    config.LoggingRules.Add(new LoggingRule("*", NLog.LogLevel.Info, loggingTarget));
    LogManager.Configuration = config;
}

static void LoadRemoteConfiguration(IConfiguration configuration)
{
    var rootPath = "/" + Constants.AppName;

    if (BookstoreConfiguration.GetSetting("Services/Database") == "aws")
    {
        using var client = new AmazonSimpleSystemsManagementClient();
        var request = new GetParameterRequest
        {
            Name = $"{rootPath}/Database/ConnectionStrings/BookstoreDatabaseConnection"
        };
        var response = client.GetParameterAsync(request).GetAwaiter().GetResult();
        BookstoreConfiguration.AddConnectionString("BookstoreDatabaseConnection", response.Parameter.Value);
    }

    if (BookstoreConfiguration.GetSetting("Services/Authentication") == "aws")
    {
        using var client = new AmazonSimpleSystemsManagementClient();
        var request = new GetParametersByPathRequest { Path = $"{rootPath}/Authentication/", Recursive = true };
        var response = client.GetParametersByPathAsync(request).GetAwaiter().GetResult();
        foreach (var parameter in response.Parameters)
        {
            BookstoreConfiguration.AddSetting(parameter.Name.Replace($"{rootPath}/", string.Empty), parameter.Value);
        }
    }

    if (BookstoreConfiguration.GetSetting("Services/FileService") == "aws")
    {
        using var client = new AmazonSimpleSystemsManagementClient();
        var request = new GetParametersByPathRequest { Path = $"{rootPath}/Files/", Recursive = true };
        var response = client.GetParametersByPathAsync(request).GetAwaiter().GetResult();
        foreach (var parameter in response.Parameters)
        {
            BookstoreConfiguration.AddSetting(parameter.Name.Replace($"{rootPath}/", string.Empty), parameter.Value);
        }
    }
}

static void ConfigureDependencyInjection(ContainerBuilder builder)
{
    // Domain services
    builder.RegisterType<BookService>().As<IBookService>();
    builder.RegisterType<OrderService>().As<IOrderService>();
    builder.RegisterType<ReferenceDataService>().As<IReferenceDataService>();
    builder.RegisterType<OfferService>().As<IOfferService>();
    builder.RegisterType<CustomerService>().As<ICustomerService>();
    builder.RegisterType<AddressService>().As<IAddressService>();
    builder.RegisterType<ShoppingCartService>().As<IShoppingCartService>();
    builder.RegisterType<ImageResizeService>().As<IImageResizeService>();

    // Repositories
    builder.RegisterType<CustomerRepository>().As<ICustomerRepository>();
    builder.RegisterType<AddressRepository>().As<IAddressRepository>();
    builder.RegisterType<BookRepository>().As<IBookRepository>();
    builder.RegisterType<OfferRepository>().As<IOfferRepository>();
    builder.RegisterType<ShoppingCartRepository>().As<IShoppingCartRepository>();
    builder.RegisterType<OrderRepository>().As<IOrderRepository>();
    builder.RegisterType<ReferenceDataRepository>().As<IReferenceDataRepository>();

    builder.RegisterGeneric(typeof(PaginatedList<>)).As(typeof(IPaginatedList<>)).InstancePerLifetimeScope();

    // File service
    if (BookstoreConfiguration.GetSetting("Services/FileService") == "aws")
    {
        builder.RegisterType<AmazonS3Client>().As<IAmazonS3>();
        builder.RegisterType<S3FileService>().As<IFileService>();
    }
    else
    {
        builder.RegisterType<LocalFileService>()
            .WithParameter("webRootPath", "wwwroot")
            .As<IFileService>();
    }

    // Image validation service
    if (BookstoreConfiguration.GetSetting("Services/ImageValidationService") == "aws")
    {
        builder.RegisterType<AmazonRekognitionClient>().As<IAmazonRekognition>();
        builder.RegisterType<RekognitionImageValidationService>().As<IImageValidationService>();
    }
    else
    {
        builder.RegisterType<LocalImageValidationService>().As<IImageValidationService>();
    }

    // Local auth middleware
    if (BookstoreConfiguration.GetSetting("Services/Authentication") != "aws")
    {
        builder.RegisterType<LocalAuthenticationMiddleware>();
    }
}
