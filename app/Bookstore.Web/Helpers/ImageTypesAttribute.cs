using Microsoft.AspNetCore.Http;
using System;
using System.ComponentModel.DataAnnotations;
using System.IO;
using System.Linq;

namespace Bookstore.Web.Helpers
{
    public class ImageTypesAttribute : ValidationAttribute
    {
        private readonly string[] imageTypes;

        public ImageTypesAttribute(string[] imageTypes)
        {
            this.imageTypes = imageTypes;
        }

        protected override ValidationResult IsValid(object value, ValidationContext validationContext)
        {
            if (value == null) return ValidationResult.Success;

            if (!(value is IFormFile file)) return ValidationResult.Success;

            var extension = Path.GetExtension(file.FileName);

            if (imageTypes.Contains(extension, StringComparer.OrdinalIgnoreCase))
            {
                return ValidationResult.Success;
            }

            return new ValidationResult(FormatErrorMessage(validationContext?.DisplayName ?? validationContext?.MemberName));
        }

        public override string FormatErrorMessage(string name)
        {
            return $"{name} must be a PNG or JPG image.";
        }
    }
}
