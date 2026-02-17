using System.Text;

namespace PdfFactory.Builders;

public static class DataUrl
{
    public static string FromSvg(string svg) => FromSvg(Encoding.UTF8.GetBytes(svg));

    public static string FromSvg(byte[] svg) =>
        "data:image/svg+xml;base64," + Convert.ToBase64String(svg);

    public static string FromPng(byte[] png) =>
        "data:image/png;base64," + Convert.ToBase64String(png);

    public static string FromJpeg(byte[] jpg) =>
        "data:image/jpeg;base64," + Convert.ToBase64String(jpg);
}
