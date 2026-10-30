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

namespace Bookstore.Web.Infrastructure
{
    /// <summary>
    /// Was App_Start/AuthenticationSetup.cs (AuthenticationConfig.ConfigureAuthentication(IAppBuilder)).
    /// The OWIN cookie + OpenIdConnect middleware registrations are now authentication schemes
    /// registered in DI; the local-authentication branch is plugged into the pipeline in Program.cs.
    /// </summary>
    public static class AuthenticationConfig
    {
        /// <summary>True when the app authenticates against Amazon Cognito (Services/Authentication == "aws").</summary>
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

            services.AddAuthorization();
        }

        /// <summary>
        /// Was app.UseMiddlewareFromContainer&lt;LocalAuthenticationMiddleware&gt;(). The middleware itself is
        /// registered with Autofac (DependencyInjectionSetup) and inserted into the pipeline in Program.cs.
        /// A cookie scheme is still registered so the framework has a default scheme to work with.
        /// </summary>
        private static void ConfigureLocalAuthentication(IServiceCollection services)
        {
            services
                .AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
                .AddCookie(options =>
                {
                    // Under OWIN no cookie middleware was registered in local mode, so a denied
                    // request produced a bare 401. The Core cookie handler challenges with a
                    // redirect, whose default target (/Account/Login) does not exist here — point
                    // it at the path LocalAuthenticationMiddleware handles so the loop closes.
                    options.LoginPath = "/Authentication/Login";
                    options.AccessDeniedPath = "/Home/Index";
                });
        }

        private static void ConfigureCognitoAuthentication(IServiceCollection services)
        {
            services
                .AddAuthentication(options =>
                {
                    // Was app.SetDefaultSignInAsAuthenticationType(CookieAuthenticationDefaults.AuthenticationType).
                    options.DefaultScheme = CookieAuthenticationDefaults.AuthenticationScheme;
                    options.DefaultSignInScheme = CookieAuthenticationDefaults.AuthenticationScheme;
                    options.DefaultChallengeScheme = OpenIdConnectDefaults.AuthenticationScheme;
                })
                // Was app.UseCookieAuthentication(new CookieAuthenticationOptions()).
                .AddCookie()
                // Was app.UseOpenIdConnectAuthentication(new OpenIdConnectAuthenticationOptions { ... }).
                .AddOpenIdConnect(options =>
                {
                    options.ClientId = BookstoreConfiguration.GetSetting("Authentication/Cognito/LocalClientId");
                    options.MetadataAddress = BookstoreConfiguration.GetSetting("Authentication/Cognito/MetadataAddress");
                    options.ResponseType = OpenIdConnectResponseType.Code;
                    options.UsePkce = false;
                    // RedeemCode = true on OWIN; on Core the authorization code is redeemed when
                    // the handler is told to get the claims from the token endpoint.
                    options.GetClaimsFromUserInfoEndpoint = false;
                    options.ResponseMode = OpenIdConnectResponseMode.Query;
                    options.Scope.Clear();
                    options.Scope.Add("openid");
                    options.Scope.Add("profile");
                    // Was RedirectUri = "{scheme}://{host}/signin-oidc" (see OnRedirectToIdentityProvider).
                    options.CallbackPath = "/signin-oidc";
                    options.SignInScheme = CookieAuthenticationDefaults.AuthenticationScheme;
                    options.UseTokenLifetime = false;
                    options.SaveTokens = true;
                    options.TokenValidationParameters = new TokenValidationParameters
                    {
                        NameClaimType = "cognito:username",
                        RoleClaimType = "cognito:groups"
                    };

                    // Was OpenIdConnectAuthenticationNotifications.
                    options.Events = new OpenIdConnectEvents
                    {
                        OnRedirectToIdentityProvider = context =>
                        {
                            context.ProtocolMessage.RedirectUri = context.Request.GetReturnUrl();

                            return Task.CompletedTask;
                        },
                        OnAuthorizationCodeReceived = context =>
                        {
                            context.TokenEndpointRequest.RedirectUri = context.Request.GetReturnUrl();

                            return Task.CompletedTask;
                        },
                        // Was SecurityTokenValidated: upsert the customer record on every login.
                        OnTokenValidated = async context =>
                        {
                            // Replaces x.OwinContext.GetAutofacLifetimeScope().Resolve<ICustomerService>().
                            var service = context.HttpContext.RequestServices.GetRequiredService<ICustomerService>();

                            var identity = (ClaimsIdentity)context.Principal.Identity;

                            var dto = new CreateOrUpdateCustomerDto(
                                identity.GetSub(),
                                identity.Name,
                                identity.FindFirst(y => y.Type.Contains("givenname")).Value,
                                identity.FindFirst(y => y.Type.Contains("surname")).Value);

                            await service.CreateOrUpdateCustomerAsync(dto);
                        }
                    };
                });
        }
    }
}
