using System;
using System.Collections.Generic;

namespace BobsBookstoreClassic.Data
{
    /// <summary>
    /// Static accessor for application settings and connection strings.
    /// Values are populated at startup from <c>IConfiguration</c> by the web host
    /// (see Bookstore.Web/Program.cs) via <see cref="AddSetting"/> and
    /// <see cref="AddConnectionString"/>. Environment variables take precedence on read,
    /// preserving the override behavior the legacy ConfigurationManager loader provided.
    /// </summary>
    public sealed class BookstoreConfiguration
    {
        private static readonly Lazy<BookstoreConfiguration> Lazy = new Lazy<BookstoreConfiguration>(() => new BookstoreConfiguration());

        private static BookstoreConfiguration Instance => Lazy.Value;

        private readonly Dictionary<string, string> _appSettings = new Dictionary<string, string>();
        private readonly Dictionary<string, string> _connectionStrings = new Dictionary<string, string>();

        private BookstoreConfiguration()
        {
        }

        public static void AddSetting(string key, string value)
        {
            Instance._appSettings[key] = value;
        }

        public static string GetSetting(string key)
        {
            // Environment variables override configured settings (legacy behavior).
            var envVar = Environment.GetEnvironmentVariable(key);
            if (envVar != null)
            {
                return envVar;
            }

            return Instance._appSettings.TryGetValue(key, out var value) ? value : string.Empty;
        }

        public static T GetSetting<T>(string key)
        {
            var value = GetSetting(key);

            if (string.IsNullOrEmpty(value))
            {
                return default!;
            }

            return (T)Convert.ChangeType(value, typeof(T));
        }

        public static void AddConnectionString(string key, string value)
        {
            Instance._connectionStrings[key] = value;
        }

        public static string GetConnectionString(string key)
        {
            var envVar = Environment.GetEnvironmentVariable(key);
            if (envVar != null)
            {
                return envVar;
            }

            return Instance._connectionStrings.TryGetValue(key, out var value) ? value : string.Empty;
        }
    }
}
