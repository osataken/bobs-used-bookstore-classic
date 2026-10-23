using Autofac;
using Autofac.Extensions.DependencyInjection;
using BobsBookstoreClassic.Data;
using Bookstore.Data;
using Bookstore.Web;
using Bookstore.Web.Helpers;
using Bookstore.Web.Middleware;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.Authorization;
using NLog.Web;

var builder = WebApplication.CreateBuilder(args);

// Must run before anything reads a setting: BookstoreConfiguration no longer self-initializes
// from ConfigurationManager, the host supplies the IConfiguration (task 13).
BookstoreConfiguration.Initialize(builder.Configuration);

// NLog programmatic configuration (AWSTarget vs DebuggerTarget), then NLog behind
// Microsoft.Extensions.Logging.
LoggingSetup.ConfigureLogging();

builder.Logging.ClearProviders();
builder.Host.UseNLog();

// AWS Systems Manager Parameter Store overrides. Asynchronous because the modern AWS SDK only
// offers async SSM operations; awaited here, before the host is built, so no setting is read
// before it is loaded.
await ConfigurationSetup.ConfigureConfigurationAsync();

// Autofac composition root (single container).
builder.Host.UseServiceProviderFactory(new AutofacServiceProviderFactory());
builder.Host.ConfigureContainer<ContainerBuilder>(containerBuilder =>
    DependencyInjectionSetup.ConfigureDependencyInjection(containerBuilder, builder.Environment));

// Replaces FilterConfig's global AuthorizeAttribute: every endpoint requires an authenticated
// user unless it opts out with [AllowAnonymous].
builder.Services.AddControllersWithViews(options =>
{
    options.Filters.Add(new AuthorizeFilter());

    // MVC5 routed to the literal method name, so actions named "ApproveAsync"/"RejectAsync"/
    // "ReceivedAsync"/"PaidAsync" (Areas/Admin OffersController) were reachable at those URLs.
    // ASP.NET Core trims the "Async" suffix by default, which would 404 the hardcoded
    // formaction URLs in Areas/Admin/Views/Offers/Index.cshtml. Keep the MVC5 naming.
    options.SuppressAsyncSuffixInActionNames = false;

    // Replaces the per-action [ValidateAntiForgeryToken] convention of the MVC5 app:
    // antiforgery validation for every unsafe (POST/PUT/PATCH/DELETE) request.
    options.Filters.Add(new AutoValidateAntiforgeryTokenAttribute());
});

builder.Services.AddHttpContextAccessor();

// Port of App_Start/AuthenticationSetup.cs: Cognito OpenID Connect when
// Services/Authentication == "aws", otherwise cookie-only for the local-auth middleware path.
AuthenticationSetup.ConfigureAuthentication(builder.Services);

builder.Services.AddAuthorization();

// Load-balancer health endpoint (mapped anonymous below, since a global AuthorizeFilter
// would otherwise answer the ALB with a 302 login challenge).
builder.Services.AddHealthChecks();

var app = builder.Build();

// Replaces the EF6 Database.SetInitializer hook, which self-triggered on first use. EF Core has
// no initializer concept, so the schema creation and reference-data seeding must be invoked
// explicitly by the host (task 23).
using (var scope = app.Services.CreateScope())
{
    var dbContext = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();

    BookstoreDbInitializer.Initialize(dbContext);
}

if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}
else
{
    app.UseExceptionHandler("/Home/Error");
}

// Logs unhandled exceptions (ported Application_Error) inside the exception handler, so the
// error view is still rendered.
app.UseMiddleware<ExceptionLoggingMiddleware>();

app.UseStaticFiles();

app.UseRouting();

app.UseAuthentication();

// Local (non-Cognito) authentication: mirrors the original ConfigureLocalAuthentication branch.
if (!AuthenticationSetup.UseCognitoAuthentication)
{
    app.UseMiddleware<LocalAuthenticationMiddleware>();
}

app.UseAuthorization();

// ECS/ALB target-group health check target; anonymous so the global authorize filter
// does not turn it into a 302.
app.MapHealthChecks("/health").AllowAnonymous();

// Named/area routes are registered before the default catch-all.
app.MapControllerRoute(
    name: "areas",
    pattern: "{area:exists}/{controller=Dashboard}/{action=Index}/{id?}");

app.MapControllerRoute(
    name: "Default",
    pattern: "{controller=Home}/{action=Index}/{id?}");

app.Run();
