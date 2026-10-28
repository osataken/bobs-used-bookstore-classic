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
    /// ASP.NET Core replacement for the OWIN AuthenticationConfig.ConfigureAuthentication(IAppBuilder).
    /// Cognito/OIDC when Services/Authentication == "aws", otherwise cookie authentication whose
    /// identity is supplied by <see cref="LocalAuthenticationMiddleware"/>.
    /// </summary>
    public static class AuthenticationSetup
    {
        public static void Configure(WebApplicationBuilder builder)
        {
            if (IsAwsAuthentication())
            {
                ConfigureCognitoAuthentication(builder);
            }
            else
            {
                ConfigureLocalAuthentication(builder);
            }
        }

        public static bool IsAwsAuthentication()
        {
            return BookstoreConfiguration.GetSetting("Services/Authentication") == "aws";
        }

        private static void ConfigureLocalAuthentication(WebApplicationBuilder builder)
        {
            // The local development identity is created per-request by LocalAuthenticationMiddleware
            // (registered in Program.cs). Cookie authentication is still registered so that the
            // authentication/authorization middleware has a valid default scheme.
            builder.Services
                .AddAuthentication(CookieAuthenticationDefaults.AuthenticationScheme)
                .AddCookie(options =>
                {
                    // Without these the fallback authorization policy challenges the cookie scheme
                    // and redirects to the ASP.NET Core default /Account/Login, which does not
                    // exist in this app. /Authentication/Login is intercepted by
                    // LocalAuthenticationMiddleware, which mints the dev identity and redirects to "/".
                    options.LoginPath = "/Authentication/Login";
                    options.AccessDeniedPath = "/Home/Error";
                });
        }

        private static void ConfigureCognitoAuthentication(WebApplicationBuilder builder)
        {
            builder.Services.AddAuthentication(options =>
                {
                    // OWIN equivalent: SetDefaultSignInAsAuthenticationType(Cookies) plus
                    // SignInAsAuthenticationType = Cookies on the OpenIdConnect options.
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
                    options.SignInScheme = CookieAuthenticationDefaults.AuthenticationScheme;
                    options.UseTokenLifetime = false;
                    options.SaveTokens = true;

                    // OWIN set Scope = "openid profile" wholesale; the Core handler uses a collection.
                    options.Scope.Clear();
                    options.Scope.Add("openid");
                    options.Scope.Add("profile");

                    options.TokenValidationParameters = new TokenValidationParameters
                    {
                        NameClaimType = "cognito:username",
                        RoleClaimType = "cognito:groups"
                    };

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
                        OnTokenValidated = async context =>
                        {
                            var service = context.HttpContext.RequestServices.GetRequiredService<ICustomerService>();

                            var identity = (ClaimsIdentity)context.Principal.Identity;

                            var dto = new CreateOrUpdateCustomerDto(
                                identity.GetSub(),
                                identity.Name,
                                identity.FindFirst(x => x.Type.Contains("givenname")).Value,
                                identity.FindFirst(x => x.Type.Contains("surname")).Value);

                            await service.CreateOrUpdateCustomerAsync(dto);
                        }
                    };
                });
        }
    }
}
