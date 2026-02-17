using RazorLight;

namespace PdfFactory.Builders;

internal static class TemplateRenderer
{
    private static readonly RazorLightEngine Engine = new RazorLightEngineBuilder()
        .UseMemoryCachingProvider()
        .Build();

    public static Task<string> RenderAsync<T>(string key, string template, T model)
    {
        return Engine.CompileRenderStringAsync(key, template, model);
    }
}
