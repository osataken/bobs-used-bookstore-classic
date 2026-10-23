using Microsoft.AspNetCore.Http;
using System;

namespace Bookstore.Web.Helpers
{
    public static class HttpContextExtensions
    {
        private const string CookieKey = "ShoppingCartId";

        /// <summary>
        /// Returns the correlation id used to associate an anonymous (or authenticated) visitor
        /// with their shopping cart / wishlist: read the cookie, fall back to the authenticated
        /// subject or a new Guid, then always write the value back with a one-year expiry rooted
        /// at "/". NOTE: the System.Web version built an identical one-year CookieOptions but
        /// never passed it to HttpCookie, so it actually emitted a session cookie. The intended
        /// one-year expiry is applied here, so anonymous carts now survive a browser restart.
        /// </summary>
        public static string GetShoppingCartCorrelationId(this HttpContext context)
        {
            var shoppingCartClientId = context.Request.Cookies[CookieKey];

            if (string.IsNullOrWhiteSpace(shoppingCartClientId))
            {
                shoppingCartClientId = (context.User?.Identity?.IsAuthenticated ?? false)
                    ? context.User.GetSub()
                    : Guid.NewGuid().ToString();
            }

            // GetSub() can still return null if the principal carries no subject claim; never
            // persist a null correlation id, otherwise every cart lookup would collapse together.
            if (string.IsNullOrWhiteSpace(shoppingCartClientId))
            {
                shoppingCartClientId = Guid.NewGuid().ToString();
            }

            context.Response.Cookies.Append(CookieKey, shoppingCartClientId, new CookieOptions
            {
                Expires = DateTimeOffset.Now.AddYears(1),
                Path = "/"
            });

            return shoppingCartClientId;
        }
    }
}
