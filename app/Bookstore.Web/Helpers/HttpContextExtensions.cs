using System;
using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    public static class HttpContextExtensions
    {
        public static string GetShoppingCartCorrelationId(this HttpContext context)
        {
            var CookieKey = "ShoppingCartId";

            // The .NET Framework version built these same options (1 year, path "/") but passed
            // them to HttpCookie without applying them, so the cookie was effectively a session
            // cookie. The options are applied here, which is the behavior the original code was
            // written to produce and what the migration plan specifies.
            var cookieOptions = new CookieOptions
            {
                Expires = DateTime.Now.AddYears(1),
                Path = "/"
            };

            var shoppingCartClientId = context.Request.Cookies[CookieKey];

            if (string.IsNullOrWhiteSpace(shoppingCartClientId))
            {
                shoppingCartClientId = context.User.Identity.IsAuthenticated ? context.User.GetSub() : Guid.NewGuid().ToString();
            }

            context.Response.Cookies.Append(CookieKey, shoppingCartClientId, cookieOptions);

            return shoppingCartClientId;
        }
    }
}
