using BobsBookstoreClassic.Data;
using Bookstore.Common;
using NLog;
using NLog.AWS.Logger;
using NLog.Config;
using NLog.Targets;

namespace Bookstore.Web
{
    public static class LoggingSetup
    {
        /// <summary>
        /// Builds the NLog configuration programmatically, preserving the original target/rule
        /// behavior: an AWS CloudWatch Logs target when Services/LoggingService is "aws",
        /// otherwise a debugger/console target, with a single "*" rule at Info level.
        /// Must be called before the first log write (and therefore before the host is built).
        /// </summary>
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

            // In the Kestrel/container hosting model the debugger target is invisible outside a
            // debugger session, so console output is added alongside it to keep the non-AWS
            // configuration observable (stdout is what the local/Docker run surfaces).
            if (!(loggingTarget is AWSTarget))
            {
                var consoleTarget = new ConsoleTarget("console")
                {
                    Layout = "${longdate}|${level:uppercase=true}|${logger}|${message}${onexception:${newline}${exception:format=tostring}}"
                };

                config.AddTarget(consoleTarget);
                config.LoggingRules.Add(new LoggingRule("*", LogLevel.Info, consoleTarget));
            }

            LogManager.Configuration = config;
        }
    }
}
