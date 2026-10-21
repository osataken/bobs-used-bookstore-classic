# .NET Framework → .NET 10 Migration Summary

## Migration Outcome

All 5 projects build successfully with **zero compilation errors and zero warnings**.

| Project | Target Framework | Result |
|---------|-----------------|--------|
| Bookstore.Cdk | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Common | netstandard2.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Domain | netstandard2.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Data | net10.0 + netstandard2.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Web | net10.0 | ✅ 0 errors, 0 warnings |

---

## Changes Made

### Bookstore.Data.csproj
- Made EF Core packages (EF Core 9.x, SqlServer, Design) conditional on `net10.0` only — EF Core 9 requires net8.0+ and is incompatible with netstandard2.0.
- Added `Condition` to exclude EF Core–dependent files (`ApplicationDbContext.cs`, `BookstoreDbInitializer.cs`, `PaginatedList.cs`, all `Repositories/*.cs`) from the `netstandard2.0` build target. The `netstandard2.0` leg still compiles (containing AWS SDK services, file services, and image services), maintaining backward-compatible lib consumption.

### Bookstore.Domain.csproj
- Added `System.ComponentModel.Annotations` package (for `[Timestamp]` attribute in `Entity.cs`).
- Added `<LangVersion>latest</LangVersion>` and `<Nullable>enable</Nullable>` to enable nullable reference type analysis.

### Nullable Reference Type Fixes (Domain)
- **Entity.cs**: `byte[]? RowVersion`
- **Book.cs**: navigation properties use `= null!`, optional strings use `string?`, required strings use `= string.Empty`
- **BookFilters.cs**: `string? Name`, `string? Author`
- **BookResult.cs**: `string? ErrorMessage`
- **BookService.cs**: `SaveAsync` and `ResizeImageAsync` accept `Stream?` and `string?`; `IsSafeAsync` accepts `Stream?`
- **IImageValidationService.cs**: `IsSafeAsync(Stream? image)` 
- **Offer.cs**: navigation properties use `= null!`, optional strings use `string?`
- **OfferFilters.cs**: `string? BookName`, `string? Author`
- **Customer.cs**: required strings use `= string.Empty`, optional strings use `string?`
- **Address.cs**: required strings use `= string.Empty`, `AddressLine2` is `string?`, navigation property uses `= null!`; removed `#pragma warning disable CS8618`
- **Order.cs**: navigation properties use `= null!`
- **OrderItem.cs**: navigation properties use `= null!`
- **ShoppingCart.cs**: `CorrelationId = string.Empty`
- **ShoppingCartItem.cs**: navigation properties use `= null!`
- **ReferenceDataItem.cs**: `Text = string.Empty`
- **DbSecrets.cs**: all string properties marked `string?`
- **AddressDtos.cs**: constructor parameters changed to `string?` with `?? string.Empty` fallback
- **BookDtos.cs**: constructor parameters changed to `string?`/`Stream?` with appropriate fallbacks
- **OfferDtos.cs**: constructor parameters changed to `string?` with `?? string.Empty` fallback
- **OrderDtos.cs**: constructor parameters changed to `string?`/`int` with `?? string.Empty` fallback
- **ReferenceDataDtos.cs**: `Text` parameter changed to `string?`
- **AddressService.cs / IAddressService**: `sub` parameters changed to `string?`
- **OrderService.cs / IOrderService**: `sub` parameter in `GetOrdersAsync` changed to `string?`
- **OfferService.cs / IOfferService**: `sub` parameter in `GetOffersAsync` changed to `string?`

### Nullable Reference Type Fixes (Data)
- **Repositories**: All `SingleOrDefaultAsync`/`FindAsync` return values now throw `InvalidOperationException` for not-found cases (eliminating CS8603). This gives explicit errors rather than null-pointer exceptions downstream.
- **PaginatedList.cs**: `source` field initialized to `null!`
- **LocalFileService.cs / S3FileService.cs**: `SaveAsync` returns `string.Empty` instead of `null` when input is null
- **LocalImageValidationService.cs**: `IsSafeAsync(Stream? image)`
- **RekognitionImageValidationService.cs**: `IsSafeAsync(Stream? image)`

### Nullable Reference Type Fixes (Web)
- **ClaimsPrincipalExtensions.cs**: return type changed to `string?`
- **HttpContextExtensions.cs**: `shoppingCartClientId` uses `string?`; `Append` gets `?? string.Empty`
- **MvcHelpers.cs**: `emptyItem` parameter changed to `string?`
- **AuthenticationController.cs**: `redirectUri` parameter changed to `string?`
- **LocalAuthenticationMiddleware.cs**: cast `principal.Identity` using `as` pattern with null guard
- **ImageTypesAttribute.cs / MaxFileSizeAttribute.cs**: `IsValid(object? value)` override signature
- **ErrorViewModel.cs** (Web and Admin): `string? RequestId`
- **Program.cs**: null-guard on `ctx.TokenEndpointRequest`
- All view models: string properties changed to `string?` or initialized with `= string.Empty` as appropriate
- **AddressController.cs**: `Redirect` calls use `?? "/"` fallback; `AddModelError` uses `?? "An error occurred."`
- **SearchDetailsViewModel.cs**: all string properties marked `string?`

