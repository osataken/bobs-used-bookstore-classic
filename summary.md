# .NET Framework → net10.0 Migration Summary

## Final Build Status

All 5 in-scope projects build successfully:

| Project | Target Framework | Errors | Warnings |
|---|---|---|---|
| Bookstore.Cdk | net10.0 | 0 | 0 |
| Bookstore.Common | net10.0 | 0 | 0 |
| Bookstore.Web | net10.0 | 0 | 0 |
| Bookstore.Domain | net10.0 | 0 | 0 |
| Bookstore.Data | net10.0 | 0 | 0 |

**Total: 0 errors, 0 warnings** (`dotnet build BobsBookstoreClassic.sln`)

---

## Changes Made

### AWS Transform Custom (atx) — Starting Point
The atx tool had already applied the bulk of the .NET Framework → net10.0 migration:
- Converted all `.csproj` files to SDK-style targeting `net10.0`
- Migrated `Global.asax` → `Program.cs`
- Replaced `System.Web` MVC with ASP.NET Core
- Migrated EF6 → EF Core 10
- Replaced `Web.config` with `appsettings.json`
- Added ASP.NET Core authentication (OpenIdConnect), NLog, Autofac DI

### Changes Applied in This Cycle

#### 1. Bookstore.Common — Target Framework Update
**File:** `app/Bookstore.Common/Bookstore.Common.csproj`  
**Change:** Updated `TargetFramework` from `netstandard2.0` → `net10.0` as required by the migration scope.

#### 2. Bookstore.Data — CVE Vulnerability Fix
**File:** `app/Bookstore.Data/Bookstore.Data.csproj`  
**Change:** Added explicit `PackageReference` for `System.Security.Cryptography.Pkcs 10.0.14` to override the vulnerable transitive version `9.0.4` (pulled in via `Microsoft.EntityFrameworkCore.SqlServer → Microsoft.Data.SqlClient 6.1.1`).  
- Fixed: GHSA-2fmv-rv2c-59x7 (moderate) and GHSA-846f-f764-76mx (high)

#### 3. Bookstore.Web — CVE Vulnerability Fix
**File:** `app/Bookstore.Web/Bookstore.Web.csproj`  
**Change:** Same `System.Security.Cryptography.Pkcs 10.0.14` pin applied to the web project (inherits the same transitive chain via `Microsoft.EntityFrameworkCore.SqlServer`).

#### 4. Bookstore.Cdk — CDK Package Updates + CloudFront API Migration
**File:** `app/Bookstore.Cdk/Bookstore.Cdk.csproj`  
**Change:** Updated packages to fix NU1901 low severity vulnerability (GHSA-464c-974j-9xm6):
- `Amazon.CDK.Lib`: 2.188.0 → **2.275.0**
- `Cdklabs.CdkNag`: 2.35.66 → **2.38.2**
- `Constructs`: 10.4.2 → **10.8.1**

**File:** `app/Bookstore.Cdk/CoreStack.cs`  
**Change:** Migrated deprecated CDK CloudFront API to modern equivalents:
- Removed: `CloudFrontWebDistributionProps`, `SourceConfiguration`, `S3OriginConfig`, `Behavior`, `CloudFrontAllowedMethods` (all obsolete)
- Added: `using Amazon.CDK.AWS.CloudFront.Origins;`
- Replaced `new CloudFrontWebDistribution(...)` → `new Distribution(...)` with `DistributionProps` / `BehaviorOptions`
- Replaced obsolete `new S3Origin(...)` → `S3BucketOrigin.WithOriginAccessIdentity(...)` with `S3BucketOriginWithOAIProps`
- Replaced `CloudFrontAllowedMethods.GET_HEAD_OPTIONS` → `AllowedMethods.ALLOW_GET_HEAD_OPTIONS`

---

## Next Steps

None required — all in-scope projects build with zero errors and zero warnings.

### Notes for Future Consideration
- `Microsoft.Data.SqlClient` is version 6.1.1 (transitive via EF Core 10.0.0). If a newer version is released that pins a non-vulnerable `System.Security.Cryptography.Pkcs` version directly, the explicit pin in Bookstore.Data and Bookstore.Web can be removed.
- CDK CloudFront was migrated from OAI (Origin Access Identity) to continue using OAI. Future migration to OAC (Origin Access Control) is possible and recommended by AWS for new deployments, but is out of scope for this framework upgrade.
