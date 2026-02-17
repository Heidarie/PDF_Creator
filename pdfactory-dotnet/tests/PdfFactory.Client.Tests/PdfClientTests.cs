using System.Net;
using System.Text;
using System.Text.Json;
using PdfFactory;
using Xunit;

namespace PdfFactory.Client.Tests;

public sealed class PdfClientTests
{
    [Fact]
    public async Task RenderAsync_ReturnsBytes_OnSuccess()
    {
        var handler = new FakeHandler((_, _) =>
        {
            var response = new HttpResponseMessage(HttpStatusCode.OK)
            {
                Content = new ByteArrayContent(Encoding.UTF8.GetBytes("%PDF-1.4"))
            };
            response.Content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("application/pdf");
            return response;
        });

        var client = new PdfClient(new HttpClient(handler) { BaseAddress = new Uri("http://localhost:8080/") });
        var bytes = await client.RenderAsync(new RenderRequest { Html = "<p>hi</p>" });

        Assert.StartsWith("%PDF-", Encoding.UTF8.GetString(bytes));
    }

    [Fact]
    public async Task RenderAsync_Throws_OnErrorResponse()
    {
        var handler = new FakeHandler((_, _) =>
        {
            var payload = JsonSerializer.SerializeToUtf8Bytes(new ErrorResponse { Error = "bad input" });
            return new HttpResponseMessage(HttpStatusCode.BadRequest)
            {
                Content = new ByteArrayContent(payload)
            };
        });

        var client = new PdfClient(new HttpClient(handler) { BaseAddress = new Uri("http://localhost:8080/") });

        var ex = await Assert.ThrowsAsync<InvalidOperationException>(() =>
            client.RenderAsync(new RenderRequest { Html = "<p>hi</p>" }));

        Assert.Contains("bad input", ex.Message, StringComparison.OrdinalIgnoreCase);
    }
}

internal sealed class FakeHandler : HttpMessageHandler
{
    private readonly Func<HttpRequestMessage, CancellationToken, HttpResponseMessage> _handler;

    public FakeHandler(Func<HttpRequestMessage, CancellationToken, HttpResponseMessage> handler)
    {
        _handler = handler;
    }

    protected override Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
    {
        return Task.FromResult(_handler(request, cancellationToken));
    }
}
