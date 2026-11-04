using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    /// <summary>
    /// ASP.NET Core replacement for the OWIN <c>IOwinRequest.GetReturnUrl()</c> extension.
    /// Builds the absolute OpenID Connect callback URI from the current request, honouring
    /// the forwarded scheme/host set by <c>UseForwardedHeaders</c> when the app runs behind
    /// a TLS-terminating load balancer.
    /// </summary>
    public static class HttpRequestExtensions
    {
        public const string OidcCallbackPath = "/signin-oidc";

        public static string GetReturnUrl(this HttpRequest request)
        {
            return $"{request.Scheme}://{request.Host}{OidcCallbackPath}";
        }
    }
}
