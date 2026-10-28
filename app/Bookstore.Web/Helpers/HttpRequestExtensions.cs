using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    public static class HttpRequestExtensions
    {
        /// <summary>
        /// The OIDC redirect URI used by both the authorization request and the code redemption
        /// request. Ported from the OWIN IOwinRequest.GetReturnUrl() extension.
        /// </summary>
        public static string GetReturnUrl(this HttpRequest request)
        {
            return $"{request.Scheme}://{request.Host}/signin-oidc";
        }
    }
}
