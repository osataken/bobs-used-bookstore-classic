# Migration Summary — BobsBookstoreClassic (.NET Framework → net10.0)

## Final build status

**PASS — 0 Error(s), 0 Warning(s)**

All 5 projects build cleanly against net10.0:

| Project | TFM | Errors | Warnings |
|---|---|---|---|
| Bookstore.Common | net10.0 | 0 | 0 |
| Bookstore.Cdk | net10.0 | 0 | 0 |
| Bookstore.Domain | net10.0 | 0 | 0 |
| Bookstore.Data | net10.0 | 0 | 0 |
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

## Next steps

None blocking. All exit criteria are met.

- The CDK stack was not deployed as part of this run. If the CloudFront distribution was previously
  deployed with an OAI-based bucket policy, a `cdk deploy` will replace the OAI with an OAC and
  update the bucket resource policy. This is a non-destructive, same-functionality change but
  should be reviewed before deploying to an existing environment.
