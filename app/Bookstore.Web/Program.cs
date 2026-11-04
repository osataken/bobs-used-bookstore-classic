using Amazon.Rekognition;
using Amazon.S3;
using Amazon.SimpleSystemsManagement;
using Amazon.SimpleSystemsManagement.Model;
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
using Microsoft.AspNetCore.HttpOverrides;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Protocols.OpenIdConnect;
using Microsoft.IdentityModel.Tokens;
using NLog;
using NLog.AWS.Logger;
using NLog.Config;
using NLog.Extensions.Logging;
using System.Security.Claims;

var builder = WebApplication.CreateBuilder(args);

// ─── Logging (NLog) ───────────────────────────────────────────────────────────
LogManager.Setup().LoadConfigurationFromSection(builder.Configuration);
builder.Logging.ClearProviders();
builder.Logging.AddNLog();

var startupLogger = LogManager.GetLogger("Startup");

// ─── Bridge IConfiguration → BookstoreConfiguration ──────────────────────────
// Populate the legacy BookstoreConfiguration from appsettings.json / env vars.
// Keys use ':' in IConfiguration; BookstoreConfiguration uses '/' as separator.
foreach (var config in builder.Configuration.AsEnumerable())
{
    if (config.Value != null)
        BookstoreConfiguration.AddSetting(config.Key.Replace(':', '/'), config.Value);
}
foreach (var cs in builder.Configuration.GetSection("ConnectionStrings").GetChildren())
{
    if (cs.Value != null)
        BookstoreConfiguration.AddConnectionString(cs.Key, cs.Value);
}

// ─── Load AWS-managed configuration (SSM Parameter Store) ────────────────────
LoadAwsConfiguration(startupLogger);

// ─── CloudWatch logging target ────────────────────────────────────────────────
// Ports App_Start/LoggingSetup.cs: when Services/LoggingService == "aws" the app
// logged to a CloudWatch log group named after the application. The appsettings
// "NLog" section only defines console/file targets, so the AWS target is added
// programmatically after the section has been loaded. NLog "autoReload" is set
// to false in appsettings.json: a reload would rebuild the configuration from
// the section alone and silently drop this target.
ConfigureAwsLoggingTarget();

// ─── EF Core / Database ───────────────────────────────────────────────────────
var connectionString = BookstoreConfiguration.GetConnectionString("BookstoreDatabaseConnection");
builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseSqlServer(connectionString));

// ─── Forwarded headers (ALB / reverse proxy TLS termination) ─────────────────
// Without this the request scheme behind the ECS load balancer resolves to
// "http", which breaks the Cognito redirect URI and any absolute URL building.
builder.Services.Configure<ForwardedHeadersOptions>(options =>
{
    options.ForwardedHeaders = ForwardedHeaders.XForwardedFor | ForwardedHeaders.XForwardedProto;
    options.KnownIPNetworks.Clear();
    options.KnownProxies.Clear();
});

// ─── MVC ──────────────────────────────────────────────────────────────────────
builder.Services.AddControllersWithViews(options =>
{
    // Global authorize filter - requires authentication for all controllers by default
    // (ports the MVC5 global AuthorizeAttribute from App_Start/FilterConfig.cs).
    options.Filters.Add(new Microsoft.AspNetCore.Mvc.Authorization.AuthorizeFilter());

    // MVC5 routed to the literal action method name, so URLs such as
    // /Admin/Offers/ApproveAsync/{id} (used as formaction values in the Admin
    // Offers view) must keep answering. ASP.NET Core trims the "Async" suffix
    // by default, which would 404 those posts.
    options.SuppressAsyncSuffixInActionNames = false;
})
.AddRazorRuntimeCompilation();

