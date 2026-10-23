using System;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Logging;

namespace Bookstore.Web
{
    /// <summary>
    /// Replacement for <c>Global.asax.cs</c> <c>Application_Error</c>, which logged every
    /// unhandled exception through NLog. The exception is logged via ILogger (NLog is still the
    /// provider behind it) and then re-thrown so that UseExceptionHandler renders the error view,
    /// preserving the behaviour of the MVC5 HandleErrorAttribute.
    /// </summary>
    public class ExceptionLoggingMiddleware : IMiddleware
    {
        private readonly ILogger<ExceptionLoggingMiddleware> _logger;

        public ExceptionLoggingMiddleware(ILogger<ExceptionLoggingMiddleware> logger)
        {
            _logger = logger;
        }

        public async Task InvokeAsync(HttpContext context, RequestDelegate next)
        {
            try
            {
                await next(context);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Unhandled exception processing {Method} {Path}", context.Request.Method, context.Request.Path);

                throw;
            }
        }
    }
}
