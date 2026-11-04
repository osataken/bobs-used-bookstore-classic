using System;
using System.Security.Claims;
using System.Threading.Tasks;
using Bookstore.Domain.Customers;
using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    /// <summary>
    /// ASP.NET Core middleware that simulates authentication for local development.
    /// Replaces the OWIN-based LocalAuthenticationMiddleware from .NET Framework.
    /// </summary>
    public class LocalAuthenticationMiddleware
    {
        private const string UserId = "FB6135C7-1464-4A72-B74E-4B63D343DD09";
        private const string LocalAuthCookieName = "LocalAuthentication";

        private readonly RequestDelegate _next;

        public LocalAuthenticationMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        public async Task InvokeAsync(HttpContext context, ICustomerService customerService)
        {
            if (context.Request.Path.Value?.StartsWith("/Authentication/Login") == true)
            {
                var principal = CreateClaimsPrincipal();
                context.User = principal;

                await SaveCustomerDetailsAsync(customerService, principal);

                context.Response.Cookies.Append(LocalAuthCookieName, "1", new CookieOptions
                {
                    Expires = DateTimeOffset.Now.AddDays(1),
                    Path = "/"
                });

                context.Response.Redirect("/");
                return;
            }
            else if (context.Request.Cookies[LocalAuthCookieName] != null)
            {
                var principal = CreateClaimsPrincipal();
                context.User = principal;

                await SaveCustomerDetailsAsync(customerService, principal);
            }

            await _next(context);
        }

        private static ClaimsPrincipal CreateClaimsPrincipal()
        {
            var identity = new ClaimsIdentity("Application");

            identity.AddClaim(new Claim(ClaimTypes.Name, "bookstoreuser"));
            identity.AddClaim(new Claim("nameidentifier", UserId));
            identity.AddClaim(new Claim("given_name", "Bookstore"));
            identity.AddClaim(new Claim("family_name", "User"));
            identity.AddClaim(new Claim(ClaimTypes.Role, "Administrators"));

            return new ClaimsPrincipal(identity);
        }

        private static async Task SaveCustomerDetailsAsync(ICustomerService customerService, ClaimsPrincipal principal)
        {
            var identity = principal.Identity as ClaimsIdentity;
            if (identity == null) return;

            var dto = new CreateOrUpdateCustomerDto(
                identity.FindFirst("nameidentifier")?.Value ?? string.Empty,
                identity.Name ?? string.Empty,
                identity.FindFirst("given_name")?.Value ?? string.Empty,
                identity.FindFirst("family_name")?.Value ?? string.Empty);

            await customerService.CreateOrUpdateCustomerAsync(dto);
        }
    }
}
