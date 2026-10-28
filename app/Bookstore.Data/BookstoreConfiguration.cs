using Microsoft.Extensions.Configuration;
using System;
using System.Collections.Generic;

namespace BobsBookstoreClassic.Data
{
    /// <summary>
    /// Static configuration facade. Originally backed by ConfigurationManager (app/web.config);
    /// now backed by Microsoft.Extensions.Configuration. The host supplies an IConfiguration via
    /// <see cref="Initialize(IConfiguration)"/>; when it does not, a default configuration built
    /// from appsettings.json plus environment variables is used, so behavior is preserved for
    /// every existing call site.
    /// </summary>
    public sealed class BookstoreConfiguration
    {
        private static readonly object SyncRoot = new object();
        private static BookstoreConfiguration _instance;

        private static BookstoreConfiguration Instance
        {
            get
            {
                if (_instance != null) return _instance;

                lock (SyncRoot)
                {
                    if (_instance == null) _instance = new BookstoreConfiguration(BuildDefaultConfiguration());
                }

                return _instance;
            }
        }

        private readonly Dictionary<string, string> _appSettings = new Dictionary<string, string>();
        private readonly Dictionary<string, string> _connectionStrings = new Dictionary<string, string>();

        private BookstoreConfiguration(IConfiguration configuration)
        {
            if (configuration == null) return;

            // Settings. Configuration keys are hierarchical with ':' separators; the application
            // addresses them with '/' separators (e.g. "Files/BucketName"), as the original
            // appSettings keys did, so both forms are registered. Environment variables keep the
            // override precedence the ConfigurationManager-based loop implemented.
            foreach (var pair in configuration.AsEnumerable())
            {
                if (pair.Value == null) continue;

                var key = pair.Key;

                if (key.StartsWith("ConnectionStrings:", StringComparison.OrdinalIgnoreCase)) continue;

                var slashKey = key.Replace(':', '/');

                _appSettings[key] = pair.Value;
                _appSettings[slashKey] = pair.Value;

                var environmentValue = Environment.GetEnvironmentVariable(slashKey) ?? Environment.GetEnvironmentVariable(key);

                if (environmentValue != null)
                {
                    _appSettings[key] = environmentValue;
                    _appSettings[slashKey] = environmentValue;
                }
            }

            // Connection strings.
            foreach (var child in configuration.GetSection("ConnectionStrings").GetChildren())
            {
                if (child.Value != null) _connectionStrings[child.Key] = child.Value;
            }
        }

        private static IConfiguration BuildDefaultConfiguration()
        {
            return new ConfigurationBuilder()
                .AddJsonFile("appsettings.json", optional: true, reloadOnChange: false)
                .AddEnvironmentVariables()
                .Build();
        }

        /// <summary>
        /// Supplies the host's configuration. Call once during startup, before the first
        /// GetSetting/GetConnectionString call.
        /// </summary>
        public static void Initialize(IConfiguration configuration)
        {
            if (configuration == null) throw new ArgumentNullException(nameof(configuration));

            lock (SyncRoot)
            {
                _instance = new BookstoreConfiguration(configuration);
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