// ─── Authentication ───────────────────────────────────────────────────────────
var authService = BookstoreConfiguration.GetSetting("Services/Authentication");
if (authService == "aws")
{
    builder.Services.AddAuthentication(options =>
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
        options.UsePkce = true;
        options.Scope.Add("openid");
        options.Scope.Add("profile");
        options.CallbackPath = HttpRequestExtensions.OidcCallbackPath;
        options.SignInScheme = CookieAuthenticationDefaults.AuthenticationScheme;
        options.UseTokenLifetime = false;
        options.SaveTokens = true;
        options.TokenValidationParameters = new TokenValidationParameters
        {
            NameClaimType = "cognito:username",
            RoleClaimType = "cognito:groups"
        };
        options.Events = new OpenIdConnectEvents
        {
            // The OWIN implementation set the redirect URI on BOTH the authorize
            // request and the token request; Cognito rejects the code exchange if
            // the two do not match exactly.
            OnRedirectToIdentityProvider = context =>
            {
                context.ProtocolMessage.RedirectUri = context.Request.GetReturnUrl();

                return Task.CompletedTask;
            },
            OnAuthorizationCodeReceived = context =>
            {
                var returnUrl = context.Request.GetReturnUrl();

                context.TokenEndpointRequest.RedirectUri = returnUrl;
                context.Properties.RedirectUri ??= returnUrl;

                return Task.CompletedTask;
            },
            OnTokenValidated = async context =>
            {
                var service = context.HttpContext.RequestServices.GetRequiredService<ICustomerService>();

                if (context.Principal?.Identity is ClaimsIdentity identity)
                {
                    // Make the validated principal available for the remainder of
                    // this request, as the OWIN handler did via Request.User.
                    context.HttpContext.User = context.Principal;

                    var dto = new CreateOrUpdateCustomerDto(
                        identity.GetSub(),
                        identity.Name ?? string.Empty,
                        identity.FindFirst(c => c.Type.Contains("givenname"))?.Value ?? string.Empty,
                        identity.FindFirst(c => c.Type.Contains("surname"))?.Value ?? string.Empty);

                    await service.CreateOrUpdateCustomerAsync(dto);
                }
            }
        };
    });
}
else
{
    // Local authentication: cookie-based with custom middleware
    builder.Services.AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
        .AddCookie(CookieAuthenticationDefaults.AuthenticationScheme, options =>
        {
            options.LoginPath = "/Authentication/Login";
        });
}

// ─── Dependency Injection ─────────────────────────────────────────────────────
builder.Services.AddScoped<IBookService, BookService>();
builder.Services.AddScoped<IOrderService, OrderService>();
builder.Services.AddScoped<IReferenceDataService, ReferenceDataService>();
builder.Services.AddScoped<IOfferService, OfferService>();
builder.Services.AddScoped<ICustomerService, CustomerService>();
builder.Services.AddScoped<IAddressService, AddressService>();
builder.Services.AddScoped<IShoppingCartService, ShoppingCartService>();
builder.Services.AddScoped<IImageResizeService, ImageResizeService>();

builder.Services.AddScoped<ICustomerRepository, CustomerRepository>();
builder.Services.AddScoped<IAddressRepository, AddressRepository>();
builder.Services.AddScoped<IBookRepository, BookRepository>();
builder.Services.AddScoped<IOfferRepository, OfferRepository>();
builder.Services.AddScoped<IShoppingCartRepository, ShoppingCartRepository>();
builder.Services.AddScoped<IOrderRepository, OrderRepository>();
builder.Services.AddScoped<IReferenceDataRepository, ReferenceDataRepository>();
builder.Services.AddScoped(typeof(IPaginatedList<>), typeof(PaginatedList<>));

if (BookstoreConfiguration.GetSetting("Services/FileService") == "aws")
{
    builder.Services.AddSingleton<IAmazonS3, AmazonS3Client>();
    builder.Services.AddScoped<IFileService, S3FileService>();
}
else
{
    // Uploads must land under the served static-file root so that the URL
    // LocalFileService returns ("/images/coverimages/...") actually resolves.
    builder.Services.AddSingleton<IFileService>(sp =>
    {
        var env = sp.GetRequiredService<IWebHostEnvironment>();
        return new LocalFileService(env.WebRootPath);
    });
}

if (BookstoreConfiguration.GetSetting("Services/ImageValidationService") == "aws")
{
    builder.Services.AddSingleton<IAmazonRekognition, AmazonRekognitionClient>();
    builder.Services.AddScoped<IImageValidationService, RekognitionImageValidationService>();
}
else
{
    builder.Services.AddScoped<IImageValidationService, LocalImageValidationService>();
}

if (authService != "aws")
{
    builder.Services.AddScoped<LocalAuthenticationMiddleware>();
}

// HTTP context accessor for middleware/services
builder.Services.AddHttpContextAccessor();

// ─── Build app ────────────────────────────────────────────────────────────────
var app = builder.Build();

// ─── Ensure writable static content directories exist ────────────────────────
// LocalFileService writes cover images here and the files are served from the
// same path; the directory does not exist in a fresh container.
Directory.CreateDirectory(Path.Combine(app.Environment.WebRootPath, "images", "coverimages"));

