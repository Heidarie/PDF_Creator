namespace PdfKit;

public interface IPdfClient
{
    Task<byte[]> RenderAsync(RenderRequest request, CancellationToken cancellationToken = default);
    Task RenderToFileAsync(RenderRequest request, string path, CancellationToken cancellationToken = default);
}
