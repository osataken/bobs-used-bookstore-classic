using System.Threading.Tasks;
using Amazon.SimpleSystemsManagement;
using Amazon.SimpleSystemsManagement.Model;
using BobsBookstoreClassic.Data;
using Bookstore.Common;

namespace Bookstore.Web
{
    /// <summary>
    /// Loads configuration from AWS Systems Manager Parameter Store, ported from
    /// App_Start/ConfigurationSetup.cs. The modern AWS SDK only exposes asynchronous SSM
    /// operations, so this method is async and is awaited from Program.cs before the host is
    /// built (no sync-over-async blocking).
    /// </summary>
    public static class ConfigurationSetup
    {
        public static async Task ConfigureConfigurationAsync()
        {
            var rootPath = "/" + Constants.AppName;

            const string databasePath = "/Database";
            const string authenticationPath = "/Authentication";
            const string fileServicePath = "/Files";

            if (BookstoreConfiguration.GetSetting("Services/Database") == "aws")
            {
                using (var client = new AmazonSimpleSystemsManagementClient())
                {
                    var request = new GetParameterRequest { Name = $"{rootPath}{databasePath}/ConnectionStrings/BookstoreDatabaseConnection" };
                    var response = await client.GetParameterAsync(request);

                    BookstoreConfiguration.AddSetting(response.Parameter.Name.Replace($"{rootPath}{databasePath}/", string.Empty), response.Parameter.Value);

                    // The original App_Start/ConfigurationSetup.cs only called AddSetting here, which
                    // stored the value under "ConnectionStrings/BookstoreDatabaseConnection" in the
                    // app-settings bag. DependencyInjectionSetup resolves the context with
                    // GetConnectionString("BookstoreDatabaseConnection"), which reads the separate
                    // connection-string bag, so on the Services/Database=aws path the SSM value was
                    // never actually consumed. Publish it where the consumer looks for it.
                    BookstoreConfiguration.AddConnectionString("BookstoreDatabaseConnection", response.Parameter.Value);
                }
            }

            if (BookstoreConfiguration.GetSetting("Services/Authentication") == "aws")
            {
                using (var client = new AmazonSimpleSystemsManagementClient())
                {
                    var request = new GetParametersByPathRequest { Path = $"{rootPath}{authenticationPath}/", Recursive = true };
                    var response = await client.GetParametersByPathAsync(request);

                    foreach (var parameter in response.Parameters)
                    {
                        BookstoreConfiguration.AddSetting(parameter.Name.Replace($"{rootPath}/", string.Empty), parameter.Value);
                    }
                }
            }

            if (BookstoreConfiguration.GetSetting("Services/FileService") == "aws")
            {
                using (var client = new AmazonSimpleSystemsManagementClient())
                {
                    var request = new GetParametersByPathRequest { Path = $"{rootPath}{fileServicePath}/", Recursive = true };
                    var response = await client.GetParametersByPathAsync(request);

                    foreach (var parameter in response.Parameters)
                    {
                        BookstoreConfiguration.AddSetting(parameter.Name.Replace($"{rootPath}/", string.Empty), parameter.Value);
                    }
                }
            }
        }
    }
}
