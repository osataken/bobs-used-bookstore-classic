using System.Threading.Tasks;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Http;
using NLog;

namespace Bookstore.Web.Infrastructure
{
    /// <summary>
    /// Replaces Global.asax's <c>Application_Error</c>: every unhandled exception is still written to
    /// NLog at Error level. Returning <c>false</c> lets <c>UseExceptionHandler("/Home/Error")</c>
    /// continue to render the error page, exactly as the global HandleErrorAttribute did.
    /// </summary>
    public class NLogExceptionHandler : IExceptionHandler
    {
        private static readonly Logger Logger = LogManager.GetCurrentClassLogger();

        public ValueTask<bool> TryHandleAsync(HttpContext httpContext, System.Exception exception, System.Threading.CancellationToken cancellationToken)
        {
            Logger.Error(exception);

            return ValueTask.FromResult(false);
        }
    }
}
