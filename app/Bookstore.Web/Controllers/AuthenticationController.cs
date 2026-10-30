using System.Threading.Tasks;
using BobsBookstoreClassic.Data;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.Controllers
{
    public class AuthenticationController : Controller
    {
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
            if (Request.Cookies["LocalAuthentication"] != null)
            {
                // Was an expired HttpCookie added to the response.
                Response.Cookies.Delete("LocalAuthentication");
            }

            return RedirectToAction("Index", "Home");
        }

        private async Task<IActionResult> CognitoSignOutAsync()
        {
            // Was expiring the ".AspNet.Cookies" cookie by hand; the cookie handler owns the
            // cookie name on Core, so sign out of the scheme instead.
            await HttpContext.SignOutAsync(CookieAuthenticationDefaults.AuthenticationScheme);

            var domain = BookstoreConfiguration.GetSetting("Authentication/Cognito/CognitoDomain");
            var clientId = BookstoreConfiguration.GetSetting("Authentication/Cognito/LocalClientId");
            var port = Request.Host.Port ?? (Request.IsHttps ? 443 : 80);
            var logoutUri = $"{Request.Scheme}://{Request.Host.Host}:{port}/";

            return Redirect($"{domain}/logout?client_id={clientId}&logout_uri={logoutUri}");
        }
    }
}
