using System.Threading.Tasks;
using BobsBookstoreClassic.Data;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.Controllers
{
    public class AuthenticationController : Controller
    {
        // Reaching this action requires an authenticated user (global authorization fallback
        // policy), which triggers the OpenIdConnect challenge in the Cognito configuration and
        // the local authentication middleware otherwise — same behavior as the MVC5 version.
        public IActionResult Login(string redirectUri = null)
        {
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
            await HttpContext.SignOutAsync(CookieAuthenticationDefaults.AuthenticationScheme);

            if (Request.Cookies.ContainsKey(".AspNetCore.Cookies"))
            {
                Response.Cookies.Delete(".AspNetCore.Cookies");
            }

            var domain = BookstoreConfiguration.GetSetting("Authentication/Cognito/CognitoDomain");
            var clientId = BookstoreConfiguration.GetSetting("Authentication/Cognito/LocalClientId");

            // Request.Host carries the port when one is present, matching the MVC5
            // "{scheme}://{host}:{port}/" logout URI.
            var logoutUri = $"{Request.Scheme}://{Request.Host}/";

            return Redirect($"{domain}/logout?client_id={clientId}&logout_uri={logoutUri}");
        }
    }
}
