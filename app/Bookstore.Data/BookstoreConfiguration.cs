using Microsoft.Extensions.Configuration;
using System;
using System.Collections.Generic;

namespace BobsBookstoreClassic.Data
{
    /// <summary>
    /// Static settings façade used across the application. On .NET Framework this read directly from the
    /// legacy framework configuration API; on modern .NET the host supplies an <see cref="IConfiguration"/>
    /// by calling <see cref="Initialize(IConfiguration)"/> during startup, before any setting is read.
    /// </summary>
    public sealed class BookstoreConfiguration
    {
        private static readonly Lazy<BookstoreConfiguration> Lazy = new Lazy<BookstoreConfiguration>(() => new BookstoreConfiguration());

        private static BookstoreConfiguration Instance => Lazy.Value;

        private readonly Dictionary<string, string> _appSettings = new Dictionary<string, string>();
        private readonly Dictionary<string, string> _connectionStrings = new Dictionary<string, string>();

        private BookstoreConfiguration() { }

        /// <summary>
        /// Loads app settings and connection strings from the host configuration. Mirrors the original
        /// behaviour: every setting may be overridden by an environment variable of the same name.
        /// </summary>
        public static void Initialize(IConfiguration configuration)
        {
            if (configuration == null) throw new ArgumentNullException(nameof(configuration));

            var instance = Instance;

            var appSettingsSection = configuration.GetSection("AppSettings");

            // "AppSettings" keeps the web.config <appSettings> keys together when one exists; when it
            // does not, the configuration root itself carries the keys.
            var settingsSource = appSettingsSection.Exists() ? appSettingsSection : configuration;

            foreach (var child in settingsSource.GetChildren())
            {
                if (string.Equals(child.Key, "ConnectionStrings", StringComparison.OrdinalIgnoreCase)) continue;
                if (child.Value == null) continue;

                var value = Environment.GetEnvironmentVariable(child.Key) ?? child.Value;

                instance._appSettings[child.Key] = value;
            }

            foreach (var connectionString in configuration.GetSection("ConnectionStrings").GetChildren())
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
