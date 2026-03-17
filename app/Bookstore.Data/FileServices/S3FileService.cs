using Amazon.S3;
using Amazon.S3.Model;
using Amazon.S3.Transfer;
using Bookstore.Domain;
using Microsoft.Extensions.Configuration;
using System.IO;
using System.Threading.Tasks;

namespace Bookstore.Data.FileServices
{
    public class S3FileService : IFileService
    {
        private readonly TransferUtility transferUtility;
        private readonly IConfiguration _configuration;

        public S3FileService(IAmazonS3 s3Client, IConfiguration configuration)
        {
            transferUtility = new TransferUtility(s3Client);
            _configuration = configuration;
        }

        public async Task DeleteAsync(string filePath)
        {
            if (string.IsNullOrWhiteSpace(filePath)) return;

            var request = new DeleteObjectRequest
            {
                BucketName = _configuration["Files:BucketName"],
                Key = Path.GetFileName(filePath)
            };

            await transferUtility.S3Client.DeleteObjectAsync(request);
        }

        public async Task<string> SaveAsync(Stream contents, string filename)
        {
            if (contents == null) return null;

            var uniqueFilename = $"{Path.GetFileNameWithoutExtension(Path.GetRandomFileName())}{Path.GetExtension(filename)}";

            var request = new TransferUtilityUploadRequest
            {
                BucketName = _configuration["Files:BucketName"],
                InputStream = contents,
                Key = uniqueFilename
            };

            await transferUtility.UploadAsync(request);

            return $"{_configuration["Files:CloudFrontDomain"]}/{uniqueFilename}";
        }
    }
}
