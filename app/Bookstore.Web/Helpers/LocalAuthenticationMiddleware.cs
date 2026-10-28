using System;
using System.Security.Claims;
using System.Threading.Tasks;
using Bookstore.Domain.Customers;
using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    /// <summary>
    /// ASP.NET Core port of the OWIN LocalAuthenticationMiddleware. Supplies a fixed development
    /// identity when Services/Authentication is not "aws".
    /// </summary>
    public class LocalAuthenticationMiddleware
    {
        public const string CookieName = "LocalAuthentication";

        private const string UserId = "FB6135C7-1464-4A72-B74E-4B63D343DD09";

        private readonly RequestDelegate _next;

        public LocalAuthenticationMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        public async Task InvokeAsync(HttpContext context, ICustomerService customerService)
        {
            if (context.Request.Path.HasValue && context.Request.Path.Value.StartsWith("/Authentication/Login"))
            {
                var principal = CreateClaimsPrincipal(context);

                await SaveCustomerDetailsAsync(principal, customerService);

                context.Response.Cookies.Append(CookieName, UserId, new CookieOptions
                {
                    Expires = DateTimeOffset.Now.AddDays(1),
                    HttpOnly = true,
                    IsEssential = true
                });

                context.Response.Redirect("/");
            }
            else if (context.Request.Cookies.ContainsKey(CookieName))
            {
                var principal = CreateClaimsPrincipal(context);

                await SaveCustomerDetailsAsync(principal, customerService);

                await _next(context);
            }
            else
            {
                await _next(context);
            }
        }

        private static ClaimsPrincipal CreateClaimsPrincipal(HttpContext context)
        {
            var identity = new ClaimsIdentity("Application", ClaimTypes.Name, ClaimTypes.Role);

            identity.AddClaim(new Claim(ClaimTypes.Name, "bookstoreuser"));
            identity.AddClaim(new Claim("nameidentifier", UserId));
            identity.AddClaim(new Claim("given_name", "Bookstore"));
            identity.AddClaim(new Claim("family_name", "User"));
            identity.AddClaim(new Claim(ClaimTypes.Role, "Administrators"));

            var principal = new ClaimsPrincipal(identity);

            context.User = principal;

            return principal;
        }

        private static async Task SaveCustomerDetailsAsync(ClaimsPrincipal principal, ICustomerService customerService)
        {
            var identity = (ClaimsIdentity)principal.Identity;

            var dto = new CreateOrUpdateCustomerDto(
                identity.FindFirst("nameidentifier").Value,
                identity.Name,
                identity.FindFirst("given_name").Value,
                identity.FindFirst("family_name").Value);

            await customerService.CreateOrUpdateCustomerAsync(dto);
        }
    }
}
