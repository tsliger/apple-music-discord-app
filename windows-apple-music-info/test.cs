// Save as NowPlaying.cs
using System;
using System.Threading.Tasks;
using Windows.Media.Control;

class Program
{
    static async Task Main()
    {
        var sessions = await GlobalSystemMediaTransportControlsSessionManager.RequestAsync();
        var currentSession = sessions.GetCurrentSession();

        if (currentSession != null)
        {
            var mediaProperties = await currentSession.TryGetMediaPropertiesAsync();
            Console.WriteLine($"{mediaProperties.Title}|{mediaProperties.Artist}|{mediaProperties.AlbumTitle}");
        }
        else
        {
            Console.WriteLine("NO_TRACK");
        }
    }
}
