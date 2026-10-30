using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    /// <summary>Was OwinRequestExtensions (extension on IOwinRequest).</summary>
    public static class HttpRequestExtensions
    {
        public static string GetReturnUrl(this HttpRequest request)
        {
            return $"{request.Scheme}://{request.Host}/signin-oidc";
        }
    }
}
