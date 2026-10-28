using System;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Logging;

namespace Bookstore.Web.Middleware
{
    /// <summary>
    /// Replaces <c>Global.asax.cs Application_Error</c>: logs the unhandled exception (previously
    /// via <c>LogManager.GetCurrentClassLogger().Error(ex)</c>, now through ILogger with NLog as
    /// the provider) and re-throws so that the outer <c>UseExceptionHandler</c> renders the error
    /// view - the behaviour the MVC5 <c>HandleErrorAttribute</c> used to provide.
    /// </summary>
    public class ExceptionLoggingMiddleware
    {
        private readonly RequestDelegate next;
        private readonly ILogger<ExceptionLoggingMiddleware> logger;

        public ExceptionLoggingMiddleware(RequestDelegate next, ILogger<ExceptionLoggingMiddleware> logger)
        {
            this.next = next;
            this.logger = logger;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            try
            {
                await next(context);
            }
            catch (Exception ex)
            {
                logger.LogError(ex, "Unhandled exception processing {Method} {Path}", context.Request.Method, context.Request.Path);

                throw;
            }
        }
    }
}
