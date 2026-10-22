# Migration Summary: .NET Framework 4.8 → .NET 10 Upgrade

## Status: ✅ Complete — All projects build with zero compilation errors

### Projects Migrated

| Project | From | To | Status |
|---|---|---|---|
| Bookstore.Domain | .NET Framework 4.8 (SDK-style, net10.0 by ATX) | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Common | netstandard2.0 | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Data | .NET Framework 4.8 (old-style csproj) | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Cdk | net10.0 (already by ATX) | net10.0 | ✅ 0 errors, 5 warnings (NU1901 only) |
| Bookstore.Web | .NET Framework 4.8 (old-style csproj) | net10.0 | ✅ 0 errors, 0 warnings |

---

## Key Changes Made

### Bookstore.Common
- Changed target from `netstandard2.0` to `net10.0`

### Bookstore.Data
- Converted from old-style `.csproj` to SDK-style targeting `net10.0`
- **EF6 → EF Core 9**: Replaced `System.Data.Entity` with `Microsoft.EntityFrameworkCore`
  - `ApplicationDbContext`: Rewrote constructor to use `DbContextOptions<ApplicationDbContext>`, updated `OnModelCreating` fluent API (`HasRequired→HasOne`, `WillCascadeOnDelete→OnDelete(NoAction)`, `HasDatabaseGeneratedOption→ValueGeneratedOnAdd`)
  - `BookstoreDbInitializer`: Changed from `DropCreateDatabaseIfModelChanges<T>` to a static `Seed(ApplicationDbContext)` method compatible with EF Core's `EnsureCreated()` seeding pattern
  - `PaginatedList<T>`: Replaced `using System.Data.Entity` with `using Microsoft.EntityFrameworkCore` for async extension methods
  - All Repositories: Updated `using System.Data.Entity` → `using Microsoft.EntityFrameworkCore`; converted string-based `.Include("Navigation")` to lambda `.Include(x => x.Navigation)`; converted nested EF6 `.Include(x => x.Items.Select(y => y.Child))` chains to EF Core `.Include(x => x.Items).ThenInclude(y => y.Child)` chains; changed `Task.Run(() => dbContext.DbSet.Add(entity))` to `await dbContext.DbSet.AddAsync(entity)`
- **BookstoreConfiguration**: Rewrote from `System.Configuration.ConfigurationManager` to `IConfiguration` pattern with a new `Initialize(IConfiguration)` static method; kept namespace `BobsBookstoreClassic.Data` and all existing API surface intact
- **Magick.NET upgraded**: From 14.6.0 to 15.0.0 to resolve vulnerability advisories (NU1901–NU1903)

