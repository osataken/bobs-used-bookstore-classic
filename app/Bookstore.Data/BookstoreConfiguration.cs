using System;
using System.Collections.Generic;
using Microsoft.Extensions.Configuration;

namespace BobsBookstoreClassic.Data
{
    public sealed class BookstoreConfiguration
    {
        private static readonly Lazy<BookstoreConfiguration> Lazy = new Lazy<BookstoreConfiguration>(() => new BookstoreConfiguration());

        private static BookstoreConfiguration Instance => Lazy.Value;

        private readonly Dictionary<string, string> _appSettings = new Dictionary<string, string>();
        private readonly Dictionary<string, string> _connectionStrings = new Dictionary<string, string>();

        private BookstoreConfiguration() { }

        /// <summary>
        /// Initializes BookstoreConfiguration from ASP.NET Core IConfiguration.
        /// Call this once at startup before using any other methods.
        /// </summary>
        public static void Initialize(IConfiguration configuration)
        {
            // Load all flat configuration values
            foreach (var kvp in configuration.AsEnumerable())
            {
                if (kvp.Value != null)
                {
                    Instance._appSettings[kvp.Key] = kvp.Value;
                }
            }

            // Load connection strings from ConnectionStrings section
            var connectionStringsSection = configuration.GetSection("ConnectionStrings");
            foreach (var item in connectionStringsSection.GetChildren())
            {
                if (item.Value != null)
                    Instance._connectionStrings[item.Key] = item.Value;
            }

            // Allow environment variable overrides using double-underscore separator
            foreach (var key in new List<string>(Instance._appSettings.Keys))
            {
                var envKey = key.Replace("/", "__").Replace(":", "__");
                var envValue = Environment.GetEnvironmentVariable(envKey);
                if (envValue != null)
                {
                    Instance._appSettings[key] = envValue;
                }
            }
        }

        public static void AddSetting(string key, string value)
        {
            Instance._appSettings[key] = value;
        }

        public static string GetSetting(string key)
        {
            if (Instance._appSettings.TryGetValue(key, out var value))
                return value ?? string.Empty;
            return string.Empty;
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
            if (Instance._connectionStrings.TryGetValue(key, out var value))
                return value ?? string.Empty;
            return string.Empty;
        }
    }
}
