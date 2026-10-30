using System.IO;
using Autofac;
using Autofac.Extensions.DependencyInjection;
using BobsBookstoreClassic.Data;
using Bookstore.Web.Helpers;
using Bookstore.Web.Infrastructure;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Authorization;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using NLog.Web;

var builder = WebApplication.CreateBuilder(args);

// MUST run before any GetSetting/GetConnectionString call and before the SSM load below.
BookstoreConfiguration.Initialize(builder.Configuration);

// Was ConfigurationSetup.ConfigureConfiguration() in the OWIN Startup class. Pulls the
// database/authentication/file-service settings from SSM Parameter Store when the matching
// Services/* setting is "aws".
await ConfigurationSetup.ConfigureConfigurationAsync();

// Was LoggingSetup.ConfigureLogging() in the OWIN Startup class.
LoggingSetup.ConfigureLogging();
builder.Logging.ClearProviders();
builder.Host.UseNLog();

// ---------------------------------------------------------------------------------------------
// Dependency injection: Autofac now backs the framework container instead of
// DependencyResolver.SetResolver / app.UseAutofacMiddleware.
// ---------------------------------------------------------------------------------------------
builder.Host.UseServiceProviderFactory(new AutofacServiceProviderFactory());

// Replaces HttpRuntime.AppDomainAppVirtualPath / AppDomainAppPath probing.
var webRootPath = builder.Environment.WebRootPath ?? Path.Combine(builder.Environment.ContentRootPath, "wwwroot");
Directory.CreateDirectory(Path.Combine(webRootPath, "images", "coverimages"));

builder.Services.AddControllersWithViews(options =>
    {
        // Was FilterConfig: global AuthorizeAttribute. Every action requires an authenticated
        // user unless it opts out with [AllowAnonymous].
        options.Filters.Add(new AuthorizeFilter());
    })
    .AddMvcOptions(options =>
    {
        // Keep the Admin Offers *Async action names in their URLs (Core trims "Async" by default).
        options.SuppressAsyncSuffixInActionNames = false;
    })
    // Replaces builder.RegisterControllers(typeof(MvcApplication).Assembly) so Autofac owns controllers.
    .AddControllersAsServices();

// Replaces Global.asax Application_Error (NLog logger.Error(ex)).
builder.Services.AddExceptionHandler<NLogExceptionHandler>();

// EF Core DbContext + AWS SDK clients.
DependencyInjectionSetup.ConfigureServices(builder.Services);

// ===== AUTH_SERVICE_REGISTRATION_HOOK =====
// Was AuthenticationConfig.ConfigureAuthentication(IAppBuilder) in App_Start/AuthenticationSetup.cs:
// cookie + Cognito OpenIdConnect when Services/Authentication == "aws", cookie-only otherwise
// (the local path is handled by LocalAuthenticationMiddleware below).
AuthenticationConfig.ConfigureAuthentication(builder.Services);
// ===== END AUTH_SERVICE_REGISTRATION_HOOK =====

// Services, repositories, paginated list, file/image-validation service branches, local auth middleware.
builder.Host.ConfigureContainer<ContainerBuilder>(containerBuilder =>
    DependencyInjectionSetup.ConfigureContainer(containerBuilder, webRootPath));

var app = builder.Build();

// Was the global HandleErrorAttribute. Web.config had debug="true" with the default
// customErrors mode="RemoteOnly", so local requests saw full exception detail and only remote
// requests got Error.cshtml. NLogExceptionHandler preserves Application_Error logging.
if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}
else
{
    app.UseExceptionHandler("/Home/Error");
    app.UseStatusCodePagesWithReExecute("/Error/Index/{0}");
}

app.UseStaticFiles();

app.UseRouting();

// ===== AUTH_MIDDLEWARE_HOOK =====
// Was app.UseMiddlewareFromContainer<LocalAuthenticationMiddleware>() in the local branch of
// AuthenticationConfig. IMiddleware-based, so Autofac supplies the scoped ICustomerService.
//
// This MUST run BEFORE UseAuthentication/UseAuthorization. [Authorize(Roles="Administrators")]
// on AdminAreaControllerBase is endpoint metadata, so AuthorizationMiddleware evaluates it at
// UseAuthorization(); if the synthetic local principal were assigned after that, the whole Admin
// area would be denied for an anonymous user. UseAuthentication does not overwrite an already
// assigned context.User (it only assigns when AuthenticateAsync yields a principal), so the
// synthetic identity survives.
if (!AuthenticationConfig.UseCognitoAuthentication)
{
    app.UseMiddleware<LocalAuthenticationMiddleware>();
}
// ===== END AUTH_MIDDLEWARE_HOOK =====

app.UseAuthentication();
app.UseAuthorization();

// Was AdminAreaRegistration + RouteConfig. IgnoreRoute("{resource}.axd/...") is dropped.
app.MapControllerRoute(
    name: "areas",
    pattern: "{area:exists}/{controller}/{action}/{id?}");

app.MapControllerRoute(
    name: "Default",
    pattern: "{controller=Home}/{action=Index}/{id?}");

app.Run();
