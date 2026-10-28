# Migration Summary — BobsBookstoreClassic (.NET Framework → net10.0)

## Final build status

**PASS — 0 Error(s), 0 Warning(s)**

All 5 projects build cleanly:

| Project | TFMs | Errors | Warnings |
|---|---|---|---|
| Bookstore.Common | net10.0 + netstandard2.0 | 0 | 0 |
| Bookstore.Cdk | net10.0 | 0 | 0 |
| Bookstore.Domain | net10.0 + netstandard2.0 | 0 | 0 |
| Bookstore.Data | net10.0 + netstandard2.0 | 0 | 0 |
| Bookstore.Web | net10.0 | 0 | 0 |

## What this cycle did

### 1. Bookstore.Domain — System.ComponentModel.DataAnnotations on netstandard2.0

The `[Timestamp]` attribute on `Entity.cs` requires `System.ComponentModel.DataAnnotations`.
In .NET 10 SDK toolchains targeting `netstandard2.0`, this assembly is no longer auto-referenced.

- `Bookstore.Domain/Bookstore.Domain.csproj`: added a conditional
  `<PackageReference Include="System.ComponentModel.Annotations" Version="5.0.0" Condition="'$(TargetFramework)'=='netstandard2.0'" />`

### 2. Bookstore.Data — EF Core 10 not compatible with netstandard2.0 (NU1202)

`Microsoft.EntityFrameworkCore` 10.x and `Magick.NET-Q8-AnyCPU` 15.x only target `net10.0`,
not `netstandard2.0`. The instruction rules require `#if` guards rather than TFM removal for
NU1202 (package incompatibility).

**Project file changes** (`Bookstore.Data.csproj`):
- Moved `Microsoft.EntityFrameworkCore`, `Microsoft.EntityFrameworkCore.Relational`,
  `Microsoft.EntityFrameworkCore.SqlServer`, `Magick.NET-Q8-AnyCPU`, and
  `System.Security.Cryptography.Pkcs` into a `Condition="'$(TargetFramework)'!='netstandard2.0'"`
  `<ItemGroup>` so they are only restored for `net10.0`.
- Kept `AWSSDK.*` and `Microsoft.Extensions.Configuration.*` packages unconditional (they support
  both TFMs).

**Code guards** — `#if !NETSTANDARD2_0` / `#endif` wrappers added to all EF Core-dependent files
so the `netstandard2.0` leg compiles cleanly (no symbols, no types):
- `ApplicationDbContext.cs` — whole file
- `BookstoreDbInitializer.cs` — whole file
- `PaginatedList.cs` — `using Microsoft.EntityFrameworkCore` directive and `PopulateAsync()` body;
  the method throws `PlatformNotSupportedException` on `netstandard2.0`
- `Repositories/AddressRepository.cs` — whole file
- `Repositories/BookRepository.cs` — whole file
- `Repositories/CustomerRepository.cs` — whole file
- `Repositories/OfferRepository.cs` — whole file
- `Repositories/OrderRepository.cs` — whole file
- `Repositories/ReferenceDataRepository.cs` — whole file
- `Repositories/ShoppingCartRepository.cs` — whole file
- `ImageResizeService/ImageResizeService.cs` — whole file (depends on Magick.NET)

### 3. Bookstore.Data — EF1002 SQL injection advisory resolved

`BookstoreDbInitializer.ExecuteSqlRaw` was called with an interpolated string containing the table
name. Table names cannot be parameterised in `SET IDENTITY_INSERT` SQL, so `ExecuteSql` cannot be
used. The table name is a hardcoded application constant (never user input), so the call is safe.
Added `#pragma warning disable/restore EF1002` around the two calls with an inline comment
explaining the rationale.

## Previous cycle (ATX + fixes)

### 1. Bookstore.Common — target framework
`netstandard2.0` → `net10.0` (scope explicitly targets net10.0).

### 2. Amazon.CDK.Lib — vulnerability advisory
`Amazon.CDK.Lib` 2.188.0 → 2.274.0; `Constructs` 10.4.2 → 10.5.0.

### 3. CloudFront — deprecated constructs replaced
`CloudFrontWebDistribution` migrated to `Distribution` / `S3BucketOrigin.WithOriginAccessControl()`
(OAI → OAC). See previous summary for full mapping table.

## Next steps

None blocking. All exit criteria are met:
- `dotnet build BobsBookstoreClassic.sln` exits with code 0, 0 errors, 0 warnings.
- All five target projects compile for their declared TFMs.

Informational notes:
- The `netstandard2.0` build of `Bookstore.Data` exposes only the non-EF-Core services
  (file services, image validation, configuration). Any .NET Framework caller that references
  this library cannot use the repository/DbContext layer — that is expected because EF Core 10
  does not run on .NET Framework.
- The CDK stack was not deployed as part of this run. The previously noted OAI→OAC CloudFront
  change remains pending deployment review.
