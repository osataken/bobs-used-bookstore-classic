# BobsBookstoreClassic — .NET Framework → net10.0 Migration Summary

## Final Build Status

```
dotnet build BobsBookstoreClassic.sln  →  Build succeeded, 0 Error(s), 0 Warning(s)
```

All five projects pass individually and as a solution:

| Project | Target | Status |
|---|---|---|
| Bookstore.Common | netstandard2.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Domain | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Data | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Cdk | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Web | net10.0 | ✅ 0 errors, 0 warnings |

---

## What Was Done in This Cycle

### Phase 1 — AWS Transform Custom CLI (ATX)

ATX performed the bulk of the migration automatically:

- **Project files**: All converted to SDK-style; Bookstore.Data, Bookstore.Domain, Bookstore.Web, Bookstore.Cdk retargeted to `net10.0`. Bookstore.Common remains on `netstandard2.0` (correct — it has no framework-specific code).
- **EF6 → EF Core**: `EntityFramework` 6.5.1 removed; `Microsoft.EntityFrameworkCore` / `.Relational` / `.SqlServer` / `.Design` added. `ApplicationDbContext` rewritten with EF Core model builder, explicit table mappings, FK configurations, and `ShoppingCartItem` composite key. Explicit `HasForeignKey(x => x.CustomerId)` on `Order → Customer` added during the verify phase to prevent a shadow FK column.
- **OWIN/MVC5 → ASP.NET Core**: `Global.asax`, `App_Start/*`, `FilterConfig`, `BundleConfig`, `RouteConfig`, `AdminAreaRegistration` all replaced by a single `Program.cs`. `AuthorizeFilter` applied globally; `SuppressAsyncSuffixInActionNames = false` preserves Admin Offers action URLs.
- **Authentication**: `Microsoft.Owin.Security.OpenIdConnect` replaced by `Microsoft.AspNetCore.Authentication.OpenIdConnect 10.0.0`. Cookie-only local branch preserved via `LocalAuthenticationMiddleware`. `AuthenticationConfig.ConfigureAuthentication()` wires either Cognito OIDC or local cookies depending on `Services/Authentication`.
- **Dependency Injection**: `Autofac.Mvc5` / `Autofac.Owin` replaced by `Autofac.Extensions.DependencyInjection` + `AutofacServiceProviderFactory`. All original per-request registrations preserved in `DependencyInjectionSetup`.
- **Configuration**: `Web.config` / `ConfigurationManager` replaced by `appsettings.json` + `IConfiguration` via `BookstoreConfiguration.Initialize()`. SSM Parameter Store loading preserved in `ConfigurationSetup`.
- **Logging**: NLog + `AWS.Logger.NLog` replaces the OWIN/CloudWatch logging setup. `NLogExceptionHandler` replaces `Global.asax Application_Error`.
- **Static assets**: 320 content files moved into `wwwroot/`. `Scripts.Render` / `Styles.Render` bundle calls removed from views; Bootstrap 5 CDN links and real `wwwroot` paths used instead.
- **Dockerfile / CDK**: Dockerfile rewritten for `aspnet:10.0` / `sdk:10.0`, port 8080, non-root `$APP_UID`. `EcsStack` updated. CDK packages: `Amazon.CDK.Lib` 2.188.0 → 2.275.0, `Constructs` 10.4.2 → 10.5.0, `Amazon.Jsii.Analyzers` pinned to 1.142.0.
- **CVE remediation**: `Magick.NET-Q8-AnyCPU` 14.6.0 → 14.15.0 (162 known vulnerabilities, highest severity HIGH).

### Phase 2 — Kiro CLI post-ATX fixes

The build ATX left behind was already clean. Kiro applied targeted improvements:

1. **EF Core aligned to net10.0**: Upgraded `Microsoft.EntityFrameworkCore`, `.Relational`, `.SqlServer`, `.Design`, and `Microsoft.Extensions.Configuration.Abstractions` from 9.0.0 → **10.0.14** to match the target framework. Using EF Core 9 on net10.0 was functional but inconsistent.

2. **ASP.NET Core patch aligned**: `Microsoft.AspNetCore.Authentication.OpenIdConnect` 10.0.0 → **10.0.14** (patch update).

3. **Newtonsoft.Json patch**: 13.0.3 → **13.0.4**.

4. **CVE: System.Security.Cryptography.Pkcs** (NU1902/NU1903): Upgrading to EF Core 10.0.14 surfaced a transitive `System.Security.Cryptography.Pkcs 9.0.11` with two known CVEs (GHSA-2fmv-rv2c-59x7 moderate, GHSA-846f-f764-76mx high). Fixed by pinning `System.Security.Cryptography.Pkcs` to **10.0.14** directly in both `Bookstore.Data.csproj` and `Bookstore.Web.csproj`, overriding the transitive 9.x resolution.

---

## Final Package Versions

### Bookstore.Data
| Package | Version |
|---|---|
| AWSSDK.Rekognition | 3.7.400.129 |
| AWSSDK.S3 | 3.7.416.5 |
| Magick.NET-Q8-AnyCPU | 14.15.0 |
| Microsoft.EntityFrameworkCore | 10.0.14 |
| Microsoft.EntityFrameworkCore.Relational | 10.0.14 |
| Microsoft.EntityFrameworkCore.SqlServer | 10.0.14 |
| Microsoft.Extensions.Configuration.Abstractions | 10.0.14 |
| System.Security.Cryptography.Pkcs | 10.0.14 |

