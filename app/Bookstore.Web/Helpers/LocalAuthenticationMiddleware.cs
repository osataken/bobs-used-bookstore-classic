using System;
using System.Security.Claims;
using System.Threading.Tasks;
using Bookstore.Domain.Customers;
using Microsoft.AspNetCore.Owin;

using Microsoft.AspNetCore.Http;


namespace Bookstore.Web.Helpers
{
public class LocalAuthenticationMiddleware     {
RequestDelegate _next = null;        private const string UserId = "FB6135C7-1464-4A72-B74E-4B63D343DD09";

        private readonly ICustomerService _customerService;
        private readonly IHttpContextAccessor _httpContextAccessor;

        public LocalAuthenticationMiddleware(RequestDelegate next, ICustomerService customerService, IHttpContextAccessor httpContextAccessor)         {
            _customerService = customerService;
            _httpContextAccessor = httpContextAccessor;
_next = next;        }
public async Task Invoke(HttpContext context)
        {
            if (context.Request.Path.Value.StartsWith("/Authentication/Login"))
            {
                CreateClaimsPrincipal(context);

                await SaveCustomerDetailsAsync();

                var cookieOptions = new CookieOptions { Expires = DateTime.Now.AddDays(1), Path = "/" };
                context.Response.Cookies.Append("LocalAuthentication", "true", cookieOptions);

                context.Response.Redirect("/");
                return;
            }
else if (/* Added by CTA: TODO: Replace HttpContext.Current with dependency injection pattern using IHttpContextAccessor. */
_httpContextAccessor.HttpContext.Request.Cookies["LocalAuthentication"] != null)
            {
                CreateClaimsPrincipal(context);

                await SaveCustomerDetailsAsync();

                await _next.Invoke(context);
                return;
            }
            else
            {
                await _next.Invoke(context);
                return;
            }
        }

        private void CreateClaimsPrincipal(HttpContext context)
        {
            var identity = new ClaimsIdentity("Application");

            identity.AddClaim(new Claim(ClaimTypes.Name, "bookstoreuser"));
            identity.AddClaim(new Claim("nameidentifier", UserId));
            identity.AddClaim(new Claim("given_name", "Bookstore"));
            identity.AddClaim(new Claim("family_name", "User"));
            identity.AddClaim(new Claim(ClaimTypes.Role, "Administrators"));

            context.User = new ClaimsPrincipal(identity);
        }

        private async Task SaveCustomerDetailsAsync()
        {
var identity = (ClaimsIdentity)_httpContextAccessor.HttpContext.User.Identity;

            var dto = new CreateOrUpdateCustomerDto(
                identity.FindFirst("nameidentifier").Value,
                identity.Name,
                identity.FindFirst("given_name").Value,
                identity.FindFirst("family_name").Value);

            await _customerService.CreateOrUpdateCustomerAsync(dto);
        }
    }
}
