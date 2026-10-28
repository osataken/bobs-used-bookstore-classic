# Migration Summary — BobsBookstoreClassic (.NET Framework → net10.0)

## Final build status

**PASS — 0 Error(s), 0 Warning(s)**

All 5 projects build cleanly (net10.0 and their netstandard2.0 legs where applicable):

| Project | TFMs | Errors | Warnings |
|---|---|---|---|
| Bookstore.Common | net10.0; netstandard2.0 | 0 | 0 |
| Bookstore.Cdk | net10.0 | 0 | 0 |
| Bookstore.Domain | net10.0; netstandard2.0 | 0 | 0 |
| Bookstore.Data | net10.0; netstandard2.0 | 0 | 0 |
| Bookstore.Web | net10.0 | 0 | 0 |

## What this cycle did (fixes on top of ATX output)

### 1. Bookstore.Common — target framework upgraded to net10.0

ATX intentionally left `Bookstore.Common` on `netstandard2.0` ("excluded by design" as a leaf
library). The task scope explicitly lists this project as a net10.0 target, so the framework
moniker was changed:

- `Bookstore.Common/Bookstore.Common.csproj`: `netstandard2.0` → `net10.0`

### 2. Amazon.CDK.Lib — vulnerability advisory resolved

The ATX run left `Amazon.CDK.Lib` at 2.188.0, which carried a known low-severity advisory
(GHSA-464c-974j-9xm6). Bumped to 2.274.0 (latest at time of run).

- `Bookstore.Cdk/Bookstore.Cdk.csproj`: `Amazon.CDK.Lib` 2.188.0 → 2.274.0
- `Bookstore.Cdk/Bookstore.Cdk.csproj`: `Constructs` 10.4.2 → 10.5.0 (required by 2.274.0;
  NU1605 downgrade error fixed)

### 3. CloudFront — deprecated constructs replaced (CS0618 / CS0612)

CDK 2.274.0 marks the old `CloudFrontWebDistribution` L2 construct and its companion types
(`SourceConfiguration`, `S3OriginConfig`, `Behavior`, `CloudFrontAllowedMethods`) as obsolete
(CS0618/CS0612). Migrated `CoreStack.CreateCloudFrontDistribution()` to the current construct set:

| Before | After |
|---|---|
| `CloudFrontWebDistribution` | `Distribution` |
| `CloudFrontWebDistributionProps` | `DistributionProps` / `BehaviorOptions` |
| `OriginAccessIdentity` + manual bucket policy | `S3BucketOrigin.WithOriginAccessControl()` (OAC) |
| `S3OriginConfig` / `S3Origin` (deprecated) | `S3BucketOrigin.WithOriginAccessControl()` |
| `CloudFrontAllowedMethods.GET_HEAD_OPTIONS` | `AllowedMethods.ALLOW_GET_HEAD_OPTIONS` |
| `distribution.DistributionDomainName` | `distribution.DomainName` |

**Infrastructure note:** OAI (Origin Access Identity) was replaced with OAC (Origin Access
Control). OAC is the AWS-recommended successor to OAI — it uses SigV4 signing and is required for
new S3 buckets in regions that do not support OAI. The bucket policy is now managed by the CDK
`S3BucketOrigin` construct rather than manually. Functional behavior (CloudFront-only access to
the S3 bucket, HTTPS redirect, GET/HEAD/OPTIONS allowed) is preserved.

Also removed the now-unused `using Amazon.CDK.AWS.IAM;` import from `CoreStack.cs`.

### 4. Bookstore.Domain — netstandard2.0 leg: missing System.ComponentModel.DataAnnotations

The `netstandard2.0` compile of `Bookstore.Domain` failed with CS0234/CS0246 because
`System.ComponentModel.DataAnnotations` (and `[Timestamp]`) is not in-box for netstandard2.0.
Fixed by adding a conditional PackageReference:

```xml
<ItemGroup Condition="'$(TargetFramework)' == 'netstandard2.0'">
  <PackageReference Include="System.ComponentModel.Annotations" Version="5.0.0" />
</ItemGroup>
```

### 5. Bookstore.Data — EF Core 10 is net10.0-only (NU1202)

EF Core 10.x packages (`Microsoft.EntityFrameworkCore`, `.Relational`, `.SqlServer` at 10.0.13)
only support `net10.0` — they emit NU1202 for the `netstandard2.0` leg.
`Microsoft.Extensions.Configuration.Abstractions/Binder` 10.0.13 also required alignment.

Resolution: split all package references into conditional `ItemGroup` blocks per TFM,
and guard every file that uses EF Core with `#if NET10_0_OR_GREATER`:

- **`Bookstore.Data.csproj`**: EF Core + Pkcs packages → `net10.0` only;
  Config Abstractions/Binder → 10.0.13 for `net10.0`, 8.0.x for `netstandard2.0`.
- **`ApplicationDbContext.cs`**: wrapped in `#if NET10_0_OR_GREATER`.
- **`BookstoreDbInitializer.cs`**: wrapped in `#if NET10_0_OR_GREATER`.
- **`PaginatedList.cs`**: EF Core-dependent members guarded; a no-op `PopulateAsync`
  stub provided for `netstandard2.0`.
- **All Repositories** (`Book`, `Order`, `Offer`, `Address`, `Customer`, `ShoppingCart`,
  `ReferenceData`): each wrapped in `#if NET10_0_OR_GREATER`.

The `netstandard2.0` target compiles cleanly (no EF Core types referenced) while the
`net10.0` target retains full data-access functionality.

### 6. BookstoreDbInitializer — EF1002 warning suppressed

`ExecuteSqlRaw` with an interpolated table name (always a hardcoded constant, never user
input) raised EF1002. Wrapped the two call sites in `#pragma warning disable/restore EF1002`
with a comment documenting why the raw call is safe.

## Next steps

None blocking. All exit criteria are met.

- The CDK stack was not deployed as part of this run. If the CloudFront distribution was previously
  deployed with an OAI-based bucket policy, a `cdk deploy` will replace the OAI with an OAC and
  update the bucket resource policy. This is a non-destructive, same-functionality change but
  should be reviewed before deploying to an existing environment.
- The `netstandard2.0` TFM in `Bookstore.Data` produces a dll where all EF Core types are excluded
  (no DbContext, no repositories). Any future .NET Framework caller would need its own data access
  implementation. If the netstandard2.0 target is no longer needed, the `<TargetFrameworks>` can
  be collapsed to `<TargetFramework>net10.0</TargetFramework>` in a subsequent cycle.
