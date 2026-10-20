# Migration Summary: .NET Framework 4.8 → .NET 10

## Status: ✅ Complete — 0 errors, 0 warnings

All five projects build successfully targeting `net10.0`.

---

## Projects Migrated

| Project | Before | After | Notes |
|---------|--------|-------|-------|
| Bookstore.Cdk | net10.0 (SDK-style, already modern) | net10.0 | Bumped Amazon.CDK.Lib → 2.273.0, Constructs → 10.5.0 to resolve NU1901 vulnerability advisory |
| Bookstore.Common | netstandard2.0 | netstandard2.0 | No changes needed; simple constants project |
| Bookstore.Domain | netstandard2.0 | netstandard2.0 | Fixed duplicate `AssemblyInfo.cs` attributes |
| Bookstore.Data | netstandard2.0 + EntityFramework 6.5.1 | net10.0 + EF Core 9.0.4 | Full EF6 → EF Core migration (see below) |
| Bookstore.Web | .NET Framework 4.8 (old-style csproj, MVC 5 / OWIN) | net10.0 (SDK.Web, ASP.NET Core MVC) | Full stack migration (see below) |

---

## Bookstore.Data Changes

### EF6 → EF Core
- Replaced `EntityFramework 6.5.1` with `Microsoft.EntityFrameworkCore 9.0.4` + `Microsoft.EntityFrameworkCore.SqlServer 9.0.4`
- Migrated `ApplicationDbContext`:
  - Constructor changed from `(string connectionString) : base(connectionString)` → `(DbContextOptions<ApplicationDbContext> options) : base(options)`
  - `OnModelCreating(DbModelBuilder)` → `OnModelCreating(ModelBuilder)`
  - Removed `PluralizingTableNameConvention` (EF Core doesn't pluralize by default)
  - `HasRequired(...).WillCascadeOnDelete(false)` → `HasOne(...).OnDelete(DeleteBehavior.NoAction)`
  - `HasDatabaseGeneratedOption(DatabaseGeneratedOption.Identity)` → `ValueGeneratedOnAdd()`
  - Removed `Database.SetInitializer(new BookstoreDbInitializer())`
  - Moved seed data from `BookstoreDbInitializer` into `OnModelCreating` using `HasData()` for ReferenceDataItems and Books
- Migrated `PaginatedList<T>`: `using System.Data.Entity` → `using Microsoft.EntityFrameworkCore`
- Migrated all repositories:
  - `using System.Data.Entity` → `using Microsoft.EntityFrameworkCore`
  - Nested EF6 includes `.Include(x => x.OrderItems.Select(y => y.Book))` → `.Include(x => x.OrderItems).ThenInclude(y => y.Book)`
  - `await Task.Run(() => dbContext.X.Add(obj))` → `await dbContext.X.AddAsync(obj)`
  - String-based `Include("Genre")` → typed `Include(x => x.Genre)`
  - Removed `using Amazon.Auth.AccessControlPolicy` from `OfferRepository` (incorrect import)
- `BookstoreConfiguration`: Replaced `System.Configuration.ConfigurationManager` with an `Initialize(IConfiguration config)` static method that populates the singleton from ASP.NET Core `IConfiguration`. The public API (`GetSetting`, `GetConnectionString`, etc.) is unchanged, so no call sites required modification.
- `BookstoreDbInitializer`: Converted from `DropCreateDatabaseIfModelChanges<T>` to a static `InitializeAsync(ApplicationDbContext)` helper that calls `EnsureCreatedAsync()`.

### Target Framework Decision
Bookstore.Data was collapsed from `netstandard2.0` to `net10.0` only. EF Core 9.0.4 targets `net8.0` only (no `netstandard2.0` support), making multi-targeting fundamentally incompatible. This project is exclusively consumed by the `net10.0` Bookstore.Web application.

---

## Bookstore.Web Changes

### Project File
- Replaced old-style `.csproj` (`.NET Framework 4.8`) with SDK-style `Microsoft.NET.Sdk.Web` targeting `net10.0`
- Removed all old-style package references (EntityFramework 6, Autofac.Mvc5, Autofac.Owin, Microsoft.Owin.*, System.Web.Mvc, etc.)
- Added ASP.NET Core packages: `Microsoft.AspNetCore.Authentication.OpenIdConnect`, `Autofac.Extensions.DependencyInjection`, `NLog.Web.AspNetCore`

### Startup / Application Lifecycle
- `Global.asax.cs` → excluded from compilation; replaced by `Program.cs`
- `Startup.cs` (OWIN) → excluded from compilation; replaced by `Program.cs`
- `App_Start/` files → converted to empty stubs; all logic moved to `Program.cs`
- Created `Program.cs` using `WebApplication.CreateBuilder` with:
  - `BookstoreConfiguration.Initialize(builder.Configuration)`
  - SSM Parameter Store loading via synchronous `GetAwaiter().GetResult()` at startup
  - `builder.Host.UseServiceProviderFactory(new AutofacServiceProviderFactory())`
  - `builder.Host.ConfigureContainer<ContainerBuilder>(ConfigureDependencyInjection)`
  - Cookie + OpenID Connect authentication for AWS Cognito; cookie-only for local
  - NLog via `builder.Host.UseNLog()`
  - `app.UseMiddleware<LocalAuthenticationMiddleware>()` for non-AWS auth
- Created `appsettings.json` with all settings from `web.config`

### Authentication
- OWIN middleware (`Microsoft.Owin.Security.Cookies`, `Microsoft.Owin.Security.OpenIdConnect`) → `Microsoft.AspNetCore.Authentication.Cookies` + `Microsoft.AspNetCore.Authentication.OpenIdConnect`
- `LocalAuthenticationMiddleware`: `OwinMiddleware` → conventional ASP.NET Core middleware (`RequestDelegate` + `InvokeAsync(HttpContext context, ICustomerService service)`)
- Authentication cookie now handled via `Response.Cookies.Delete()` / `Response.Cookies.Append()`

### MVC / Controllers
- All controllers: `using System.Web.Mvc` → `using Microsoft.AspNetCore.Mvc`
- `ActionResult` return types changed to `IActionResult`
- `[RouteArea("Admin")]` → `[Area("Admin")]`
- `AreaRegistration` class converted to empty comment stub (areas auto-discovered in ASP.NET Core)
- `model.CoverImage?.InputStream` → `model.CoverImage?.OpenReadStream()`

### Models
- `HttpPostedFileBase CoverImage` → `IFormFile CoverImage` in `InventoryCreateUpdateViewModel`
- `using System.Web.Mvc` → `using Microsoft.AspNetCore.Mvc.Rendering` for `SelectListItem` in all view models

### Helpers
- `HttpContextExtensions`: `HttpContextBase` → `HttpContext` (ASP.NET Core)
- `HttpCookie` → `CookieOptions`; `context.Response.Cookies.Add(new HttpCookie(...))` → `context.Response.Cookies.Append(...)`
- `MvcHelpers`: `HtmlHelper` → `IHtmlHelper` extension
- `MaxFileSizeAttribute`: `HttpPostedFileBase` → `IFormFile`; `file.ContentLength` → `file.Length`
- `ImageTypesAttribute`: `HttpPostedFileBase` → `IFormFile`
- `IOwinRequestExtensions`: `IOwinRequest` → `HttpRequest`

### Views
- Created `Views/_ViewImports.cshtml` with correct ASP.NET Core usings (`Microsoft.AspNetCore.Mvc.Rendering`, `Microsoft.AspNetCore.Routing`)
- Created `Areas/Admin/Views/_ViewImports.cshtml`
- Replaced `@{Html.RenderPartial("...")}` → `@{await Html.RenderPartialAsync("...")}` (RenderPartial removed in ASP.NET Core)
- Replaced `@Html.EnumDropDownListFor(...)` (MVC5-only, not in ASP.NET Core) → `Html.GetEnumSelectList<TEnum>()` with inline `foreach`
- `RouteValueDictionary` in views now resolves via `using Microsoft.AspNetCore.Routing` in `_ViewImports.cshtml`

### Packages Updated / Removed
- Removed: `Autofac.Mvc5`, `Autofac.Owin`, `Microsoft.Owin.*`, `Owin`, `System.Web.Mvc`, `EntityFramework 6`, `WebGrease`, `Antlr3.Runtime`, `Microsoft.AspNet.Web.Optimization`, `Microsoft.CodeDom.Providers.DotNetCompilerPlatform`
- Upgraded: `Magick.NET-Q8-AnyCPU` 14.6.0 → 15.0.0 (resolves ~100 NU190x vulnerability advisories)
- Updated: `NLog.Web.AspNetCore` 5.3.16 → 5.4.0 (resolves NU1603)

---

## Next Steps

- **EF Core Migrations**: The current setup uses `EnsureCreatedAsync()` which creates the schema if absent but doesn't manage migrations. For production use, generate an initial EF Core migration with `dotnet ef migrations add InitialCreate` and use `context.Database.MigrateAsync()` instead.
- **HTTPS**: For production deployments, add `app.UseHttpsRedirection()` and configure HTTPS certificates.
- **Antiforgery**: The forms currently use `<form method="post">` without explicit antiforgery tokens. Tag helper forms (`<form asp-action="...">`) emit tokens automatically, but hand-written `<form>` elements (e.g. in admin views) should be audited and `@Html.AntiForgeryToken()` added.
- **Bookstore.Data netstandard2.0**: Removed. This project is consumed only by the net10.0 web application. If future callers on .NET Framework need access, a separate compatibility shim would be needed.
- **LocalFileService webRootPath**: In Autofac DI registration, `LocalFileService` is registered with `webRootPath = "wwwroot"` (relative path). For production, inject `IWebHostEnvironment` and use `env.WebRootPath` for an absolute path.
- **Amazon.CDK.Lib `CloudFrontWebDistribution` obsolete**: `CoreStack.cs` uses the deprecated `CloudFrontWebDistribution` class. Migrate to `Distribution` in a future pass.
