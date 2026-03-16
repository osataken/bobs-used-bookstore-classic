using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    public static class IOwinRequestExtensions
    {
        public static string GetReturnUrl(this HttpRequest request)
        {
            return $"{request.Scheme}://{request.Host}{request.PathBase}{request.Path}";
        }
    }
}
