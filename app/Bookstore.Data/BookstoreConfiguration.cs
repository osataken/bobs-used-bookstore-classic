using Microsoft.Extensions.Configuration;
using System;
using System.Collections.Generic;

namespace BobsBookstoreClassic.Data
{
    public sealed class BookstoreConfiguration
    {
        private static readonly Lazy<BookstoreConfiguration> Lazy = new Lazy<BookstoreConfiguration>(() => new BookstoreConfiguration());

        private static BookstoreConfiguration Instance => Lazy.Value;

        private static IConfiguration _configuration;

        private readonly Dictionary<string, string> _appSettings = new Dictionary<string, string>();
        private readonly Dictionary<string, string> _connectionStrings = new Dictionary<string, string>();

        /// <summary>
        /// Loads the application settings and connection strings from the host's
        /// <see cref="IConfiguration"/>. This replaces the .NET Framework
        /// app-settings/connection-string configuration load and must be called once during
        /// host startup, before any setting is read.
        /// </summary>
        public static void Initialize(IConfiguration configuration)
        {
            _configuration = configuration ?? throw new ArgumentNullException(nameof(configuration));

            Instance.Load(configuration);
        }

        private BookstoreConfiguration()
        {
            if (_configuration != null) Load(_configuration);
        }

        private void Load(IConfiguration configuration)
        {
            // AppSettings previously came from web.config's <appSettings> section; they now live
            // under the "AppSettings" section of appsettings.json. Keys are kept verbatim
            // (including the "Services/Database" style names used throughout the application).
            foreach (var setting in configuration.GetSection("AppSettings").AsEnumerable(true))
            {
                if (setting.Value == null) continue;

                _appSettings[setting.Key] = setting.Value;

                // An environment variable of the same name takes precedence over the configured
                // value, exactly as it did in the .NET Framework implementation.
                var environmentValue = Environment.GetEnvironmentVariable(setting.Key);

                if (environmentValue != null)
                {
                    _appSettings[setting.Key] = environmentValue;
                }
            }

            foreach (var connectionString in configuration.GetSection("ConnectionStrings").AsEnumerable(true))
            {
                if (connectionString.Value == null) continue;

                _connectionStrings[connectionString.Key] = connectionString.Value;
            }
        }

        public static void AddSetting(string key, string value)
        {
            Instance._appSettings[key] = value;
        }

        public static string GetSetting(string key)
        {
            if (Instance._appSettings.TryGetValue(key, out var value)) return value;

            // Settings supplied purely through the environment remain resolvable, matching the
            // environment-variable precedence of the original implementation.
            var environmentValue = Environment.GetEnvironmentVariable(key);

            if (environmentValue != null) return environmentValue;

            throw new KeyNotFoundException($"The configuration setting '{key}' was not found.");
        }

        public static T GetSetting<T>(string key)
        {
            var value = GetSetting(key);

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
