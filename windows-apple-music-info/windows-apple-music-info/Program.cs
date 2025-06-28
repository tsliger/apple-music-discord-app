using System;
using System.Linq;
using System.Text.Json;
using System.Threading.Tasks;
using Windows.Media.Capture;
using Windows.Media.Control;

class Program
{
    static GlobalSystemMediaTransportControlsSessionManager? sessionManager;
    static GlobalSystemMediaTransportControlsSession? appleMusicSession;

    static string? lastTrackId = null;
    static bool trackJustChanged = false;
    static TimeSpan? lastDuration = null;
    static TimeSpan? lastPosition = null;

    public static async Task Main()
    {
        sessionManager = await GlobalSystemMediaTransportControlsSessionManager.RequestAsync();

        appleMusicSession = sessionManager.GetSessions()
            .FirstOrDefault(s =>
                s.SourceAppUserModelId?.Contains("AppleMusic", StringComparison.OrdinalIgnoreCase) == true);

        if (appleMusicSession == null)
        {
            Console.WriteLine("{\"status\":\"NO_APPLE_MUSIC_SESSION\"}");
            return;
        }

        appleMusicSession.MediaPropertiesChanged += OnMediaOrPlaybackChanged;
        appleMusicSession.PlaybackInfoChanged += OnMediaOrPlaybackChanged;
        appleMusicSession.TimelinePropertiesChanged += OnTimelinePropertyChange;

        await PrintMediaProperties(appleMusicSession);
        await Task.Delay(-1);
    }

    private static async void OnMediaOrPlaybackChanged(GlobalSystemMediaTransportControlsSession session, object args)
    {
        await Task.Run(async delegate
        {
            await Task.Delay(1000);
            await PrintMediaProperties(session);
        });
    }

    static TimeSpan? prevPosition = null;
    private static async void OnTimelinePropertyChange(GlobalSystemMediaTransportControlsSession session, object args)
    {
        var timeline = session.GetTimelineProperties();
        TimeSpan? currentPosition = timeline.Position;

        if (currentPosition.HasValue && prevPosition.HasValue)
        {
            var diff = currentPosition.Value - prevPosition.Value;
            prevPosition = currentPosition;

            if (Math.Abs(diff.TotalSeconds) > 15)
            {
                await Task.Delay(1000);
                await PrintMediaProperties(session);
            }
        }
        else if (currentPosition.HasValue && prevPosition == null)
        {
            prevPosition = currentPosition;
        }
    }

    private static async Task PrintMediaProperties(GlobalSystemMediaTransportControlsSession session)
    {
        var mediaProperties = await session.TryGetMediaPropertiesAsync();
        var playbackInfo = session.GetPlaybackInfo();
        var timeline = session.GetTimelineProperties();

        string artist = mediaProperties.Artist ?? "";
        string album = mediaProperties.AlbumTitle ?? "";
        string title = mediaProperties.Title ?? "";
        string status = playbackInfo?.PlaybackStatus.ToString() ?? "Unknown";


        TimeSpan duration = timeline.EndTime;
        TimeSpan position = timeline.Position;

        // Generate track ID
        string newTrackId = $"{title}|{artist}|{album}";
        bool isNewTrack = newTrackId != lastTrackId;

        if (isNewTrack)
        {
            lastTrackId = newTrackId;
            trackJustChanged = true;

            lastDuration = null;
            lastPosition = null;
        }

        if (trackJustChanged)
        {
            if (duration.TotalSeconds == 0 || string.IsNullOrWhiteSpace(title))
            {
                return;
            }
            trackJustChanged = false;
        }

        // Fallback split artist/album
        if (string.IsNullOrWhiteSpace(album) && artist.Contains("—"))
        {
            var parts = artist.Split(new[] { '—' }, 2, StringSplitOptions.TrimEntries);
            if (parts.Length == 2)
            {
                artist = parts[0];
                album = parts[1];
            }
        }

        int durationSec = (int)Math.Round(duration.TotalSeconds);
        int positionSec = (int)Math.Round(position.TotalSeconds);
        int lastDurationSec = lastDuration.HasValue ? (int)Math.Round(lastDuration.Value.TotalSeconds) : -1;
        int lastPositionSec = lastPosition.HasValue ? (int)Math.Round(lastPosition.Value.TotalSeconds) : -1;

        if (durationSec == 0 && status != "Stopped")
        {
            return;
        }

        if (durationSec != lastDurationSec || positionSec != lastPositionSec)
        {
            lastDuration = duration;
            lastPosition = position;

            var output = new
            {
                track_name = title,
                artist_name = artist,
                album_name = album,
                player_state = status,
                playhead_time = positionSec.ToString(),
                track_length = durationSec.ToString()
            };

            Console.WriteLine(JsonSerializer.Serialize(output));
        }
    }

}
