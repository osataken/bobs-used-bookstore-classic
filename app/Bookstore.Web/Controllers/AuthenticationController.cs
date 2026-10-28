using System;
using System.Threading.Tasks;
using BobsBookstoreClassic.Data;
using Bookstore.Web.Helpers;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.Controllers
{
    // Anonymous access is required: Program.cs installs a global FallbackPolicy requiring an
    // authenticated user, and these endpoints run before/after sign-in.
    [AllowAnonymous]
    public class AuthenticationController : Controller
    {
        public IActionResult Login(string redirectUri = null)
        {
            if (string.IsNullOrWhiteSpace(redirectUri)) return RedirectToAction("Index", "Home");

            return Redirect(redirectUri);
        }

        public Task<IActionResult> LogOut()
        {
            return AuthenticationSetup.IsAwsAuthentication() ? CognitoSignOutAsync() : Task.FromResult(LocalSignOut());
        }

        private IActionResult LocalSignOut()
        {
            if (Request.Cookies.ContainsKey(LocalAuthenticationMiddleware.CookieName))
            {
                // Core equivalent of writing an already-expired HttpCookie.
                Response.Cookies.Delete(LocalAuthenticationMiddleware.CookieName);
            }

            return RedirectToAction("Index", "Home");
        }

        private async Task<IActionResult> CognitoSignOutAsync()
        {
            // Replaces the manual expiry of the ".AspNet.Cookies" OWIN cookie.
            await HttpContext.SignOutAsync(CookieAuthenticationDefaults.AuthenticationScheme);

            var domain = BookstoreConfiguration.GetSetting("Authentication/Cognito/CognitoDomain");
            var clientId = BookstoreConfiguration.GetSetting("Authentication/Cognito/LocalClientId");
            var logoutUri = $"{Request.Scheme}://{Request.Host.Host}:{Request.Host.Port ?? (Request.IsHttps ? 443 : 80)}/";

            return Redirect($"{domain}/logout?client_id={clientId}&logout_uri={logoutUri}");
        }
    }
}
