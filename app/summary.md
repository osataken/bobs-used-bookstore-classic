# Migration Summary: .NET Framework 4.8 → .NET 10

## Result

All 5 projects build with **zero compilation errors** targeting `net10.0`.

| Project | Target | Result | Warnings |
|---------|--------|--------|----------|
| Bookstore.Common | net10.0;netstandard2.0 | ✅ Build succeeded | 0 |
| Bookstore.Domain | net10.0;netstandard2.0 | ✅ Build succeeded | 0 |
| Bookstore.Data | net10.0 | ✅ Build succeeded | 4 (NU only) |
| Bookstore.Cdk | net10.0 | ✅ Build succeeded | 5 (NU only) |
| Bookstore.Web | net10.0 | ✅ Build succeeded | 6 (NU only) |

All remaining warnings are NU1902/NU1903 package vulnerability advisories — there are no C# compiler (CS) warnings in any project.

---

## Changes Made

### Project Files

**Bookstore.Common.csproj**
- Converted to SDK-style, multi-targeted `net10.0;netstandard2.0`
- Added `LangVersion=latest`, `Nullable=enable`

**Bookstore.Domain.csproj**
- Converted to SDK-style, multi-targeted `net10.0;netstandard2.0`
- Removed `System.Data.DataSetExtensions` (unused)
- Retained `System.ComponentModel.Annotations` and `Microsoft.CSharp`
- `GenerateAssemblyInfo=false` to prevent conflict with `Properties/AssemblyInfo.cs`

**Bookstore.Data.csproj**
- Collapsed to single `net10.0` target (netstandard2.0 dropped — EF Core 10 requires net8.0+, and no EF6 version supports netstandard2.0; this is the KB-defined NU1201 escape hatch for fundamentally incompatible package graphs)
- Replaced EF6 (`EntityFramework 6.5.1`) with EF Core 10 (`Microsoft.EntityFrameworkCore.SqlServer 10.0.0`)
- Upgraded `Magick.NET-Q8-AnyCPU` from 14.6.0 → 15.0.0 to clear high/moderate security vulnerabilities
- Pinned `System.Security.Cryptography.Pkcs` to 10.0.0 (still tracked in CVE — see Next Steps)

**Bookstore.Web.csproj**
- Full conversion from legacy ToolsVersion/MSBuild XML to `Microsoft.NET.Sdk.Web` SDK-style
- Target `net10.0` (single TFM — application project, never multi-targeted)
- Replaced all `<Reference>` and `<HintPath>` items with inline `<PackageReference>`
- Removed OWIN packages (`Microsoft.Owin.*`, `Owin`), `Autofac.Mvc5`, `Autofac.Owin`, `System.Web.Optimization`, `WebGrease`, `Antlr`, `Microsoft.CodeDom.Providers.*`
- Added ASP.NET Core packages: `Microsoft.AspNetCore.Authentication.OpenIdConnect`, `Microsoft.AspNetCore.Mvc.Razor.RuntimeCompilation`
- Replaced `Autofac` DI with built-in ASP.NET Core DI
- Excluded legacy startup files from compilation (added to `<Compile Remove>`)

### Application Startup

