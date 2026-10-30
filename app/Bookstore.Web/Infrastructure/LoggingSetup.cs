using BobsBookstoreClassic.Data;
using Bookstore.Common;
using NLog;
using NLog.AWS.Logger;
using NLog.Config;
using NLog.Targets;

namespace Bookstore.Web.Infrastructure
{
    /// <summary>
    /// Ported from App_Start/LoggingSetup.cs. Builds the NLog configuration programmatically,
    /// preserving the AWS CloudWatch vs Debugger target switch and the "*"/Info logging rule.
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

            config.LoggingRules.Add(new LoggingRule("*", LogLevel.Info, loggingTarget));

            LogManager.Configuration = config;
        }
    }
}
