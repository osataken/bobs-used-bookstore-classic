using System.Linq;
using System.Security.Claims;
using System.Threading.Tasks;
using BobsBookstoreClassic.Data;
using Bookstore.Domain.Customers;
using Bookstore.Web.Helpers;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authentication.OpenIdConnect;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.IdentityModel.Protocols.OpenIdConnect;
using Microsoft.IdentityModel.Tokens;

namespace Bookstore.Web
{
    /// <summary>
    /// Port of App_Start/AuthenticationSetup.cs (OWIN AuthenticationConfig).
    /// The original chose between a Cognito OpenID Connect pipeline and the local
    /// authentication middleware based on the Services/Authentication setting; the same
    /// branch is preserved here (service registration side) and in Program.cs
    /// (middleware side, task 34).
    /// </summary>
    public static class AuthenticationSetup
    {
        public static bool UseCognitoAuthentication =>
            BookstoreConfiguration.GetSetting("Services/Authentication") == "aws";

        public static void ConfigureAuthentication(IServiceCollection services)
        {
            if (UseCognitoAuthentication)
            {
                ConfigureCognitoAuthentication(services);
            }
            else
            {
                ConfigureLocalAuthentication(services);
            }
        }

        // The local path authenticates through LocalAuthenticationMiddleware (registered in
        // Program.cs). A cookie scheme is still registered so that the authentication and
        // authorization middleware always have a backing default scheme.
        private static void ConfigureLocalAuthentication(IServiceCollection services)
        {
            services.AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
                .AddCookie();
        }

        private static void ConfigureCognitoAuthentication(IServiceCollection services)
        {
            // SetDefaultSignInAsAuthenticationType(Cookies) + UseCookieAuthentication +
            // UseOpenIdConnectAuthentication. DefaultChallengeScheme makes the OIDC handler
            // service challenges (the OWIN pipeline did this implicitly by ordering).
            services.AddAuthentication(options =>
            {
                options.DefaultScheme = CookieAuthenticationDefaults.AuthenticationScheme;
                options.DefaultSignInScheme = CookieAuthenticationDefaults.AuthenticationScheme;
                options.DefaultChallengeScheme = OpenIdConnectDefaults.AuthenticationScheme;
            })
            .AddCookie()
            .AddOpenIdConnect(options =>
            {
                options.ClientId = BookstoreConfiguration.GetSetting("Authentication/Cognito/LocalClientId");
                options.MetadataAddress = BookstoreConfiguration.GetSetting("Authentication/Cognito/MetadataAddress");
                options.ResponseType = OpenIdConnectResponseType.Code;
                options.CallbackPath = "/signin-oidc";
                options.Scope.Clear();
                options.Scope.Add("openid");
                options.Scope.Add("profile");
                options.SignInScheme = CookieAuthenticationDefaults.AuthenticationScheme;
                options.UseTokenLifetime = false;
                options.SaveTokens = true;
                options.TokenValidationParameters = new TokenValidationParameters
                {
                    NameClaimType = "cognito:username",
                    RoleClaimType = "cognito:groups"
                };

                // Notifications.RedirectToIdentityProvider
                options.Events.OnRedirectToIdentityProvider = context =>
                {
                    context.ProtocolMessage.RedirectUri = context.Request.GetReturnUrl();

                    return Task.CompletedTask;
                };

                // Notifications.AuthorizationCodeReceived
                options.Events.OnAuthorizationCodeReceived = context =>
                {
                    if (context.TokenEndpointRequest != null)
                    {
                        context.TokenEndpointRequest.RedirectUri = context.Request.GetReturnUrl();
                    }

                    return Task.CompletedTask;
                };

                // Notifications.SecurityTokenValidated
                options.Events.OnTokenValidated = async context =>
                {
                    var customerService = context.HttpContext.RequestServices.GetRequiredService<ICustomerService>();

                    var identity = (ClaimsIdentity)context.Principal.Identity;

                    context.HttpContext.User = new ClaimsPrincipal(identity);

                    var dto = new CreateOrUpdateCustomerDto(
                        identity.GetSub(),
                        identity.Name,
                        identity.Claims.FirstOrDefault(y => y.Type.Contains("givenname") || y.Type == "given_name")?.Value,
                        identity.Claims.FirstOrDefault(y => y.Type.Contains("surname") || y.Type == "family_name")?.Value);

                    await customerService.CreateOrUpdateCustomerAsync(dto);
                };

                options.Events.OnAuthenticationFailed = context =>
                {
                    // The OWIN version had no AuthenticationFailed notification; keep default
                    // behaviour (the exception surfaces to the error pipeline).
                    return Task.CompletedTask;
                };
            });
        }
    }
}
