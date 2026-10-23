using System.Security.Claims;
using System.Threading.Tasks;
using BobsBookstoreClassic.Data;
using Bookstore.Domain.Customers;
using Bookstore.Web.Helpers;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.AspNetCore.Authentication.OpenIdConnect;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.IdentityModel.Protocols.OpenIdConnect;
using Microsoft.IdentityModel.Tokens;

namespace Bookstore.Web
{
    /// <summary>
    /// Port of the .NET Framework App_Start/AuthenticationSetup.cs (OWIN
    /// <c>AuthenticationConfig.ConfigureAuthentication</c> OWIN-era setup). The "aws" branch wires
    /// cookie + OpenIdConnect (Cognito) authentication handlers; any other value keeps the
    /// local development authentication middleware.
    /// </summary>
    public static class AuthenticationSetup
    {
        public static void ConfigureAuthentication(WebApplicationBuilder builder)
        {
            if (IsCognitoAuthentication())
            {
                ConfigureCognitoAuthentication(builder);
            }
            else
            {
                ConfigureLocalAuthentication(builder);
            }
        }

        /// <summary>
        /// The OWIN application had no authentication handler in the "local" service mode —
        /// LocalAuthenticationMiddleware populated the user itself. ASP.NET Core still needs a
        /// registered scheme, otherwise UseAuthentication() and the global authorization fallback
        /// policy throw "No authenticationScheme was specified". A cookie scheme whose LoginPath is
        /// /Authentication/Login reproduces the MVC5 behavior: an unauthenticated request to a
        /// protected action is challenged to /Authentication/Login, which the local middleware
        /// intercepts to sign the development user in.
        /// </summary>
        private static void ConfigureLocalAuthentication(WebApplicationBuilder builder)
        {
            builder.Services
                .AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
                .AddCookie(options =>
                {
                    options.LoginPath = "/Authentication/Login";
                    options.Cookie.HttpOnly = true;
                    options.Cookie.IsEssential = true;
                });
        }

        /// <summary>
        /// Replaces <c>app.UseMiddlewareFromContainer&lt;LocalAuthenticationMiddleware&gt;()</c>.
        /// The middleware is resolved from the Autofac container (registered AsSelf, scoped) by
        /// the factory-middleware support of UseMiddleware for <see cref="Microsoft.AspNetCore.Http.IMiddleware"/>.
        /// </summary>
        public static void UseLocalAuthentication(WebApplication app)
        {
            if (IsCognitoAuthentication()) return;

            app.UseMiddleware<LocalAuthenticationMiddleware>();
        }

        private static bool IsCognitoAuthentication()
        {
            return BookstoreConfiguration.GetSetting("Services/Authentication") == "aws";
        }

        private static void ConfigureCognitoAuthentication(WebApplicationBuilder builder)
        {
            builder.Services
                .AddAuthentication(options =>
                {
                    // OWIN registered the OIDC middleware in Active mode, so it answered the 401
                    // raised by the global AuthorizeAttribute and redirected to Cognito.
                    // DefaultChallengeScheme = OpenIdConnect reproduces that; sign-in and the
                    // resulting principal still live in the cookie scheme.
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

                    options.Events = new OpenIdConnectEvents
                    {
                        // OWIN: Notifications.RedirectToIdentityProvider
                        OnRedirectToIdentityProvider = context =>
                        {
                            context.ProtocolMessage.RedirectUri = context.Request.GetReturnUrl();

                            return Task.CompletedTask;
                        },

                        // OWIN: Notifications.AuthorizationCodeReceived
                        OnAuthorizationCodeReceived = context =>
                        {
                            if (context.TokenEndpointRequest != null)
                            {
                                context.TokenEndpointRequest.RedirectUri = context.Request.GetReturnUrl();
                            }

                            return Task.CompletedTask;
                        },

                        // OWIN: Notifications.SecurityTokenValidated
                        OnTokenValidated = async context =>
                        {
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
