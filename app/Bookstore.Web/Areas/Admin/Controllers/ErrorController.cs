using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.Areas.Admin.Controllers
{
    [AllowAnonymous]
    public class ErrorController : AdminAreaControllerBase
    {
        // Both the two-segment absolute path (used by the Program.cs error wires)
        // and the MVC5 area URL (/Admin/Error/Index?code=404) must resolve.
        [Route("/Error/Index/{code:int}")]
        [Route("/Error/Index")]
        [Route("/Admin/Error/Index/{code:int?}")]
        public IActionResult Index(int code)
        {
            ViewData["StatusCode"] = code;

            return View();
        }

        [Route("/error")]
        [Route("/Admin/Error/Support")]
        public IActionResult Support()
        {
            return View("~/Areas/Admin/Views/Error/Index.cshtml");
        }
    }
}
