using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    /// <summary>
    /// Replaces Helpers/IOwinRequestExtensions.cs (OwinRequestExtensions.GetReturnUrl).
    /// The OIDC redirect URI behaviour is preserved exactly: scheme + host + /signin-oidc.
    /// </summary>
    public static class HttpRequestExtensions
    {
        public static string GetReturnUrl(this HttpRequest request)
        {
            return $"{request.Scheme}://{request.Host}/signin-oidc";
        }
    }
}
