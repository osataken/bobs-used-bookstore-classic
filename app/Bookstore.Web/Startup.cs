using Microsoft.AspNetCore.Owin;

using Microsoft.AspNetCore.Builder;




namespace Bookstore.Web
{
    public static class LoggingSetup
    {
        public static void ConfigureLogging()
        {
        }
    }

    public class Startup
    {
        public void Configuration(IApplicationBuilder app)
        {
            LoggingSetup.ConfigureLogging();

            // ConfigurationSetup.ConfigureConfiguration();

            // DependencyInjectionSetup.ConfigureDependencyInjection(app);

            // AuthenticationConfig.ConfigureAuthentication(app);
        }
    }
}
