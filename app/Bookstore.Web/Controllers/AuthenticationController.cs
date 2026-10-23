using BobsBookstoreClassic.Data;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authentication.OpenIdConnect;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using System.Threading.Tasks;

namespace Bookstore.Web.Controllers
{
    // Reachable while signed out: Login must be able to start the challenge and LogOut must be
    // able to complete even after the auth cookie is gone (the global AuthorizeFilter would
    // otherwise block both).
    [AllowAnonymous]
    public class AuthenticationController : Controller
    {
        public async Task<IActionResult> Login(string redirectUri = null)
        {
            // Under OWIN the global AuthorizeAttribute triggered the OIDC challenge before this
            // action ran. With [AllowAnonymous] the challenge has to be issued explicitly.
            if (BookstoreConfiguration.GetSetting("Services/Authentication") == "aws" &&
                !(User.Identity?.IsAuthenticated ?? false))
            {
                return Challenge(
                    new AuthenticationProperties { RedirectUri = string.IsNullOrWhiteSpace(redirectUri) ? "/" : redirectUri },
                    OpenIdConnectDefaults.AuthenticationScheme);
            }

            await Task.CompletedTask;

            if (string.IsNullOrWhiteSpace(redirectUri)) return RedirectToAction("Index", "Home");

            return Redirect(redirectUri);
        }

        public async Task<IActionResult> LogOut()
        {
            return BookstoreConfiguration.GetSetting("Services/Authentication") == "aws"
                ? await CognitoSignOutAsync()
                : LocalSignOut();
        }

        private IActionResult LocalSignOut()
        {
            if (Request.Cookies.ContainsKey("LocalAuthentication"))
            {
                Response.Cookies.Delete("LocalAuthentication");
            }

            return RedirectToAction("Index", "Home");
        }

        private async Task<IActionResult> CognitoSignOutAsync()
        {
            // Replaces the manual expiry of the ".AspNet.Cookies" ticket cookie.
            await HttpContext.SignOutAsync(CookieAuthenticationDefaults.AuthenticationScheme);

            var domain = BookstoreConfiguration.GetSetting("Authentication/Cognito/CognitoDomain");
            var clientId = BookstoreConfiguration.GetSetting("Authentication/Cognito/LocalClientId");
            var logoutUri = $"{Request.Scheme}://{Request.Host}/";

            // Cognito's /logout endpoint ends the federated session, so the OIDC handler's own
            // sign-out redirect is not used here (matching the original behaviour).
            return Redirect($"{domain}/logout?client_id={clientId}&logout_uri={logoutUri}");
        }
    }
}
