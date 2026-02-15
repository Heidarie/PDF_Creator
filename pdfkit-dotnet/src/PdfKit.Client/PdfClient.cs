using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace PdfKit;

internal sealed class PdfClient : IPdfClient
{
    private readonly HttpClient _httpClient;
    private readonly JsonSerializerOptions _jsonOptions;
    private readonly long _maxResponseBytes;

    public PdfClient(HttpClient httpClient, JsonSerializerOptions? jsonOptions = null, long maxResponseBytes = 25 * 1024 * 1024)
    {
        _httpClient = httpClient ?? throw new ArgumentNullException(nameof(httpClient));
        _jsonOptions = jsonOptions ?? new JsonSerializerOptions
        {
            DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
            PropertyNamingPolicy = null
        };
        _maxResponseBytes = maxResponseBytes;
    }

    public async Task<byte[]> RenderAsync(RenderRequest request, CancellationToken cancellationToken = default)
    {
        if (request is null)
        {
            throw new ArgumentNullException(nameof(request));
        }

        if (string.IsNullOrWhiteSpace(request.Html))
        {
            throw new ArgumentException("Html is required", nameof(request));
        }

        var payload = JsonSerializer.Serialize(request, _jsonOptions);
        using var content = new StringContent(payload, Encoding.UTF8, "application/json");
        using var response = await _httpClient.PostAsync("render", content, cancellationToken).ConfigureAwait(false);

        var body = await ReadLimitedAsync(response.Content, _maxResponseBytes, cancellationToken).ConfigureAwait(false);

        if (!response.IsSuccessStatusCode)
        {
            var error = TryParseError(body);
            throw new InvalidOperationException($"Render failed: {error}");
        }

        if (body.Length == 0)
        {
            throw new InvalidOperationException("Empty PDF response");
        }

        return body;
    }

    public async Task RenderToFileAsync(RenderRequest request, string path, CancellationToken cancellationToken = default)
    {
        var pdf = await RenderAsync(request, cancellationToken).ConfigureAwait(false);
        var directory = Path.GetDirectoryName(path);
        if (!string.IsNullOrWhiteSpace(directory))
        {
            Directory.CreateDirectory(directory);
        }
        await File.WriteAllBytesAsync(path, pdf, cancellationToken).ConfigureAwait(false);
    }

    private static async Task<byte[]> ReadLimitedAsync(HttpContent content, long limitBytes, CancellationToken cancellationToken)
    {
        if (content is null)
        {
            return Array.Empty<byte>();
        }

        using var stream = await content.ReadAsStreamAsync(cancellationToken).ConfigureAwait(false);
        using var ms = new MemoryStream();
        var buffer = new byte[8192];
        long total = 0;
        while (true)
        {
            var read = await stream.ReadAsync(buffer, cancellationToken).ConfigureAwait(false);
            if (read == 0)
            {
                break;
            }
            total += read;
            if (limitBytes > 0 && total > limitBytes)
            {
                throw new InvalidOperationException("Response too large");
            }
            ms.Write(buffer, 0, read);
        }
        return ms.ToArray();
    }

    private string TryParseError(byte[] body)
    {
        try
        {
            var response = JsonSerializer.Deserialize<ErrorResponse>(body, _jsonOptions);
            if (!string.IsNullOrWhiteSpace(response?.Error))
            {
                return response!.Error!;
            }
        }
        catch
        {
            // ignore parse errors
        }

        if (body.Length > 0)
        {
            return Encoding.UTF8.GetString(body);
        }

        return "Unknown error";
    }
}
