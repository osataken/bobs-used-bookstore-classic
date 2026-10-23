using System.IO;
using Autofac;
using Autofac.Extensions.DependencyInjection;
using BobsBookstoreClassic.Data;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace Bookstore.Web
{
    public class Program
    {
        public static async System.Threading.Tasks.Task Main(string[] args)
        {
            var builder = WebApplication.CreateBuilder(args);

            // MUST be the first call: BookstoreConfiguration is the application's settings
            // singleton (previously backed by Web.config/ConfigurationManager) and every step
            // below reads settings through it.
            BookstoreConfiguration.Initialize(builder.Configuration);

            // Same order as the OWIN Startup.Configuration: logging, configuration, DI, auth.
            LoggingSetup.ConfigureLogging(builder);

            await ConfigurationSetup.ConfigureConfigurationAsync();

            // Core's form tag helper / Html.BeginForm emit an anti-forgery token on every POST
            // form automatically; AutoValidateAntiforgeryToken validates it on every unsafe verb
            // (the MVC5 equivalent of decorating each POST action with [ValidateAntiForgeryToken]).
            builder.Services.AddControllersWithViews(options =>
            {
                options.Filters.Add(new AutoValidateAntiforgeryTokenAttribute());
            });
            builder.Services.AddHttpContextAccessor();
            builder.Services.AddScoped<ExceptionLoggingMiddleware>();

            // Replaces the global AuthorizeAttribute registered by App_Start/FilterConfig.cs.
            // A fallback policy (rather than an unconditional filter) keeps [AllowAnonymous]
            // controllers and actions publicly reachable, exactly as in MVC5.
            builder.Services.AddAuthorization(options =>
            {
                options.FallbackPolicy = new AuthorizationPolicyBuilder()
                    .RequireAuthenticatedUser()
                    .Build();
            });

            // Single Autofac container, created by the host's service provider factory.
            builder.Host.UseServiceProviderFactory(new AutofacServiceProviderFactory());
            builder.Host.ConfigureContainer<ContainerBuilder>(containerBuilder =>
            {
                containerBuilder.RegisterModule(new ApplicationModule(builder.Environment.WebRootPath));
            });

            // Cookie / OpenIdConnect (Cognito) or local authentication services.
            AuthenticationSetup.ConfigureAuthentication(builder);

            var app = builder.Build();

            // Cover images saved by LocalFileService used to be created lazily under Content;
            // ensure the wwwroot folder exists so the files are served as static content.
            Directory.CreateDirectory(Path.Combine(app.Environment.WebRootPath, "images", "coverimages"));

            if (app.Environment.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                app.UseExceptionHandler("/Home/Error");
                app.UseHsts();
            }

            // Logs unhandled exceptions the way Global.asax Application_Error did.
            app.UseMiddleware<ExceptionLoggingMiddleware>();

            app.UseHttpsRedirection();
            app.UseStaticFiles();

            app.UseRouting();

            app.UseAuthentication();

            // Local (non-Cognito) authentication middleware, in the same pipeline position as
            // the OWIN version: after the authentication handlers, before authorization.
            AuthenticationSetup.UseLocalAuthentication(app);

            app.UseAuthorization();

            // Area route first — Admin/... must win over the default route.
            app.MapControllerRoute(
                name: "Admin_default",
                pattern: "Admin/{controller}/{action=Index}/{id?}",
                defaults: new { area = "Admin" });

            app.MapControllerRoute(
                name: "Default",
                pattern: "{controller=Home}/{action=Index}/{id?}");

            await app.RunAsync();
        }
    }
}