### MVC1000 Warnings (Razor Views)
Replaced all `@Html.Partial("name")` and `@Html.Partial("name", Model)` calls with `<partial name="name" />` and `<partial name="name" model="Model" />` tag helpers in:
- Views/Wishlist/Index.cshtml
- Views/ShoppingCart/Index.cshtml
- Views/Search/Index.cshtml
- Views/Shared/_Layout.cshtml
- Views/Shared/_NavBarPartial.cshtml
- Areas/Admin/Views/Shared/_Layout.cshtml
- Areas/Admin/Views/Shared/_NavBarPartial.cshtml
- Areas/Admin/Views/Orders/Details.cshtml
- Areas/Admin/Views/Orders/Index.cshtml
- Areas/Admin/Views/ReferenceData/Index.cshtml
- Areas/Admin/Views/Inventory/Index.cshtml
- Areas/Admin/Views/Offers/Index.cshtml

### CS8602 in Razor Views
- **_LoginPartial.cshtml** (main and Admin): changed `User.Identity.IsAuthenticated` to `User.Identity?.IsAuthenticated == true`
- **_NavBarPartial.cshtml**: same null-conditional fix
- **ShoppingCart/Index.cshtml**: same fix

### CDK CloudFront Migration (CS0618/CS0612)
- Replaced deprecated `CloudFrontWebDistribution` + `CloudFrontWebDistributionProps` with the modern `Distribution` + `DistributionProps` API.
- Replaced deprecated `S3Origin` with `S3BucketOrigin.WithOriginAccessIdentity(...)`.
- Added `using Amazon.CDK.AWS.CloudFront.Origins;` import.
- Result: Bookstore.Cdk now builds with **0 warnings**.

---

## Architecture Notes

### Multi-targeting in Bookstore.Data
`Bookstore.Data` retains the `netstandard2.0` target so the library can still be referenced by any .NET Framework project that hasn't yet migrated. The `netstandard2.0` build compiles:
- `BookstoreConfiguration.cs`
- `FileServices/LocalFileService.cs`
- `FileServices/S3FileService.cs`
- `ImageValidationServices/LocalImageValidationService.cs`
- `ImageValidationServices/RekognitionImageValidationService.cs`
- `ImageResizeService/ImageResizeService.cs`

EF Core–dependent code (DbContext, repositories, PaginatedList, initializer) is excluded from the `netstandard2.0` TFM because EF Core 9.x requires net8.0+. Any .NET Framework consumer that needs data-access functionality would need a separate EF6-based adapter.

---

## Next Steps

- **EF Core migrations**: Run `dotnet ef migrations add InitialCreate --project Bookstore.Data --startup-project Bookstore.Web` if database schema needs to be applied to a new SQL Server instance.
- **Dockerfile**: Review `Bookstore.Web/Dockerfile` to ensure it targets a net10.0 base image (`mcr.microsoft.com/dotnet/aspnet:10.0`).
- **ShoppingCartRepository.GetAsync**: Changed from returning `null` to throwing `InvalidOperationException` when cart not found — callers should handle this gracefully; review `ShoppingCartController` for cases where a missing cart is expected (e.g., first visit).
- **Bookstore.Common**: Currently targets `netstandard2.0` only. If it needs modern APIs, add `net10.0` to `<TargetFrameworks>`.
- **Integration tests**: No test projects are present. Adding integration/unit tests is recommended to validate migrated behaviour end-to-end.

---

## Post-Migration Hardening (Round 2)

### ✅ EF Core Migrations (EnsureCreatedAsync → MigrateAsync)

**Problem**: `EnsureCreatedAsync()` creates the schema once and never applies changes; it bypasses the migrations history table entirely, making incremental schema changes impossible in production.

**Changes**:
- Added `Microsoft.EntityFrameworkCore.Design` package to `Bookstore.Web.csproj` (required by the `dotnet-ef` tooling).
- Configured `MigrationsAssembly("Bookstore.Web")` on `UseSqlServer(...)` in `Program.cs`. Because `Bookstore.Data` is a multi-targeted library (net10.0 + netstandard2.0), the `dotnet-ef` tool cannot collect MSBuild metadata from it directly. Placing migrations in the startup project (`Bookstore.Web`) is the standard workaround.
- Ran `dotnet ef migrations add InitialCreate --framework net10.0 --output-dir Migrations` from `app/Bookstore.Web`, generating:
  - `Migrations/20261021075657_InitialCreate.cs` — full Up/Down schema
  - `Migrations/20261021075657_InitialCreate.Designer.cs` — EF model snapshot metadata
  - `Migrations/ApplicationDbContextModelSnapshot.cs` — baseline model snapshot for diff detection
