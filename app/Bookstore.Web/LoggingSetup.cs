using BobsBookstoreClassic.Data;
using Bookstore.Common;
using NLog;
using NLog.AWS.Logger;
using NLog.Config;
using NLog.Targets;

namespace Bookstore.Web
{
    /// <summary>
    /// Programmatic NLog configuration, ported from App_Start/LoggingSetup.cs. NLog remains the
    /// logging provider, but is now hosted behind Microsoft.Extensions.Logging via NLog.Web.
    /// </summary>
    public static class LoggingSetup
    {
        public static void ConfigureLogging()
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

            config.LoggingRules.Add(new LoggingRule("*", NLog.LogLevel.Info, loggingTarget));

            LogManager.Configuration = config;
        }
    }
}
