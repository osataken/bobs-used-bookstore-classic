using Microsoft.Extensions.Configuration;
using System;
using System.Collections.Generic;

namespace BobsBookstoreClassic.Data
{
    /// <summary>
    /// Application settings and connection strings.
    ///
    /// The .NET Framework version of this class read every entry of the appSettings and
    /// connectionStrings sections of Web.config (allowing an environment variable of the
    /// same name to override an app setting) and exposed them through the static accessors
    /// below.
    ///
    /// On modern .NET the values come from IConfiguration (appsettings.json, environment
    /// variables, parameter store values added at run time through AddSetting, ...). The host
    /// calls <see cref="Initialize(IConfiguration)"/> once during start up; the environment
    /// variable precedence of the original implementation is preserved.
    ///
    /// Host requirements: Initialize must be called (e.g.
    /// BookstoreConfiguration.Initialize(builder.Configuration) in Program.cs) before any
    /// setting is read, otherwise GetSetting throws KeyNotFoundException. The setting names
    /// used by the application are flat (for example "Services/Database"), so they must stay
    /// top-level properties in appsettings.json — nesting them under a section would change
    /// their configuration key to "Section:Name" and the lookups would miss.
    /// </summary>
    public sealed class BookstoreConfiguration
    {
        private static readonly Lazy<BookstoreConfiguration> Lazy = new Lazy<BookstoreConfiguration>(() => new BookstoreConfiguration());

        private static BookstoreConfiguration Instance => Lazy.Value;

        private readonly Dictionary<string, string> _appSettings = new Dictionary<string, string>();
        private readonly Dictionary<string, string> _connectionStrings = new Dictionary<string, string>();

        private BookstoreConfiguration() { }

        /// <summary>
        /// Loads the application settings and connection strings from the supplied
        /// configuration. Must be called once during application start up, before any
        /// setting is read. Calling it again merges in (and refreshes) the values.
        /// </summary>
        public static void Initialize(IConfiguration configuration)
        {
            if (configuration == null) throw new ArgumentNullException(nameof(configuration));

            var instance = Instance;

            const string connectionStringsSection = "ConnectionStrings";

            foreach (var entry in configuration.AsEnumerable())
            {
                if (entry.Value == null) continue;

                var key = entry.Key;

                if (key.StartsWith(connectionStringsSection + ":", StringComparison.OrdinalIgnoreCase))
                {
                    continue;
                }

                instance._appSettings[key] = entry.Value;

                // An environment variable of the same name takes precedence over the
                // configured app setting, matching the original Web.config behaviour.
                var environmentValue = Environment.GetEnvironmentVariable(key);

                if (environmentValue != null)
                {
                    instance._appSettings[key] = environmentValue;
                }
            }

            foreach (var connectionString in configuration.GetSection(connectionStringsSection).GetChildren())
            {
                if (connectionString.Value == null) continue;

                instance._connectionStrings[connectionString.Key] = connectionString.Value;
            }
        }

        public static void AddSetting(string key, string value)
        {
            Instance._appSettings[key] = value;
        }

        public static string GetSetting(string key)
        {
            return Instance._appSettings[key];
        }

        public static T GetSetting<T>(string key)
        {
            var value = Instance._appSettings[key];

            return (T)Convert.ChangeType(value, typeof(T));
        }

        public static void AddConnectionString(string key, string value)
        {
            Instance._connectionStrings[key] = value;
        }

        public static string GetConnectionString(string key)
        {
            return Instance._connectionStrings[key];
        }
    }
}
