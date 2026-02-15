using Microsoft.Extensions.DependencyInjection;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace PdfKit;

public static class ServiceCollectionExtensions
{
    public static IServiceCollection AddPdfKitClient(this IServiceCollection services, Action<PdfClientOptions> configure)
    {
        if (configure is null)
        {
            throw new ArgumentNullException(nameof(configure));
        }

        var options = new PdfClientOptions();
        configure(options);

        services.AddHttpClient(PdfClientOptions.HttpClientName, client =>
        {
            if (!string.IsNullOrWhiteSpace(options.BaseAddress))
            {
                client.BaseAddress = new Uri(options.BaseAddress, UriKind.Absolute);
            }
            if (!string.IsNullOrWhiteSpace(options.UserAgent))
            {
                client.DefaultRequestHeaders.UserAgent.ParseAdd(options.UserAgent);
            }
        });

        services.AddTransient<IPdfClient>(serviceProvider =>
        {
            var factory = serviceProvider.GetRequiredService<IHttpClientFactory>();
            var httpClient = factory.CreateClient(PdfClientOptions.HttpClientName);
            var jsonOptions = options.JsonOptions ?? new JsonSerializerOptions
            {
                DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
                PropertyNamingPolicy = null
            };
            return new PdfClient(httpClient, jsonOptions, options.MaxResponseBytes);
        });

        return services;
    }

    public static IServiceCollection AddPdfKitClient(this IServiceCollection services, string baseAddress)
    {
        if (string.IsNullOrWhiteSpace(baseAddress))
        {
            throw new ArgumentException("Base address must be a non-empty absolute URL.", nameof(baseAddress));
        }

        return services.AddPdfKitClient(options => options.BaseAddress = baseAddress);
    }
}
