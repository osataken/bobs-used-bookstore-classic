using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    /// <summary>
    /// Port of the OWIN <c>OwinRequestExtensions</c> (IOwinRequest → HttpRequest).
    /// </summary>
    public static class HttpRequestExtensions
    {
        public static string GetReturnUrl(this HttpRequest request)
        {
            return $"{request.Scheme}://{request.Host}/signin-oidc";
        }
    }
}
