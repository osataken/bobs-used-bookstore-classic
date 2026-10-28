using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.Areas.Admin.Controllers
{
    [AllowAnonymous]
    public class ErrorController : AdminAreaControllerBase
    {
        // Preserves the original MVC5 absolute route while keeping the conventional
        // area URL (/Admin/Error/Index) reachable.
        [Route("/Error/Index/{code:int}")]
        [Route("/Error/Index")]
        [Route("/Admin/Error/Index")]
        public IActionResult Index(int code = 0)
        {
            var exception = HttpContext.Features.Get<IExceptionHandlerPathFeature>();
            ViewData["Path"] = exception?.Path;
            ViewData["StatusCode"] = code != 0 ? code : HttpContext.Response.StatusCode;

            return View();
        }

        [Route("/error")]
        [Route("/Admin/Error/Support")]
        public IActionResult Support()
        {
            var exception = HttpContext.Features.Get<IExceptionHandlerPathFeature>();
            ViewData["Path"] = exception?.Path;
            ViewData["StatusCode"] = HttpContext.Response.StatusCode;

            return View("~/Areas/Admin/Views/Error/Index.cshtml");
        }
    }
}
