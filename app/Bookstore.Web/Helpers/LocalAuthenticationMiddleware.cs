using System;
using System.Security.Claims;
using System.Threading.Tasks;
using Bookstore.Domain.Customers;
using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    /// <summary>
    /// Core port of the OWIN LocalAuthenticationMiddleware. Behaviour is unchanged:
    /// a request to /Authentication/Login signs the well-known local user in (cookie marker
    /// + customer upsert) and redirects to "/", a request carrying the LocalAuthentication
    /// cookie is authenticated for the duration of the request, anything else passes through
    /// untouched.
    /// </summary>
    public class LocalAuthenticationMiddleware
    {
        private const string UserId = "FB6135C7-1464-4A72-B74E-4B63D343DD09";
        private const string CookieName = "LocalAuthentication";

        private readonly RequestDelegate _next;

        public LocalAuthenticationMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        // ICustomerService is scoped, so it is injected per-request as an InvokeAsync parameter.
        public async Task InvokeAsync(HttpContext context, ICustomerService customerService)
        {
            if (context.Request.Path.StartsWithSegments("/Authentication/Login", StringComparison.OrdinalIgnoreCase))
            {
                CreateClaimsPrincipal(context);

                await SaveCustomerDetailsAsync(context, customerService);

                context.Response.Cookies.Append(CookieName, UserId, new CookieOptions
                {
                    Expires = DateTimeOffset.Now.AddDays(1),
                    Path = "/"
                });

                context.Response.Redirect("/");

                // Short-circuit: the original never invoked Next on this branch.
                return;
            }

            if (!string.IsNullOrEmpty(context.Request.Cookies[CookieName]))
            {
                CreateClaimsPrincipal(context);

                await SaveCustomerDetailsAsync(context, customerService);

                await _next(context);

                return;
            }

            await _next(context);
        }

        private static void CreateClaimsPrincipal(HttpContext context)
        {
            var identity = new ClaimsIdentity("Application", ClaimTypes.Name, ClaimTypes.Role);

            identity.AddClaim(new Claim(ClaimTypes.Name, "bookstoreuser"));
            identity.AddClaim(new Claim("nameidentifier", UserId));
            identity.AddClaim(new Claim("given_name", "Bookstore"));
            identity.AddClaim(new Claim("family_name", "User"));
            identity.AddClaim(new Claim(ClaimTypes.Role, "Administrators"));

            context.User = new ClaimsPrincipal(identity);
        }

        private static async Task SaveCustomerDetailsAsync(HttpContext context, ICustomerService customerService)
        {
            var identity = (ClaimsIdentity)context.User.Identity;

            var dto = new CreateOrUpdateCustomerDto(
                identity.FindFirst("nameidentifier").Value,
                identity.Name,
                identity.FindFirst("given_name").Value,
                identity.FindFirst("family_name").Value);

            await customerService.CreateOrUpdateCustomerAsync(dto);
        }
    }
}