### Bookstore.Web
- Converted from old-style `.csproj` to `Microsoft.NET.Sdk.Web` SDK-style targeting `net10.0`
- Excluded legacy files from compilation (retained on disk for reference): `Startup.cs`, `Global.asax.cs`, `App_Start/RouteConfig.cs`, `App_Start/BundleConfig.cs`, `App_Start/FilterConfig.cs`, `App_Start/DependencyInjectionSetup.cs`, `App_Start/AuthenticationSetup.cs`, `Areas/Admin/AdminAreaRegistration.cs`, `Properties/AssemblyInfo.cs`
- **All MVC Controllers**: Replaced `System.Web.Mvc` with `Microsoft.AspNetCore.Mvc`; `ActionResult` → `IActionResult`
- **AuthenticationController**: Replaced `HttpCookie`/manual cookie manipulation with `HttpContext.SignOutAsync(CookieAuthenticationDefaults.AuthenticationScheme)` for local sign-out; Cognito sign-out now uses `Request.Scheme`/`Request.Host`
- **LocalAuthenticationMiddleware**: Rewrote from OWIN `OwinMiddleware` to ASP.NET Core `IMiddleware`; cookie-based sign-in via `context.SignInAsync(CookieAuthenticationDefaults.AuthenticationScheme, ...)`
- **HttpContextExtensions**: Replaced `HttpContextBase` (old MVC) with `Microsoft.AspNetCore.Http.HttpContext`; cookie API updated to `Response.Cookies.Append()`
- **IOwinRequestExtensions**: Replaced `IOwinRequest`-based extension with `HttpRequest`-based `GetReturnUrl()`
- **MvcHelpers**: Replaced `HtmlHelper` with `IHtmlHelper` from `Microsoft.AspNetCore.Mvc.Rendering`
- **MaxFileSizeAttribute / ImageTypesAttribute**: Replaced `HttpPostedFileBase` with `IFormFile`; `ContentLength` → `Length`; `InputStream` → `OpenReadStream()`
- **ControllerExtensions**: Updated to `Microsoft.AspNetCore.Mvc.Controller`
- **All models using SelectListItem**: Replaced `System.Web.Mvc.SelectListItem` with `Microsoft.AspNetCore.Mvc.Rendering.SelectListItem`
- **InventoryCreateUpdateViewModel**: Replaced `HttpPostedFileBase CoverImage` with `IFormFile CoverImage`
- **ConfigurationSetup**: Converted synchronous AWS SDK calls (`GetParameter`, `GetParametersByPath`) to async (`GetParameterAsync`, `GetParametersByPathAsync`); updated `Program.cs` to `await ConfigurationSetup.ConfigureConfigurationAsync()`
- **LoggingSetup**: Updated to use `AWS.Logger.NLog` package (was incorrectly referencing `NLog.AWS.Logger`); resolved `LogLevel` ambiguity between `NLog.LogLevel` and `Microsoft.Extensions.Logging.LogLevel`
- **Program.cs**: Added `using Bookstore.Web;` for `ConfigurationSetup`/`LoggingSetup` resolution; updated to `await ConfigureConfigurationAsync()`
- **Views — Html.Partial/RenderPartial**: All `@Html.Partial(...)` and `Html.RenderPartial(...)` calls converted to `<partial name="..." />` tag helper syntax (eliminated MVC1000 warnings)
- **Views — EnumDropDownListFor**: Replaced `@Html.EnumDropDownListFor(...)` (not available in ASP.NET Core) with `@Html.DropDownListFor(..., Html.GetEnumSelectList<T>(), ...)`
- **Admin Views — RouteValueDictionary**: Added `Areas/Admin/Views/_ViewImports.cshtml` with `@using Microsoft.AspNetCore.Routing` to resolve `RouteValueDictionary` type in Dashboard view
- **Views/_ViewImports.cshtml**: Already present with `@addTagHelper *, Microsoft.AspNetCore.Mvc.TagHelpers`

---

## Next Steps

1. **Amazon.CDK.Lib 2.188.0 — NU1901 (low severity)**: The CDK project has 5 remaining NU1901 advisory warnings for `Amazon.CDK.Lib 2.188.0`. Upgrading CDK to a newer version is a separate infrastructure concern and could change CDK construct APIs; this should be handled as a dedicated CDK version upgrade task.

2. **Runtime validation**: The migration targets SQL Server via EF Core. Ensure the connection string in `appsettings.json` is valid for the target environment before deploying. The `EnsureCreated()` + `Seed()` pattern in Program.cs will create and seed the database on first run.

3. **Authentication (Local mode)**: The `LocalAuthenticationMiddleware` now uses ASP.NET Core cookie authentication via `SignInAsync`. Verify that the `/Authentication/Login` path correctly triggers sign-in in the dev environment.

4. **Cognito HTTPS redirect**: As noted in README.md, the Cognito Hosted UI requires HTTPS redirects (except localhost). The ECS deployment uses HTTP, so AWS authentication remains disabled for ECS-deployed instances.

5. **Dockerfile**: The existing Dockerfile targets Windows containers (`FROM mcr.microsoft.com/dotnet/framework/aspnet`). Since the app now targets net10.0 and runs on Linux, the Dockerfile should be updated to use a Linux-based image (e.g., `mcr.microsoft.com/dotnet/aspnet:10.0`). See `19-linux-containerization.md` knowledge base article.