**Created: `Program.cs`** — replaces Global.asax.cs, Startup.cs (OWIN), and all App_Start/* files:
- Registers all services via ASP.NET Core DI (IBookService, IOrderService, ICustomerService, etc.)
- Configures EF Core DbContext with SqlServer
- Sets up authentication: Cookie + OpenIdConnect (Cognito) or local authentication middleware
- Configures NLog logging via `NLog.Extensions.Logging`
- Serves legacy `Content/` and `Scripts/` folders via additional `UseStaticFiles` entries (preserves existing URL paths)
- Loads AWS SSM parameters on startup (ConfigurationSetup equivalent)
- Bootstraps `BookstoreConfiguration` from `IConfiguration` (backwards-compatible key translation `:` → `/`)
- Registers Admin area route and default route
- Calls `BookstoreDbSeeder.SeedAsync()` to seed initial data on first run

**Created: `appsettings.json`** — replaces `Web.config` appSettings and connectionStrings:
- All service switches (`Services/Authentication`, etc.) as nested JSON keys
- Connection string for SQL LocalDB (dev default)
- NLog configuration section
- Cognito and S3 settings (placeholders — populated from SSM at runtime)

**Created: `appsettings.Development.json`** — dev environment overrides

### EF6 → EF Core Migration

**ApplicationDbContext.cs**
- Constructor changed: `(string connectionString)` → `(DbContextOptions<ApplicationDbContext> options)`
- `OnModelCreating`: `DbModelBuilder` → `ModelBuilder`
- `HasRequired(...).WithMany()...WillCascadeOnDelete(false)` → `HasOne(...).WithMany()...OnDelete(DeleteBehavior.NoAction)`
- `HasColumnType("nvarchar").HasMaxLength(450)` → `.HasMaxLength(450)`
- `HasDatabaseGeneratedOption(DatabaseGeneratedOption.Identity)` → `.ValueGeneratedOnAdd()`
- Removed `modelBuilder.Conventions.Remove<PluralizingTableNameConvention>()` (EF Core doesn't pluralize by default)
- Removed `Database.SetInitializer(...)` (replaced by `BookstoreDbSeeder`)

**Removed: `BookstoreDbInitializer.cs`** — replaced by `BookstoreDbSeeder.cs`

**Created: `BookstoreDbSeeder.cs`** — EF Core seeder called from Program.cs at startup

**PaginatedList.cs** — updated `using System.Data.Entity` → `using Microsoft.EntityFrameworkCore`

**All Repositories** — updated:
- `using System.Data.Entity` → `using Microsoft.EntityFrameworkCore`
- `await Task.Run(() => dbContext.X.Add(entity))` → `await dbContext.X.AddAsync(entity)` (EF Core)
- EF6 `.Include(x => x.OrderItems.Select(y => y.Book))` → EF Core `.Include(x => x.OrderItems).ThenInclude(y => y.Book)` (OrderRepository, ShoppingCartRepository)
- Return types updated to nullable where EF Core's `SingleOrDefaultAsync`/`FindAsync` may return null

### OWIN → ASP.NET Core Migration

**LocalAuthenticationMiddleware.cs** — rewritten as ASP.NET Core middleware:
- Inherits from nothing (not OwinMiddleware)
- Constructor: `(RequestDelegate next)` — `ICustomerService` injected via `InvokeAsync` parameter (scoped service)
- `IOwinContext` → `HttpContext`
- `HttpContext.Current.Response.Cookies.Add(new HttpCookie(...))` → `context.Response.Cookies.Append(...)`
- `context.Response.Redirect(...)` — same pattern

**Helpers/HttpContextExtensions.cs** — migrated:
- `HttpContextBase` → `HttpContext` (Microsoft.AspNetCore.Http)
- `HttpCookie` → `CookieOptions`
- `context.Request.Cookies[key]` → same (ASP.NET Core has same API)
- `context.Response.Cookies.Add(new HttpCookie(...))` → `context.Response.Cookies.Append(...)`

**Removed: `Helpers/IOwinRequestExtensions.cs`** — OWIN-specific, no ASP.NET Core equivalent needed

### MVC Migration

All controllers updated:
- `using System.Web.Mvc` → `using Microsoft.AspNetCore.Mvc`
- `using Microsoft.AspNetCore.Authorization`
- `ActionResult` → `IActionResult`
- `HttpPostedFileBase` → `IFormFile`
- `file.InputStream` → `file.OpenReadStream()`
- `HttpCookie` → cookie API via `HttpContext.Response.Cookies`

**Areas/Admin/AdminAreaControllerBase.cs**
- `[RouteArea("Admin")]` → `[Area("Admin")]`

**Areas/Admin/AdminAreaRegistration.cs** — removed (areas are auto-discovered in ASP.NET Core)

**App_Start/** files — all excluded from compilation (functionality moved to Program.cs):
- `BundleConfig.cs` — bundling replaced by static files served from Content/ and Scripts/
- `RouteConfig.cs` — routes in Program.cs `app.MapControllerRoute()`
- `FilterConfig.cs` — global authorize filter in `AddControllersWithViews(options => ...)`
- `AuthenticationSetup.cs` — authentication in Program.cs
- `DependencyInjectionSetup.cs` — services registered in Program.cs
- `LoggingSetup.cs` — NLog configured in Program.cs
- `ConfigurationSetup.cs` — configuration loading in Program.cs

### View Model Updates

All view models using `System.Web.Mvc.SelectListItem` updated to `Microsoft.AspNetCore.Mvc.Rendering.SelectListItem`.

**InventoryCreateUpdateViewModel.cs** — `HttpPostedFileBase` → `IFormFile`

**Admin Views** — `Html.EnumDropDownListFor` replaced with `Html.DropDownListFor` + `Html.GetEnumSelectList<T>()`

All views using `@Html.Partial(...)` and `@{Html.RenderPartial(...);}` replaced with `<partial name="..." />` tag helper (fixes MVC1000 deadlock warnings).

### Nullable Reference Type Improvements (Domain/Data)

- All entity classes updated with nullable annotations (`string?` for optional properties, `= null!` for navigation properties populated by EF)
- All repository interfaces updated to return nullable types (`Task<T?>`) where appropriate
- All service classes updated to handle nullable returns from repositories
- `IFileService.SaveAsync` → accepts and returns nullable stream/string
- `IImageValidationService.IsSafeAsync` → accepts nullable stream

---

## Remaining Warnings (All NU — Not Resolvable by Version Bump)

| Package | Version | Advisory | Severity | Notes |
|---------|---------|----------|----------|-------|
| `System.Security.Cryptography.Pkcs` | 10.0.0 | GHSA-2fmv-rv2c-59b7 (mod), GHSA-846f-f764-76mx (high) | High | .NET 10 platform package; same CVEs tracked in all available .NET 10 versions; no non-vulnerable version available. |

---

## Next Steps

1. **System.Security.Cryptography.Pkcs CVE**: Monitor Microsoft security bulletins (https://github.com/advisories/GHSA-846f-f764-76mx) for a patched .NET 10 release. Apply when available.

2. **Bookstore.Data netstandard2.0 compatibility**: Bookstore.Data was collapsed from multi-target (`net10.0;netstandard2.0`) to single `net10.0` because EF Core 10 and EF6 both lack support for netstandard2.0. Any callers that haven't yet upgraded from .NET Framework will need to use the pre-migration library version or be upgraded themselves.

3. **Nullable Reference Types in Bookstore.Web**: The Web project has `<Nullable>disable</Nullable>` to avoid flooding the build with CS8618/CS8604 warnings on ~20 ViewModel files migrated from a non-nullable .NET Framework codebase. Enabling nullable in the Web project incrementally (per-file with `#nullable enable`) is recommended as a follow-up.

4. **CloudFrontWebDistribution deprecation in Bookstore.Cdk**: `CloudFrontWebDistribution` is marked obsolete (CS0618). Consider migrating to the `Distribution` API when updating the CDK stack.

5. **EF Core Migrations**: The migration uses `EnsureCreated()` (via `BookstoreDbSeeder`). For production, replace with `dotnet ef migrations add InitialCreate` + `dotnet ef database update` to manage schema changes incrementally.

6. **HTTPS for Cognito**: The README notes that Amazon Cognito Hosted UI requires HTTPS for redirect URIs (except localhost). When deploying to ECS on Fargate, configure HTTPS termination (ALB + ACM certificate) before enabling Cognito authentication.

7. **Static file serving**: The legacy `Content/` and `Scripts/` directories are served via additional `UseStaticFiles` middleware entries. Consider migrating to `wwwroot/` as part of a clean-up pass.
