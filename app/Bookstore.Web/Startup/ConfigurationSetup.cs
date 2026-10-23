using System.Threading.Tasks;
using Amazon.SimpleSystemsManagement;
using Amazon.SimpleSystemsManagement.Model;
using BobsBookstoreClassic.Data;
using Bookstore.Common;

namespace Bookstore.Web
{
    /// <summary>
    /// Port of the .NET Framework App_Start/ConfigurationSetup.cs. Loads the AWS Systems Manager
    /// Parameter Store values for the Database, Authentication and File service paths into
    /// <see cref="BookstoreConfiguration"/>, keeping the three "aws" conditionals and the same
    /// key-stripping behaviour. AWS SDK v4 only exposes the asynchronous Parameter Store APIs,
    /// so the step is awaited during start up (before the DI container is built).
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

                    // The original application read this value back through
                    // BookstoreConfiguration.GetConnectionString("BookstoreDatabaseConnection"),
                    // which on .NET Framework resolved against the connectionStrings section.
                    BookstoreConfiguration.AddConnectionString(
                        response.Parameter.Name.Replace($"{rootPath}{databasePath}/ConnectionStrings/", string.Empty),
                        response.Parameter.Value);
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
