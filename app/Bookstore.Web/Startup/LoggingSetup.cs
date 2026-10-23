using BobsBookstoreClassic.Data;
using Bookstore.Common;
using Microsoft.AspNetCore.Builder;
using NLog;
using NLog.AWS.Logger;
using NLog.Config;
using NLog.Targets;
using NLog.Web;

namespace Bookstore.Web
{
    /// <summary>
    /// Port of the .NET Framework App_Start/LoggingSetup.cs. The NLog configuration is still
    /// built programmatically (AWS CloudWatch target when Services/LoggingService == "aws",
    /// debugger target otherwise, single "*" rule at Info) and is additionally hooked up as the
    /// Microsoft.Extensions.Logging provider for the host.
    /// </summary>
    public static class LoggingSetup
    {
        public static void ConfigureLogging(WebApplicationBuilder builder)
        {
            var config = new LoggingConfiguration();

            Target loggingTarget;

            if (BookstoreConfiguration.GetSetting("Services/LoggingService") == "aws")
            {
                loggingTarget = new AWSTarget { LogGroup = Constants.AppName };
            }
            else
            {
                loggingTarget = new DebuggerTarget();
            }

            config.AddTarget("aws", loggingTarget);

            config.LoggingRules.Add(new LoggingRule("*", LogLevel.Info, loggingTarget));

            LogManager.Configuration = config;

            // Route Microsoft.Extensions.Logging (ILogger<T>) through the same NLog configuration.
            builder.Host.UseNLog(new NLogAspNetCoreOptions { ReplaceLoggerFactory = true });
        }
    }
}
