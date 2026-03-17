using System;
using Microsoft.AspNetCore.Http;

namespace Bookstore.Web.Helpers
{
    public static class HttpContextExtensions
    {
        public static Guid GetShoppingCartCorrelationId(this HttpContext context)
        {
            const string key = "ShoppingCartCorrelationId";
            
            if (context.Session.TryGetValue(key, out var bytes))
            {
                return new Guid(bytes);
            }
            
            var correlationId = Guid.NewGuid();
            context.Session.Set(key, correlationId.ToByteArray());
            return correlationId;
        }
    }
}
