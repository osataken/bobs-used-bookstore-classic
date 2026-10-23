using Microsoft.AspNetCore.Http;
using System.ComponentModel.DataAnnotations;

namespace Bookstore.Web.Helpers
{
    public class MaxFileSizeAttribute : ValidationAttribute
    {
        private readonly int maxFileSize;

        public MaxFileSizeAttribute(int maxFileSize)
        {
            this.maxFileSize = maxFileSize;
        }

        protected override ValidationResult IsValid(object value, ValidationContext validationContext)
        {
            if (value == null) return ValidationResult.Success;

            if (!(value is IFormFile file)) return ValidationResult.Success;

            if (file.Length <= maxFileSize) return ValidationResult.Success;

            return new ValidationResult(FormatErrorMessage(validationContext?.DisplayName ?? validationContext?.MemberName));
        }

        public override string FormatErrorMessage(string name)
        {
            return $"{name} cannot exceed {maxFileSize.ToStorageSize()}";
        }
    }
}
