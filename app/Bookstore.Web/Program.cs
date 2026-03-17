using Amazon.Rekognition;
using Amazon.S3;
using Amazon.SimpleSystemsManagement;
using Amazon.SimpleSystemsManagement.Model;
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
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authentication.OpenIdConnect;
using Microsoft.AspNetCore.Builder;
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
using NLog.Targets;
using NLog.Web;
using System.Collections.Generic;
using System.IO;
using System.Security.Claims;
using WebOptimizer;

var builder = WebApplication.CreateBuilder(args);
var configuration = builder.Configuration;

LoadAwsParameters(configuration);
ConfigureLogging(configuration);

builder.Services.AddControllersWithViews();

var connectionString = configuration.GetConnectionString("BookstoreDatabaseConnection");
if (configuration["Services:Database"] == "aws")
{
    builder.Services.AddDbContext<ApplicationDbContext>(options =>
        options.UseSqlServer(connectionString));
}
else
{
    builder.Services.AddDbContext<ApplicationDbContext>(options =>
        options.UseSqlite(connectionString));
}

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

if (configuration["Services:FileService"] == "aws")
{
    builder.Services.AddScoped<IAmazonS3, AmazonS3Client>();
    builder.Services.AddScoped<IFileService, S3FileService>();
}
else
{
    var webRootPath = builder.Environment.WebRootPath ?? Path.Combine(builder.Environment.ContentRootPath, "wwwroot");
    builder.Services.AddSingleton<IFileService>(new LocalFileService(webRootPath));
}

if (configuration["Services:ImageValidationService"] == "aws")
{
    builder.Services.AddScoped<IAmazonRekognition, AmazonRekognitionClient>();
    builder.Services.AddScoped<IImageValidationService, RekognitionImageValidationService>();
}
else
{
    builder.Services.AddScoped<IImageValidationService, LocalImageValidationService>();
}

if (configuration["Services:Authentication"] == "aws")
{
    builder.Services.AddAuthentication(options =>
    {
        options.DefaultScheme = CookieAuthenticationDefaults.AuthenticationScheme;
        options.DefaultChallengeScheme = OpenIdConnectDefaults.AuthenticationScheme;
    })
    .AddCookie()
    .AddOpenIdConnect(options =>
    {
        options.ClientId = configuration["Authentication:Cognito:LocalClientId"];
        options.MetadataAddress = configuration["Authentication:Cognito:MetadataAddress"];
        options.ResponseType = OpenIdConnectResponseType.Code;
        options.Scope.Add("openid");
        options.Scope.Add("profile");
        options.SaveTokens = true;
        options.TokenValidationParameters = new TokenValidationParameters
        {
            NameClaimType = "cognito:username",
            RoleClaimType = "cognito:groups"
        };
        options.Events = new OpenIdConnectEvents
        {
            OnTokenValidated = async context =>
            {
                var service = context.HttpContext.RequestServices.GetRequiredService<ICustomerService>();
                var identity = (ClaimsIdentity)context.Principal.Identity;

                var dto = new CreateOrUpdateCustomerDto(
                    identity.FindFirst("sub").Value,
                    identity.Name,
                    identity.FindFirst("given_name")?.Value,
                    identity.FindFirst("family_name")?.Value);

                await service.CreateOrUpdateCustomerAsync(dto);
            }
        };
    });
}
else
{
    builder.Services.AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
        .AddCookie(options =>
        {
            options.LoginPath = "/Authentication/Login";
            options.LogoutPath = "/Authentication/Logout";
        });
}

builder.Services.AddWebOptimizer(pipeline =>
{
    pipeline.AddCssBundle("/Content/css", "Content/css/site.css", "Content/css/styles.css", "Content/css/custom-style.css");
});

builder.Logging.ClearProviders();
builder.Services.AddHealthChecks();

builder.Host.UseNLog();

var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    var dbContext = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
    BookstoreDbInitializer.Initialize(dbContext);
}

if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Home/Error");
    app.UseHsts();
}

app.UseHttpsRedirection();

app.Use(async (context, next) =>
{
    context.Response.Headers["X-Content-Type-Options"] = "nosniff";
    context.Response.Headers["X-Frame-Options"] = "DENY";
    context.Response.Headers["Referrer-Policy"] = "strict-origin-when-cross-origin";
    await next();
});

app.UseStaticFiles();
app.UseWebOptimizer();

app.UseRouting();

app.UseAuthentication();

if (configuration["Services:Authentication"] != "aws")
{
    app.UseMiddleware<Bookstore.Web.Helpers.LocalAuthenticationMiddleware>();
}

app.UseAuthorization();

app.MapHealthChecks("/health");

app.MapControllerRoute(
    name: "areas",
    pattern: "{area:exists}/{controller=Dashboard}/{action=Index}/{id?}");

app.MapControllerRoute(
    name: "default",
    pattern: "{controller=Home}/{action=Index}/{id?}");

app.Run();

void ConfigureLogging(IConfiguration config)
{
    var loggingConfig = new LoggingConfiguration();

    NLog.Targets.Target loggingTarget;

    if (config["Services:LoggingService"] == "aws")
    {
        loggingTarget = new AWSTarget { LogGroup = Constants.AppName };
    }
    else
    {
        loggingTarget = new DebuggerTarget();
    }

    loggingConfig.AddTarget("aws", loggingTarget);
    loggingConfig.LoggingRules.Add(new LoggingRule("*", NLog.LogLevel.Info, loggingTarget));

    LogManager.Configuration = loggingConfig;
}

void LoadAwsParameters(IConfiguration config)
{
    var rootPath = "/" + Constants.AppName;

    const string databasePath = "/Database";
    const string authenticationPath = "/Authentication";
    const string fileServicePath = "/Files";

    if (config["Services:Database"] == "aws")
    {
        using (var client = new AmazonSimpleSystemsManagementClient())
        {
            var request = new GetParameterRequest { Name = $"{rootPath}{databasePath}/ConnectionStrings/BookstoreDatabaseConnection" };
            var response = client.GetParameterAsync(request).GetAwaiter().GetResult();

            config[response.Parameter.Name.Replace($"{rootPath}{databasePath}/", string.Empty).Replace("/", ":")] = response.Parameter.Value;
        }
    }

    if (config["Services:Authentication"] == "aws")
    {
        using (var client = new AmazonSimpleSystemsManagementClient())
        {
            var request = new GetParametersByPathRequest { Path = $"{rootPath}{authenticationPath}/", Recursive = true };
            var response = client.GetParametersByPathAsync(request).GetAwaiter().GetResult();

            foreach (var parameter in response.Parameters)
            {
                config[parameter.Name.Replace($"{rootPath}/", string.Empty).Replace("/", ":")] = parameter.Value;
            }
        }
    }

    if (config["Services:FileService"] == "aws")
    {
        using (var client = new AmazonSimpleSystemsManagementClient())
        {
            var request = new GetParametersByPathRequest { Path = $"{rootPath}{fileServicePath}/", Recursive = true };
            var response = client.GetParametersByPathAsync(request).GetAwaiter().GetResult();

            foreach (var parameter in response.Parameters)
            {
                config[parameter.Name.Replace($"{rootPath}/", string.Empty).Replace("/", ":")] = parameter.Value;
            }
        }
    }
}