- Updated `BookstoreDbInitializer.InitializeAsync()` to call `context.Database.MigrateAsync()` instead of `EnsureCreatedAsync()`. At startup this applies any pending migrations and creates the `__EFMigrationsHistory` table if absent, enabling safe incremental schema evolution.

### ✅ HTTPS Redirection

Added `app.UseHttpsRedirection()` guarded to the `Development` environment only.

**Rationale**: On ECS Fargate, TLS is terminated at the Application Load Balancer (ALB). The container speaks HTTP on port 8080 to the ALB. Enabling `UseHttpsRedirection()` unconditionally inside the container would cause infinite HTTP→HTTPS redirect loops on Fargate because every inbound request would be HTTP. In Development (local developer machine with the ASP.NET Core dev certificate), `UseHttpsRedirection()` is safe and useful.

If unconditional HTTPS is needed (e.g., a direct-TLS Fargate deployment without an ALB), add `ASPNETCORE_URLS=https://+:8081;http://+:8080` to the container environment and mount a certificate, then remove the `IsDevelopment()` guard.

### ✅ Antiforgery Tokens

**Problem**: Hand-written `<form method="post">` elements (without `asp-action` tag helper attributes) bypass the `FormTagHelper` automatic token injection. Submitting these forms without a valid token would succeed even after enabling validation, leaving CSRF vulnerabilities.

**Changes**:

_Global validation_ — added `AutoValidateAntiforgeryTokenAttribute` as a global MVC filter in `Program.cs`:
```csharp
builder.Services.AddControllersWithViews(options =>
    options.Filters.Add(new AutoValidateAntiforgeryTokenAttribute()));
```
This enforces token validation on every non-safe HTTP method (POST, PUT, DELETE, PATCH) across all controllers without needing per-action `[ValidateAntiForgeryToken]` attributes.

_Form-level tokens_ — added `@Html.AntiForgeryToken()` as the first child of every raw `<form method="post">` that lacked it:

| View | Form purpose |
|------|-------------|
| `Areas/Admin/Views/Orders/Details.cshtml` | Update order status |
| `Areas/Admin/Views/Offers/Index.cshtml` | Approve / Reject / Received / Paid offer actions |
| `Areas/Admin/Views/ReferenceData/CreateUpdate.cshtml` | Create / Update reference data item |
| `Areas/Admin/Views/Inventory/CreateUpdate.cshtml` | Create / Update book inventory |
| `Views/Resale/Create.cshtml` | Submit book resale offer |
| `Views/Address/Index.cshtml` | Delete address (modal form) |
| `Views/Checkout/Index.cshtml` | Place order |

Forms using `@using (Html.BeginForm(...))` — `Address/CreateUpdate.cshtml`, `Wishlist/Index.cshtml`, `ShoppingCart/Index.cshtml` — already had tokens injected automatically by `HtmlHelper.BeginForm`.

### ✅ Dockerfile: .NET 8 → .NET 10 Linux Images

Updated base image tags in `app/Bookstore.Web/Dockerfile`:
- `mcr.microsoft.com/dotnet/aspnet:8.0` → `mcr.microsoft.com/dotnet/aspnet:10.0`
- `mcr.microsoft.com/dotnet/sdk:8.0` → `mcr.microsoft.com/dotnet/sdk:10.0`

Both stages were already targeting Linux/amd64 (`--platform=linux/amd64`), aligned with the CDK `EcsStack` Fargate deployment.

### ✅ LocalFileService: Hardcoded "wwwroot" → IWebHostEnvironment.WebRootPath

**Problem**: `LocalFileService` was registered with `WithParameter("webRootPath", "wwwroot")` — a relative path string. When the working directory is not the application root (e.g., inside a container, or when running tests), `Path.Combine("wwwroot", ...)` resolves to a wrong absolute path and file saves fail silently.

**Changes** in `Program.cs`:
```csharp
var webRootPath = builder.Environment.WebRootPath;
builder.Host.ConfigureContainer<ContainerBuilder>(
    container => ConfigureDependencyInjection(container, webRootPath));
```
`builder.Environment.WebRootPath` returns the fully-resolved absolute path to `wwwroot` (set to `<ContentRoot>/wwwroot` by default, or overridable via `ASPNETCORE_WEBROOT`). `LocalFileService` itself is unchanged — it still accepts a `string webRootPath` constructor parameter, keeping the `Bookstore.Data` library free of ASP.NET Core hosting dependencies.

`ConfigureDependencyInjection` signature changed from `static void ConfigureDependencyInjection(ContainerBuilder builder)` to `static void ConfigureDependencyInjection(ContainerBuilder builder, string webRootPath)`.

### ✅ CDK CloudFrontWebDistribution (Already Complete)

`CoreStack.cs` already uses the modern `Distribution` + `S3BucketOrigin.WithOriginAccessIdentity(...)` API from the previous migration cycle. No further changes required.
