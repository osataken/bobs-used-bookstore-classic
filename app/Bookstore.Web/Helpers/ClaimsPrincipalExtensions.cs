using System.Linq;
using System.Security.Claims;
using System.Security.Principal;

namespace Bookstore.Web.Helpers
{
    public static class ClaimsPrincipalExtensions
    {
        public static string GetSub(this IPrincipal claimsPrincipal)
        {
            return (claimsPrincipal as ClaimsPrincipal)?.Identities
                .Select(x => x.GetSub())
                .FirstOrDefault(x => x != null);
        }

        public static string GetSub(this ClaimsIdentity identity)
        {
            // ASP.NET Core does not apply the legacy inbound claim-type mapping, so the Cognito
            // path leaves the subject in a raw "sub" claim while the local-auth path emits
            // "nameidentifier". Both (and the full SOAP ClaimTypes.NameIdentifier URI) resolve.
            return identity?.Claims
                .FirstOrDefault(x => x.Type.Contains("nameidentifier") || x.Type == "sub")
                ?.Value;
        }
    }
}
