using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.Areas.Admin.Controllers
{
    [AllowAnonymous]
    public class ErrorController : AdminAreaControllerBase
    {
        public IActionResult Index(int code)
        {
            var exception = HttpContext.Features.Get<IExceptionHandlerPathFeature>();

            ViewData["Path"] = exception?.Path;
            ViewData["StatusCode"] = code;

            return View();
        }

        public IActionResult Support()
        {
            var exception = HttpContext.Features.Get<IExceptionHandlerPathFeature>();

            ViewData["Path"] = exception?.Path;
            ViewData["StatusCode"] = HttpContext.Response?.StatusCode;

            return View("~/Views/Error/Index.cshtml");
        }
    }
}
