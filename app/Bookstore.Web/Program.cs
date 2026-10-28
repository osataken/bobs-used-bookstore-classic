using System;
using System.Threading.Tasks;
using Autofac;
using Autofac.Extensions.DependencyInjection;
using BobsBookstoreClassic.Data;
using Bookstore.Data;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using NLog.Web;

namespace Bookstore.Web
{
    public static class Program
    {
        public static async Task Main(string[] args)
        {
            var builder = WebApplication.CreateBuilder(args);

            // 1. Configuration facade first: everything below reads BookstoreConfiguration.
            BookstoreConfiguration.Initialize(builder.Configuration);

            // 2. NLog configuration before the first log write, then route all ASP.NET Core
            //    logging through NLog (replaces the Global.asax/OWIN Startup logging bootstrap).
            LoggingSetup.ConfigureLogging();

            builder.Logging.ClearProviders();
            builder.Host.UseNLog();

            // 3. AWS Systems Manager parameters must be loaded before anything composes services
            //    or authentication from those settings. Awaited here, so it is complete before
            //    the registrations below run.
            await ConfigurationSetup.ConfigureConfigurationAsync();

            // 4. Dependency injection (Autofac). The container registrations live in
            //    Startup/DependencyInjectionSetup.cs.
            builder.Host.UseServiceProviderFactory(new AutofacServiceProviderFactory());
            builder.Host.ConfigureContainer<ContainerBuilder>(DependencyInjectionSetup.Configure);

            builder.Services.AddHttpContextAccessor();

            builder.Services.AddDbContext<ApplicationDbContext>(options =>
                options.UseSqlServer(BookstoreConfiguration.GetConnectionString("BookstoreDatabaseConnection")));

            // SuppressAsyncSuffixInActionNames is disabled deliberately: views post to literal
            // /Admin/Offers/ApproveAsync|RejectAsync|ReceivedAsync|PaidAsync/{id} URLs, which the
            // default suffix trimming would turn into 404s.
            builder.Services.AddControllersWithViews(options =>
            {
                options.SuppressAsyncSuffixInActionNames = false;
            });

            // Equivalent of the MVC5 global AuthorizeAttribute: every endpoint requires an
            // authenticated user unless it opts out with [AllowAnonymous].
            builder.Services.AddAuthorization(options =>
            {
                options.FallbackPolicy = new AuthorizationPolicyBuilder()
                    .RequireAuthenticatedUser()
                    .Build();
            });

            // Authentication (Cognito/OIDC or local) lives in Startup/AuthenticationSetup.cs.
            AuthenticationSetup.Configure(builder);

            builder.Services.AddHealthChecks();

            var app = builder.Build();

            // 5. Database creation/seeding, once, on a scoped context (the EF Core replacement for
            //    the EF6 Database.SetInitializer registration, which armed itself lazily on first
            //    use). Retried with backoff: in ECS the database may not accept connections yet,
            //    and the EF6 behaviour was to fail the first request, not the whole process.
            using (var scope = app.Services.CreateScope())
            {
                var context = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
                var logger = scope.ServiceProvider.GetRequiredService<ILogger<WebApplication>>();

                const int maxAttempts = 5;

                for (var attempt = 1; attempt <= maxAttempts; attempt++)
                {
                    try
                    {
                        BookstoreDbInitializer.Initialize(context);

                        break;
                    }
                    catch (Exception ex) when (attempt < maxAttempts)
                    {
                        logger.LogWarning(ex, "Database initialization attempt {Attempt} of {MaxAttempts} failed; retrying.", attempt, maxAttempts);

                        await Task.Delay(TimeSpan.FromSeconds(5));
                    }
                    catch (Exception ex)
                    {
                        // Do not crash the host: the ALB health check must be able to answer so the
                        // task is not killed before the database becomes reachable.
                        logger.LogError(ex, "Database initialization failed after {MaxAttempts} attempts; continuing startup.", maxAttempts);
                    }
                }
            }

            // 6. Middleware pipeline.
            // Admin-area exceptions render the area-local error page under the Admin layout, as
            // the MVC5 HandleErrorAttribute did by resolving Error.cshtml through the area view
            // hierarchy. Must be registered before the site-wide handler.
            app.UseWhen(
                context => context.Request.Path.StartsWithSegments("/Admin"),
                branch => branch.UseExceptionHandler("/Error/Index"));

            app.UseExceptionHandler("/Home/Error");
            app.UseStatusCodePages();

            // Equivalent of Global.asax Application_Error: log every unhandled exception before
            // the error view is rendered by the exception handler registered above.
            app.Use(async (context, next) =>
            {
                try
                {
                    await next(context);
                }
                catch (Exception ex)
                {
                    var logger = context.RequestServices.GetRequiredService<ILogger<WebApplication>>();

                    logger.LogError(ex, "Unhandled exception processing {Method} {Path}", context.Request.Method, context.Request.Path);

                    throw;
                }
            });

            app.UseStaticFiles();

            app.UseRouting();

            app.UseAuthentication();

            // Local development identity (OWIN UseMiddlewareFromContainer equivalent). Only when
            // Cognito is not configured; see Startup/AuthenticationSetup.cs.
            if (!AuthenticationSetup.IsAwsAuthentication())
            {
                app.UseMiddleware<Helpers.LocalAuthenticationMiddleware>();
            }

            app.UseAuthorization();

            // Area route first, so Admin/* resolves before the default route.
            app.MapControllerRoute("admin_area", "{area:exists}/{controller=Dashboard}/{action=Index}/{id?}");
            app.MapControllerRoute("default", "{controller=Home}/{action=Index}/{id?}");

            app.MapHealthChecks("/health").AllowAnonymous();

            await app.RunAsync();
        }
    }
}
