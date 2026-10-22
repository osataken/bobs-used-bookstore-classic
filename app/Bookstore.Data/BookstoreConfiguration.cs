using System;
using System.Collections.Generic;
using Microsoft.Extensions.Configuration;

namespace BobsBookstoreClassic.Data
{
    public sealed class BookstoreConfiguration
    {
        private static readonly Lazy<BookstoreConfiguration> Lazy = new Lazy<BookstoreConfiguration>(() => new BookstoreConfiguration());

        private static BookstoreConfiguration Instance => Lazy.Value;

        private readonly Dictionary<string, string> _appSettings = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);
        private readonly Dictionary<string, string> _connectionStrings = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);

        private BookstoreConfiguration() { }

        public static void Initialize(IConfiguration configuration)
        {
            // Load all configuration values from IConfiguration
            foreach (var kvp in configuration.AsEnumerable())
            {
                if (kvp.Value == null) continue;

                if (kvp.Key.StartsWith("ConnectionStrings:", StringComparison.OrdinalIgnoreCase))
                {
                    var csName = kvp.Key.Substring("ConnectionStrings:".Length);
                    Instance._connectionStrings[csName] = kvp.Value;
                }
                else
                {
                    Instance._appSettings[kvp.Key] = kvp.Value;
                }

                // Allow environment variables to override (using the same key)
                var envValue = Environment.GetEnvironmentVariable(kvp.Key);
                if (envValue != null)
                {
                    Instance._appSettings[kvp.Key] = envValue;
                }
            }
        }

        public static void AddSetting(string key, string value)
        {
            Instance._appSettings[key] = value;
        }

        public static string GetSetting(string key)
        {
            return Instance._appSettings.TryGetValue(key, out var value) ? value : null;
        }

        public static T GetSetting<T>(string key)
        {
            var value = GetSetting(key);
            if (value == null) return default;
            return (T)Convert.ChangeType(value, typeof(T));
        }

        public static void AddConnectionString(string key, string value)
        {
            Instance._connectionStrings[key] = value;
        }

        public static string GetConnectionString(string key)
        {
            return Instance._connectionStrings.TryGetValue(key, out var value) ? value : null;
        }
    }
}
