# Migration Summary — .NET Framework → .NET 10

## Status
✅ All 5 projects build with **0 errors and 0 warnings**.

| Project | TFM | Result |
|---|---|---|
| Bookstore.Common | netstandard2.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Domain | net10.0; netstandard2.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Data | net10.0; netstandard2.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Cdk | net10.0 | ✅ 0 errors, 0 warnings |
| Bookstore.Web | net10.0 | ✅ 0 errors, 0 warnings |

---

## Changes Made

### Bookstore.Domain — `Bookstore.Domain.csproj`
- Added `System.ComponentModel.Annotations 5.0.0` package reference.  
  The .NET 10 SDK reference assemblies for netstandard2.0 require an explicit package reference for `System.ComponentModel.DataAnnotations` (used by `[Timestamp]` in `Entity.cs`).

### Bookstore.Data — `Bookstore.Data.csproj`
- Made EF Core 10.x package references conditional on `'$(TargetFramework)'=='net10.0'`:
  - `Microsoft.EntityFrameworkCore 10.0.14`
  - `Microsoft.EntityFrameworkCore.Relational 10.0.14`
  - `Microsoft.EntityFrameworkCore.SqlServer 10.0.14`

  EF Core 10.x only ships a `net10.0` TFM; the library must multi-target to `netstandard2.0` so .NET Framework callers can still reference it. The remaining packages (`AWSSDK.*`, `Magick.NET`, `Microsoft.Extensions.Configuration.Abstractions`, `System.Security.Cryptography.Pkcs`) all publish `netstandard2.0` assemblies and remain unconditional.

### Bookstore.Data — EF Core code files (`#if NET10_0_OR_GREATER` guards)
All files whose top-level EF Core `using` directives would fail to resolve on `netstandard2.0` were wrapped with `#if NET10_0_OR_GREATER` / `#endif`:

- `ApplicationDbContext.cs`
- `BookstoreDbInitializer.cs`
- `PaginatedList.cs`
- `Repositories/AddressRepository.cs`
- `Repositories/BookRepository.cs`
- `Repositories/CustomerRepository.cs`
- `Repositories/OfferRepository.cs`
- `Repositories/OrderRepository.cs`
- `Repositories/ReferenceDataRepository.cs`
- `Repositories/ShoppingCartRepository.cs`

The netstandard2.0 build of Bookstore.Data therefore exports the AWS, Magick.NET, and configuration service classes (which work across TFMs) while the EF Core data-access layer is only present in the net10.0 build.

### Bookstore.Cdk — `CoreStack.cs`
Removed use of the obsolete `CloudFrontWebDistribution` API (warnings CS0618/CS0612) by migrating to the modern CDK `Distribution` construct with `S3BucketOrigin.WithOriginAccessControl`:

- **Before**: `CloudFrontWebDistribution` + `OriginAccessIdentity` + manual S3 bucket policy  
- **After**: `Distribution` + `S3BucketOrigin.WithOriginAccessControl()` (CDK manages the OAC and bucket policy automatically)
- `distribution.DistributionDomainName` → `distribution.DomainName`
- Removed unused `Amazon.CDK.AWS.IAM` using directive and all OAI/policy-statement code

  This is functionally equivalent — S3 access is still restricted to traffic from the CloudFront distribution — but uses OAC (Origin Access Control) which is the AWS-recommended successor to OAI.

---

## Next Steps

None required for the build to succeed. The following items are recorded for awareness:

1. **EF Core netstandard2.0 stub**: Bookstore.Data's `netstandard2.0` output omits all EF Core data-access types. If a legacy .NET Framework caller needs concrete repository implementations (rather than just the AWS/image service classes) it will need to upgrade to net10.0 first.

2. **CDK OAI → OAC migration**: Existing deployed CloudFormation stacks that used the old `CloudFrontWebDistribution` / OAI resources will need a CDK diff and deploy to pick up the new `Distribution` / OAC resources. This is a non-breaking change for new deployments but requires a stack update for existing deployments.
