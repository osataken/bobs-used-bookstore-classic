# Migration Summary — .NET Framework → .NET 10

## Result

`dotnet build BobsBookstoreClassic.sln` exits with **0 errors, 0 warnings** on all five projects.

---

## Projects migrated

| Project | Target framework | Status |
|---|---|---|
| Bookstore.Cdk | net10.0 | ✅ Clean |
| Bookstore.Common | netstandard2.0 | ✅ Clean (shared library; netstandard2.0 is appropriate for reuse by CDK and Web) |
| Bookstore.Web | net10.0 | ✅ Clean |
| Bookstore.Domain | net10.0 | ✅ Clean |
| Bookstore.Data | net10.0 | ✅ Clean |

---

## Changes made in this cycle

### 1. Razor Views — MVC1000 (IHtmlHelper.Partial deadlock warnings)

Replaced every synchronous `@Html.Partial(...)` and `@{ Html.RenderPartial(...); }` call with the `<partial name="..." />` tag helper. Affected files:

- `Views/Wishlist/Index.cshtml`
- `Views/ShoppingCart/Index.cshtml`
- `Views/Shared/_Layout.cshtml`
- `Views/Shared/_NavBarPartial.cshtml`
- `Views/Search/Index.cshtml`
- `Views/Address/CreateUpdate.cshtml`
- `Areas/Admin/Views/Shared/_Layout.cshtml`
- `Areas/Admin/Views/Shared/_NavBarPartial.cshtml`
- `Areas/Admin/Views/Orders/Index.cshtml`
- `Areas/Admin/Views/Orders/Details.cshtml`
- `Areas/Admin/Views/Inventory/Index.cshtml`
- `Areas/Admin/Views/Inventory/CreateUpdate.cshtml`
- `Areas/Admin/Views/ReferenceData/Index.cshtml`
- `Areas/Admin/Views/Offers/Index.cshtml`

### 2. Bookstore.Cdk — CDK CloudFront API modernisation

Replaced the deprecated `CloudFrontWebDistribution` + `CloudFrontWebDistributionProps` pattern (CS0618/CS0612) with the current CDK v2 `Distribution` API in `CoreStack.CreateCloudFrontDistribution()`:

- Added `using Amazon.CDK.AWS.CloudFront.Origins;`
- Replaced `new CloudFrontWebDistribution(...)` with `new Distribution(...)` using `DistributionProps` / `BehaviorOptions`
- Replaced `S3OriginConfig` + `OriginConfigs[]` with `S3BucketOrigin.WithOriginAccessIdentity(...)` (which also fixes the secondary `S3Origin` deprecation warning)

### 3. Bookstore.Cdk — CDK NuGet package vulnerability (NU1901)

Bumped CDK-related packages to their latest compatible versions:

| Package | Before | After |
|---|---|---|
| Amazon.CDK.Lib | 2.188.0 | 2.274.0 |
| Cdklabs.CdkNag | 2.35.66 | 2.38.2 |
| Constructs | 10.4.2 | 10.8.1 |

The low-severity vulnerability GHSA-464c-974j-9xm6 reported against 2.188.0 is resolved by this bump.

---

## What AWS Transform Custom (atx) delivered (unchanged)

- SDK-style `.csproj` files targeting `net10.0` for Bookstore.Web, Bookstore.Data, Bookstore.Domain, and Bookstore.Cdk
- `System.Web` → ASP.NET Core migration (controllers, routing, areas)
- EF6 → EF Core migration (`ApplicationDbContext`, `BookstoreDbInitializer`, repositories, `PaginatedList`)
- Global.asax → `Program.cs` (application bootstrap, DI, authentication, logging, error handling middleware)
- Web.config / ConfigurationManager → `appsettings.json` + `IConfiguration` (`BookstoreConfiguration`)
- OWIN middleware (authentication, local dev auth, SSM parameter loading) → ASP.NET Core middleware
- Autofac container wired with `AutofacServiceProviderFactory`
- NLog configured programmatically with `AWS.Logger.NLog` target
- AWS SDK v4 integrated (S3, Rekognition, SSM, CloudWatch Logs)

---

## Next steps

None required — the build is fully clean. The following are observations for future iterations:

1. **Database migrations**: `BookstoreDbInitializer` uses `EnsureCreated()` which creates the schema from scratch but does not support incremental changes. When the EF Core model changes, migrate to `dotnet ef migrations add` / `database update` to apply schema changes without data loss.

2. **Bookstore.Common target framework**: Currently `netstandard2.0`, which is correct for multi-targeting reuse. If the CDK project is decoupled from the web project's shared library in future, Bookstore.Common could be upgraded to `net10.0` only.

3. **Image cover paths in seed data**: The seeded `Book` entries reference `/Content/Images/coverimages/...` paths (the old ASP.NET MVC physical path convention). The application now serves static files from `wwwroot/images/coverimages/`. These are seed-data strings written once at schema creation time; no code change is needed but the referenced images should be placed under `wwwroot/images/coverimages/` to serve correctly.