// ─── Seed database ────────────────────────────────────────────────────────────
using (var scope = app.Services.CreateScope())
{
    var dbContext = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
    try
    {
        await BookstoreDbSeeder.SeedAsync(dbContext);
    }
    catch (Exception ex)
    {
        var logger = scope.ServiceProvider.GetRequiredService<ILogger<Program>>();
        logger.LogError(ex, "An error occurred seeding the database.");
    }
}

// ─── Middleware pipeline ──────────────────────────────────────────────────────
// Forwarded headers must run first so every later component (OIDC redirect
// building, cookie Secure policy, URL generation) sees the external scheme.
app.UseForwardedHeaders();

if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Home/Error");
}

app.UseStaticFiles();

app.UseRouting();

if (authService != "aws")
{
    // Populates HttpContext.User for local development; must run before
    // UseAuthorization so [Authorize]/[Authorize(Roles = ...)] sees the user.
    app.UseMiddleware<LocalAuthenticationMiddleware>();
}

app.UseAuthentication();
app.UseAuthorization();

// Area routing (Admin)
app.MapControllerRoute(
    name: "Admin_default",
    pattern: "Admin/{controller=Dashboard}/{action=Index}/{id?}",
    defaults: new { area = "Admin" },
    constraints: new { },
    dataTokens: new { area = "Admin" });

app.MapControllerRoute(
    name: "default",
    pattern: "{controller=Home}/{action=Index}/{id?}");

app.Run();

// ─── Helpers ──────────────────────────────────────────────────────────────────
static void ConfigureAwsLoggingTarget()
{
    if (BookstoreConfiguration.GetSetting("Services/LoggingService") != "aws") return;

    var config = LogManager.Configuration ?? new LoggingConfiguration();

    var awsTarget = new AWSTarget { Name = "aws", LogGroup = Constants.AppName };

    config.AddTarget(awsTarget);
    config.LoggingRules.Add(new LoggingRule("*", NLog.LogLevel.Info, awsTarget));

    LogManager.Configuration = config;
}

static void LoadAwsConfiguration(Logger logger)
{
    var rootPath = "/" + Constants.AppName;
    const string databasePath = "/Database";
    const string authenticationPath = "/Authentication";
    const string fileServicePath = "/Files";

    if (BookstoreConfiguration.GetSetting("Services/Database") == "aws")
    {
        try
        {
            using var client = new AmazonSimpleSystemsManagementClient();
            var request = new GetParameterRequest { Name = $"{rootPath}{databasePath}/ConnectionStrings/BookstoreDatabaseConnection" };
            var response = client.GetParameterAsync(request).GetAwaiter().GetResult();
            BookstoreConfiguration.AddSetting(response.Parameter.Name.Replace($"{rootPath}{databasePath}/", string.Empty), response.Parameter.Value);
        }
        catch (Exception ex)
        {
            // The MVC5 implementation let this propagate and fail startup. Swallowing it
            // would leave the app running against the local placeholder connection string.
            logger.Error(ex, "Unable to load the database configuration from SSM Parameter Store path {0}{1}.", rootPath, databasePath);
            throw;
        }
    }

    if (BookstoreConfiguration.GetSetting("Services/Authentication") == "aws")
    {
        try
        {
            using var client = new AmazonSimpleSystemsManagementClient();
            var request = new GetParametersByPathRequest { Path = $"{rootPath}{authenticationPath}/", Recursive = true };
            var response = client.GetParametersByPathAsync(request).GetAwaiter().GetResult();
            foreach (var parameter in response.Parameters)
            {
                BookstoreConfiguration.AddSetting(parameter.Name.Replace($"{rootPath}/", string.Empty), parameter.Value);
            }
        }
        catch (Exception ex)
        {
            logger.Error(ex, "Unable to load the authentication configuration from SSM Parameter Store path {0}{1}.", rootPath, authenticationPath);
            throw;
        }
    }

    if (BookstoreConfiguration.GetSetting("Services/FileService") == "aws")
    {
        try
        {
            using var client = new AmazonSimpleSystemsManagementClient();
            var request = new GetParametersByPathRequest { Path = $"{rootPath}{fileServicePath}/", Recursive = true };
            var response = client.GetParametersByPathAsync(request).GetAwaiter().GetResult();
            foreach (var parameter in response.Parameters)
            {
                BookstoreConfiguration.AddSetting(parameter.Name.Replace($"{rootPath}/", string.Empty), parameter.Value);
            }
        }
        catch (Exception ex)
        {
            logger.Error(ex, "Unable to load the file service configuration from SSM Parameter Store path {0}{1}.", rootPath, fileServicePath);
            throw;
        }
    }
}