### Bookstore.Web
| Package | Version |
|---|---|
| Autofac | 8.2.1 |
| Autofac.Extensions.DependencyInjection | 10.0.0 |
| AWS.Logger.Core | 3.3.3 |
| AWS.Logger.NLog | 3.3.4 |
| AWSSDK.CloudWatchLogs | 3.7.410.17 |
| AWSSDK.Core | 3.7.402.35 |
| AWSSDK.Extensions.NETCore.Setup | 3.7.400.2 |
| AWSSDK.Rekognition | 3.7.400.129 |
| AWSSDK.S3 | 3.7.416.5 |
| AWSSDK.SimpleSystemsManagement | 3.7.404.10 |
| Microsoft.AspNetCore.Authentication.OpenIdConnect | 10.0.14 |
| Microsoft.EntityFrameworkCore.Design | 10.0.14 |
| Microsoft.EntityFrameworkCore.SqlServer | 10.0.14 |
| Newtonsoft.Json | 13.0.4 |
| NLog | 5.5.0 |
| NLog.Extensions.Logging | 5.5.0 |
| NLog.Web.AspNetCore | 5.5.0 |
| System.Security.Cryptography.Pkcs | 10.0.14 |

### Bookstore.Cdk
| Package | Version |
|---|---|
| Amazon.CDK.Lib | 2.275.0 |
| Cdklabs.CdkNag | 2.35.66 |
| Constructs | 10.5.0 |
| Amazon.Jsii.Analyzers | 1.142.0 |

---

## Next Steps

These items require runtime, deployment, or product decisions and cannot be resolved by a compile-only pass.

### Runtime / Database

1. **Database bootstrap** (CRITICAL): `Program.cs` does not call `Database.Migrate()` or `EnsureCreated()`. The `BookstoreDbInitializer.Seed()` data is declared via `HasData` and will only reach the database when the schema is created. For production: run `db-scripts/bobs-used-bookstore-classic-db.sql` against SQL Server. For development with an empty DB: call `new BookstoreDbInitializer().Initialize(db)` from `Program.cs` after `app.Build()`. See the ATX `nextsteps.md` §8 for full setup instructions.

2. **EF Core 10 lazy-loading** (CRITICAL): EF Core does not lazy-load navigations. All 7 repositories should be re-tested against a live DB to confirm every read path has the required `.Include()` / `.ThenInclude()` chains — missing ones produce nulls/empty collections, not exceptions.

3. **`PaginatedList.PopulateAsync()`**: Already applies `OrderBy(x => x.Id)` before `Skip`/`Take`. Verify this ordering matches the application's expected sort in all list views.

4. **EF Core model vs schema drift**: Run `dotnet ef dbcontext script --project app/Bookstore.Data --startup-project app/Bookstore.Web` and diff against `db-scripts/bobs-used-bookstore-classic-db.sql`. The `Order → Customer` explicit FK was already fixed; check for any remaining shadow-column mismatches.

### Security / Auth

5. **Cognito OIDC round-trip** (CRITICAL): Sign-in, sign-out, token refresh and claim mapping must be tested end-to-end. The `OnTokenValidated` handler upserts the customer via `ICustomerService`; confirm it executes correctly in the Core pipeline.

6. **Local auth path**: `LocalAuthenticationMiddleware` runs before `UseAuthentication`; confirm `[Authorize(Roles="Administrators")]` on `AdminAreaControllerBase` is satisfied correctly in both the local and Cognito branches.

### Packages deferred (intentional)

7. **AWSSDK v3 → v4**: All five AWSSDK packages remain on v3 (3.7.x). The v4 upgrade introduces credential/region API changes affecting `S3FileService`, `RekognitionImageValidationService`, and `ConfigurationSetup`. Defer to a dedicated upgrade cycle.

8. **Autofac 8 → 9**: Autofac 9.x is available (9.3.4). Contains breaking changes in lifetime-scope APIs. Defer; test existing registrations first.

9. **NLog 5 → 6**: NLog 6.x is available (6.2.2). Contains breaking changes in the configuration API. `LoggingSetup.cs` and `NLogExceptionHandler` would need review. Defer.

10. **Cdklabs.CdkNag 2 → 3**: Breaking major version. Defer until after CDK stacks are validated with a `cdk synth`.

11. **Constructs 10.5.0 → 10.8.1**: Minor update; safe to take alongside a CDK review cycle.

12. **Magick.NET 14.15.0 → 15.0.0**: Major API changes possible. Verify Linux native libs load in the `aspnet:10.0` container before upgrading.

### Deployment

13. **Container build / ECS**: The Dockerfile exposes port 8080 with non-root `$APP_UID`. Update load-balancer target group, security groups, and health checks. Run `cdk synth` and diff against the pre-migration template before deploying.

14. **Static assets casing**: 320 files moved to `wwwroot/` on a case-sensitive Linux filesystem. Confirm all `<link>` and `<script>` `href`/`src` values use exact matching casing (no mixed-case paths that worked on Windows/IIS).

15. **Parameter Store key shapes**: The `Services/*`, `Authentication/Cognito/*`, `Files/*` keys in `appsettings.json` must match the deployed parameter names exactly. Verify before the first AWS-backed run.
