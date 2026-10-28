using System.Threading.Tasks;
using Amazon.SimpleSystemsManagement;
using Amazon.SimpleSystemsManagement.Model;
using BobsBookstoreClassic.Data;
using Bookstore.Common;

namespace Bookstore.Web
{
    public static class ConfigurationSetup
    {
        /// <summary>
        /// Pulls the Database / Authentication / FileService parameters from AWS Systems Manager
        /// Parameter Store when the matching Services/* setting is "aws" and feeds them into
        /// <see cref="BookstoreConfiguration"/>. The modern AWS SDK only exposes the asynchronous
        /// operations, so the original synchronous calls became awaited calls; the set of
        /// parameters fetched and the keys they are stored under are unchanged.
        /// Must complete before dependency injection / authentication composition reads settings.
        /// </summary>
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
